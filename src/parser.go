package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

// SetArgs holds the arguments for the 'set' command
type SetArgs struct {
	Dc       string
	Username string
	Password string
	Domain   string
	Target   string
	Lapspass string
}

// GetArgs holds the arguments for the 'get' command
type GetArgs struct {
	Dc       string
	Username string
	Password string
	Domain   string
	Filter   string
	Outfile  string // To write the results to a csv file
}

// Args holds the parsed command-line arguments
type Args struct {
	set *SetArgs
	get *GetArgs
}

// Partially implemented by Francois Reynaud and modified by Felipe Molina
func parseArgs() *Args {
	parser := argparse.NewParser("golaps", "Get or set the LAPS password of the computers in the domain.")

	// Parse "get" subcommand
	getCmd := parser.NewCommand("get", "")
	getCmd_dc := getCmd.String("D", "dc", &argparse.Options{Required: true, Help: "<IP|FQDN> of the Domain Controller to query."})
	getCmd_username := getCmd.String("u", "username", &argparse.Options{Required: true, Help: "Username to authenticate with."})
	getCmd_password := getCmd.String("p", "password", &argparse.Options{Required: true, Help: "Password to authenticate with."})
	getCmd_domain := getCmd.String("d", "domain", &argparse.Options{Required: true, Help: "Domain of the user authenticating."})
	getCmd_filter := getCmd.String("f", "filter", &argparse.Options{Required: false, Help: "Substring of the computer name (samAccountName) to search for."})
	getCmd_outfile := getCmd.String("o", "out", &argparse.Options{Required: false, Help: "File name of the csv file to write the results."})

	// Parse "set" subcommand
	setCmd := parser.NewCommand("set", "")
	setCmd_dc := setCmd.String("D", "dc", &argparse.Options{Required: true, Help: "<IP|FQDN> of the Domain Controller to target."})
	setCmd_username := setCmd.String("u", "username", &argparse.Options{Required: true, Help: "Username to authenticate with."})
	setCmd_password := setCmd.String("p", "password", &argparse.Options{Required: true, Help: "Password to authenticate with."})
	setCmd_domain := setCmd.String("d", "domain", &argparse.Options{Required: true, Help: "Domain of the user authenticating."})
	setCmd_targetcomputer := setCmd.String("t", "target", &argparse.Options{Required: true, Help: "FQDN of the computer to set the LAPS password."})
	setCmd_lapspass := setCmd.String("P", "lapspass", &argparse.Options{Required: true, Help: "Password to set."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	args := &Args{}
	if setCmd.Happened() {
		args.set = &SetArgs{
			Dc:       *setCmd_dc,
			Username: *setCmd_username,
			Password: *setCmd_password,
			Domain:   *setCmd_domain,
			Target:   *setCmd_targetcomputer,
			Lapspass: *setCmd_lapspass,
		}
	}

	if getCmd.Happened() {
		args.get = &GetArgs{
			Dc:       *getCmd_dc,
			Username: *getCmd_username,
			Password: *getCmd_password,
			Domain:   *getCmd_domain,
			Filter:   *getCmd_filter,
			Outfile:  *getCmd_outfile,
		}
	}

	return args
}
