package main

import (
	. "zeebe-go-spencer/zeebeutils"
	"time"
	"math/rand"
	"zeebe-go-spencer/heros"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	client := CreateNewClient()

	DeployProcess(client)

	go heros.InitHero(client, "t", PlayerSetting{NormalAttack:10, SpecialAttack: 30, AdditionalRange: 5})
	go heros.InitHero(client, "b", PlayerSetting{NormalAttack:10, SpecialAttack: 40, AdditionalRange: 10})
	go heros.InitHero(client, "h7", PlayerSetting{NormalAttack:0, SpecialAttack: 50, AdditionalRange: 20})

	play()
}

func play() {
	for {
		client := CreateNewClient()
		StartProcess(client)
	}
}