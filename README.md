[![Go Report Card](https://goreportcard.com/badge/github.com/holunda-io/zeebe-go-spencer)](https://goreportcard.com/report/github.com/holunda-io/zeebe-go-spencer)
[![Build Status](https://travis-ci.org/holunda-io/zeebe-go-spencer.svg?branch=master)](https://travis-ci.org/holunda-io/zeebe-go-spencer)

# Zeebe with Go - Fight of Bud Spencer, Terrence Hill and H7-25

## Setup

1) Install Go => See [Go Setup Doc](https://golang.org/doc/install), don't forget GOPATH!

2) Install Zeebe zbctl => Use latest release from https://github.com/zeebe-io/zbctl/releases

3) If you are not on MacOS replace 0.0.0.0 in docker-compose.yml and main.go with your docker ip.

## Run it

1) Start broker with: `docker-compose up`

2) To run the program: `go run main/main.go`
This will deploy and start an easy process.

3) Use `localhost:8080/fight` to start a process

## Monitoring

1) Download latest simple monitor: https://github.com/zeebe-io/zeebe-simple-monitor/releases

2) Start Monitor `java -jar zeebe-simple-monitor-0.3.0.jar`

3) Check Monitor on: http://127.0.0.1:8080/

4) Add Broker with "[DOCKER_IP]:51015" (e.g. 0.0.0.0:51015 on MacOS)

More infos: https://docs.zeebe.io/go-client/get-started.html

## Notes 

Bud Spencer - Normale Schläge / Dampfhammer
Terrence Hill - Normale Schläge / Multi Ohrfeige
H7-25 - Photonenkanone
