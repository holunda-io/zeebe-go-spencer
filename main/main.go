package main

import (
	. "zeebe-go-spencer/zeebeutils"
	"zeebe-go-spencer/heros/terence"
)

func main() {
	zbClient := CreateNewClient()

	DeployProcess(zbClient)

	StartProcess(zbClient)

	terence.InitTerence()
}