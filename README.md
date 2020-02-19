# collectCertInfo

Given a list of hosts, this tool will attempt to collect some basic information from their SSL certificate (if any). It will attempt to fetch the target's Common Name, Organization Name, and Subject Alternative Names. The results will be stored in a SQLite database for further use. 

Host lists can be passed as domain names or IP addresses (such as *google.com*, or *127.0.0.1*).

## Installation

```
go get github.com/mattn/go-sqlite3
go build
```

## Usage

```
./collectCertInfo hosts.txt
```

## Output

Output is stored in the 'certs.sqlite3' database. All future runs of this script will add to the existing database, and will skip over any already existing hosts.

Database fields:

```
host TEXT NOT NULL PRIMARY KEY
common_name TEXT
organization_name TEXT
subject_alt_name TEXT
```
