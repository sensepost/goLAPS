# goLAPS
Retrieve LAPS passwords from a domain. 
The tools is inspired in [pyLAPS](https://github.com/p0dalirius/pyLAPS). This project was just a personal excuse to learn Golang.

# Capabilities
* It can get all LAPS passwords from a domain controler using the "get" command
* It can set the target computer LAPS password using the "set" command
* For now, it only works with simple binding on LDAP and LDAPS protocols
* You can provide a filter (-f, --filter) to retrieve computers in the domain that follow a specific patter on their samAccountName

# Usage
## Get LAPS passwords
```bash
./golaps get -h

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

usage: golaps get [-h|--help] -D|--dc "<value>" -u|--username "<value>"
              -p|--password "<value>" -d|--domain "<value>" [-f|--filter
              "<value>"] [-o|--out "<value>"]

              

Arguments:

  -h  --help      Print help information
  -D  --dc        <IP|FQDN> of the Domain Controller to query.
  -u  --username  Username to authenticate with.
  -p  --password  Password to authenticate with.
  -d  --domain    Domain of the user authenticating.
  -f  --filter    Substring of the computer name (samAccountName) to search
                  for.
  -o  --out       File name of the csv file to write the results.
```

## Set LAPS password
```bash
./golaps set -h
usage: golaps set [-h|--help] [-D|--dc "<value>"] [-u|--username "<value>"]
              [-p|--password "<value>"] [-d|--domain "<value>"] [-t|--target
              "<value>"] [-P|--lapspass "<value>"]
Arguments:

  -h  --help      Print help information
  -D  --dc        <IP|FQDN> of the Domain Controller to target.
  -u  --username  Username to authenticate with.
  -p  --password  Password to authenticate with.
  -d  --domain    Domain of the user authenticating.
  -t  --target    FQDN of the computer to set the LAPS password.
  -P  --lapspass  Password to set.
```

# Version
27/05/2024 - SenseCon 2024 Edition

# Authors
Felipe Molina de la Torre ([@felmoltor](https://infosec.exchange/@felmoltor)).
Help from François Reinaud on the argument parsing functionality and Deon Wilemse on the testing infrastructure.
