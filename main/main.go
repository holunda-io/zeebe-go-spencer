package main

import (
	. "zeebe-go-spencer/zeebeutils"
	"zeebe-go-spencer/heros/terence"
	"time"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	zbClient := CreateNewClient()

	DeployProcess(zbClient)

	StartProcess(zbClient)

	terence.InitTerence()
}