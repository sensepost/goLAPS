package main

import (
	"crypto/tls"
	csv "encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	filepath "path/filepath"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// Struct to hold the credentials of the GoLaps client
type Credentials struct {
	Username string
	Password string
	Domain   string
}

// Struct to hold the configuration of the GoLaps client
type GoLaps struct {
	Credentials        Credentials // The credentials to bind to the DC
	TargetDC           string      // This could be an IP or the FQDN of the DC
	BaseDN             string      // The base DN to search for the computers
	ComputernameFilter string      // The filter to search for the computernames to read the LAPS password from
	Targetcomputer     string      // The FQDN of the computer to set the LAPS password
}

func (gl GoLaps) InitLogging() {
	// Create the logs folder if it does not exits
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	// Log to a file
	logPath := filepath.Join("logs", time.Now().Format("2006-01-02_150405")+"_golaps.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	// defer file.Close() // I'm new to this language. I guess I shouldn't close the file stream if I want to keep writing to the log file, innit?

	// Direct logs to both stdout and the file
	log.SetOutput(io.MultiWriter(os.Stdout, file))
}

// ConnectToDC connects to the DC using LDAPS or LDAP
func (gl GoLaps) ConnectToDC() (*ldap.Conn, error) {
	ldapsURL := "ldaps://" + gl.TargetDC + ":636"
	l, err := ldap.DialURL(ldapsURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true})) // I don't care about the cert, take me to jail, officer
	if err != nil {
		log.Println("Error connecting to LDAP server using LDAPS: " + err.Error())
		ldapURL := "ldap://" + gl.TargetDC + ":389"
		l, err = ldap.DialURL(ldapURL)
		if err != nil {
			log.Fatal("Error connecting to LDAP server using LDAP: " + err.Error())
		}
	}
	// defer l.Close() // Close the connection when the function returns
	return l, err
}

func (gl GoLaps) CloseConnection(l *ldap.Conn) {
	l.Close()
}

// Bind to the DC using the username and password if they are set
func (gl GoLaps) BindToDC(l *ldap.Conn) error {
	// Choose the type of binding, Unauthenticated or Authenticated
	if gl.Credentials.Username == "" || gl.Credentials.Password == "" {
		log.Println("Username or password not set. Using unauthenticated bind")
		err := l.UnauthenticatedBind("")
		if err != nil {
			log.Println("Error unauthenticated binding to LDAP server" + err.Error())
		}
	} else {
		err := l.Bind(gl.Credentials.Username+"@"+gl.Credentials.Domain, gl.Credentials.Password)
		if err != nil {
			log.Fatal("Error binding to LDAP server" + err.Error())
		}
		return err
	}
	return nil
}

// Receive the Pointer to the LDAP connection and return and updates the domain DN
func (gl *GoLaps) GetDomainDN(l *ldap.Conn) (string, error) {
	// Get the domain DN (Distinguished Name)
	// Create a search request
	searchRequest := ldap.NewSearchRequest(
		"", // Base DN (empty string for Root DSE)
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",                // Search filter for domain object
		[]string{"defaultNamingContext"}, // Attribute to retrieve (DN)
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal("Error searching LDAP server" + err.Error())
	}

	domainDN := ""
	// Print the DN of the domain
	if len(sr.Entries) > 0 {
		domainDN = sr.Entries[0].GetAttributeValue("defaultNamingContext")
		log.Println("DN of the domain:", domainDN)
	} else {
		log.Println("Domain not found")
	}

	gl.BaseDN = domainDN

	return domainDN, err
}

func (gl GoLaps) SearchComputersWithLaps(l *ldap.Conn, csvfile string) (*ldap.SearchResult, error) {
	searchFilter := fmt.Sprintf("(&(objectCategory=computer)(ms-MCS-AdmPwd=*)(sAMAccountName=*%s*))", ldap.EscapeFilter(gl.ComputernameFilter))
	if gl.ComputernameFilter == "" {
		searchFilter = "(&(objectCategory=computer)(ms-MCS-AdmPwd=*))"
	}
	// Construct LDAP search request
	searchRequest := ldap.NewSearchRequest(
		gl.BaseDN, // Base DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter, // Search filter
		[]string{"sAMAccountName", "ms-Mcs-AdmPwd"}, // Attributes to retrieve
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal("Error searching for Computers with ms-Mcs-AdmPwd attributes" + err.Error())
	}

	// Open the CSV file to write the search results
	var w *csv.Writer
	if csvfile != "" {
		// Export to CSV
		file, err := os.OpenFile(csvfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening CSV file" + err.Error())
		}
		defer file.Close()
		w = csv.NewWriter(file)
		defer w.Flush()
		row := []string{"sAMAccountNamE", "ms-Mcs-AdmPwd"}
		w.Write(row)
	}

	// Print search results
	log.Println("Search Results:")
	for _, entry := range sr.Entries {
		log.Printf("DN: %s\n", entry.DN)
		name := entry.GetAttributeValue("sAMAccountName")
		pass := entry.GetAttributeValue("ms-Mcs-AdmPwd")
		log.Printf(" %s: %s\n", name, pass)
		// Write to CSV if the writer is not nil
		if w != nil {
			row := []string{name, pass}
			w.Write(row)
		}
	}

	return sr, err
}

func (gl GoLaps) SetLapsPassword(l *ldap.Conn, LapsPassword string) error {
	// query for the target computer
	searchFilter := fmt.Sprintf("(&(objectCategory=computer)(sAMAccountName=%s))", ldap.EscapeFilter(gl.Targetcomputer))
	searchRequest := ldap.NewSearchRequest(
		gl.BaseDN, // Base DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter, // Search filter
		[]string{"sAMAccountName", "ms-Mcs-AdmPwd"}, // Attributes to retrieve
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal("Error searching for target computer %s" + err.Error())
		return err
	}
	targetComputerDN := ""
	for _, entry := range sr.Entries {
		targetComputerDN = entry.DN
		log.Printf("Target computer DN: %s\n", targetComputerDN)
		for _, attr := range entry.Attributes {
			log.Printf("  %s: %v\n", attr.Name, attr.Values)
		}
	}
	// TODO: Finish this
	log.Println("Set the LAPS password of " + gl.Targetcomputer + " to " + LapsPassword)
	modifyRequest := ldap.ModifyRequest{
		DN:       targetComputerDN,
		Changes:  nil,
		Controls: nil,
	}
	modifyRequest.Replace("ms-Mcs-AdmPwd", []string{LapsPassword})
	// Perform the modification
	err = l.Modify(&modifyRequest)
	if err != nil {
		log.Fatalf("Failed to modify attribute 'ms-Mcs-AdmPwd': %v", err)
		return err
	}

	log.Println("Successfully modified the ms-Mcs-AdmPwd attribute")

	return nil
}
