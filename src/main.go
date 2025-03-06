package main

import (
	"fmt"
	"log"
)

func printBanner() {
	fmt.Print(`
    ________  ________  ___       ________  ________  ________      
    |\   ____\|\   __  \|\  \     |\   __  \|\   __  \|\   ____\     
    \ \  \___|\ \  \|\  \ \  \    \ \  \|\  \ \  \|\  \ \  \___|_    
     \ \  \  __\ \  \\\  \ \  \    \ \   __  \ \   ____\ \_____  \   
      \ \  \|\  \ \  \\\  \ \  \____\ \  \ \  \ \  \___|\|____|\  \  
       \ \_______\ \_______\ \_______\ \__\ \__\ \__\     ____\_\  \ 
        \|_______|\|_______|\|_______|\|__|\|__|\|__|    |\_________\
                                                         \|_________|
    Retrieve LAPS passwords from a domain controler
    (author: @felmoltor)
    Inspired by pyLAPS (https://github.com/p0dalirius/pyLAPS)            

`)
}

func getLapsPassword(args *GetArgs) {
	gl := GoLaps{
		Credentials: Credentials{
			Username: args.Username,
			Password: args.Password,
			Domain:   args.Domain,
		},
		TargetDC:           args.Dc,
		ComputernameFilter: args.Filter,
	}
	gl.InitLogging()
	l, err := gl.ConnectToDC()
	if err != nil {
		log.Fatal("Error connecting to DC" + err.Error())
	}
	err = gl.BindToDC(l)
	if err != nil {
		log.Fatal("Error binding to DC" + err.Error())
	}
	_, err = gl.GetDomainDN(l)
	if err != nil {
		log.Fatal("Error getting domain DN" + err.Error())
	}
	_, err = gl.SearchComputersWithLaps(l, args.Outfile)
	if err != nil {
		log.Fatal("Error searching for computers with LAPS" + err.Error())
	}
	gl.CloseConnection(l)
}

func setLapsPassword(args *SetArgs) {
	gl := GoLaps{
		Credentials: Credentials{
			Username: args.Username,
			Password: args.Password,
			Domain:   args.Domain,
		},
		TargetDC:       args.Dc,
		Targetcomputer: args.Target,
	}
	gl.InitLogging()
	l, err := gl.ConnectToDC()
	if err != nil {
		log.Fatal("Error connecting to DC" + err.Error())
	}
	err = gl.BindToDC(l)
	if err != nil {
		log.Fatal("Error binding to DC" + err.Error())
	}
	_, err = gl.GetDomainDN(l)
	if err != nil {
		log.Fatal("Error getting domain DN" + err.Error())
	}
	// TODO: Complete the method to set the LAPS password of this computer
	err = gl.SetLapsPassword(l, args.Lapspass)
	if err != nil {
		log.Fatal("Error setting the LAPS password" + err.Error())
	}
	gl.CloseConnection(l)
}

// main
func main() {
	printBanner()
	args := parseArgs()

	if args.get != nil {
		getLapsPassword(args.get)
	} else if args.set != nil {
		setLapsPassword(args.set)
	} else {
		fmt.Println("No command specified")
	}
}
