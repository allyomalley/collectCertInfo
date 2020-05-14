package main

import (
	"bufio"
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"strings"
)

var db *sql.DB

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please pass in a file.")
		return
	}

	filename := os.Args[1]

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transCfg}

	var err error
	db, err = sql.Open("sqlite3", "./certs.sqlite3")
	checkErr(err)

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS targets (host TEXT NOT NULL PRIMARY KEY, common_name TEXT, organization_name TEXT, subject_alt_name TEXT)")
	statement.Exec()

	targetFile, err := os.Open(filename)
	checkErr(err)
	scanner := bufio.NewScanner(targetFile)
	scanner.Split(bufio.ScanLines)
	var targets []string
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}
	targetFile.Close()

	for _, host := range targets {
		fetchCertificate(host, client)
	}

	db.Close()
}

func fetchCertificate(host string, client *http.Client) {
	fmt.Println("https://" + host)
	response, err := client.Get("https://" + host)

	if err != nil {
		return
	}

	var commonName = ""
	var org = ""
	var altDns = ""
	for _, cert := range response.TLS.PeerCertificates {

		if strings.Join(cert.Subject.Organization, ",") != "" {
			org = strings.Join(cert.Subject.Organization, ",")
		}

		if cert.Subject.CommonName != "" {
			commonName = cert.Subject.CommonName
		}

		dnsString := strings.Join(cert.DNSNames, ",")
		if dnsString != "" {
			altDns = dnsString
		}

		if commonName != "" || org != "" || altDns != "" {
			break
		}
	}

	if commonName != "" || org != "" || altDns != "" {
		insertRow(host, commonName, org, altDns)
	}

	response.Body.Close()
}

func insertRow(host string, common_name string, organization_name string, subject_alt_name string) {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO targets (host, common_name, organization_name, subject_alt_name) VALUES (?, ?, ?, ?)")
	checkErr(err)
	_, err = stmt.Exec(host, common_name, organization_name, subject_alt_name)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
