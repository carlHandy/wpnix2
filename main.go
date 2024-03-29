package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	nginxAvailable = "/etc/nginx/sites-available"
	nginxEnabled   = "/etc/nginx/sites-enabled"
	webDir         = "/var/www"
	webUser        = "www-data:www-data"
)

func showHelp() {
	fmt.Println(`Usage: wpnix [-h] [-d DOMAIN] [-u DBUSER] [-p DBPASS] [-n DBNAME] [-H DBHOST] [-e EMAIL]

This program installs WordPress and sets up an Nginx server block.

    -h          display this help and exit
	-e EMAIL    email address for Let's Encrypt SSL
    -d DOMAIN   specify the domain name
    -u DBUSER   database user
    -p DBPASS   database password
    -n DBNAME   database name
    -H DBHOST   database host`)
}

func main() {
	var (
		domain, dbUser, dbPass, dbName, dbHost, email string
		help                                          bool
	)

	flag.StringVar(&domain, "d", "", "Domain name")
	flag.StringVar(&dbUser, "u", "", "Database user")
	flag.StringVar(&dbPass, "p", "", "Database password")
	flag.StringVar(&dbName, "n", "", "Database name")
	flag.StringVar(&dbHost, "H", "", "Database host")
	flag.StringVar(&email, "e", "", "Email address for Let's Encrypt SSL")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	if domain == "" || dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || email == "" {
		fmt.Println("Error: All parameters are required.")
		showHelp()
		return
	}

	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root.")
		return
	}

	checkAndInstallPHP()
	checkAndInstallPerl()
	checkAndInstallNginx()
	installWordPress(domain, dbUser, dbPass, dbName, dbHost)
	checkAndInstallCertbot()
	createNginxConfig(domain)
	createSymlinkIfNotExists(filepath.Join(nginxAvailable, domain), filepath.Join(nginxEnabled, domain))
	validateAndReloadNginx()
	configureLetsEncryptSSL(domain, email)
	finalizeSetupAndRestartNginx(domain)
}
