[![Go Report Card](https://goreportcard.com/badge/github.com/holunda-io/zeebe-go-spencer)](https://goreportcard.com/report/github.com/holunda-io/zeebe-go-spencer)
[![Build Status](https://travis-ci.org/holunda-io/zeebe-go-spencer.svg?branch=master)](https://travis-ci.org/holunda-io/zeebe-go-spencer)

# Zeebe with Go - Fight of Bud Spencer, Terrence Hill and H7-25

## Setup

1) Install Go => See [Go Setup Doc](https://golang.org/doc/install), don't forget GOPATH!

2) Checkout our project to `$GOPATH/src/github.com/holunda-io/zeebe-go-spencer`

3) Run `setupZbctl.sh` to install latest released Zeebe zbctl

## Run it

1) Build the main-app with: `./createContainer.sh mainapp`

2) Build the hero-app with: `./createContainer.sh heroapp`

3) Start everything with: `docker-compose up`

3) Use `[DOCKER_IP]:8080/fight` to start a process

## Monitoring

1) In docker-compose there is a zeebe-simple-monitor running
   in a dedicated container and should automatically have been started.

2) Check Monitor on: http://[DOCKER_IP]:9080/

3) Add Broker with "zeebe:51015"

More infos: https://docs.zeebe.io/go-client/get-started.html

## Notes 

Bud Spencer - Normale Schläge / Dampfhammer
Terrence Hill - Normale Schläge / Multi Ohrfeige
H7-25 - Photonenkanone
