package main

import (
	. "zeebe-go-spencer/zeebeutils"
	"zeebe-go-spencer/heros"
	"time"
	"math/rand"
	"bufio"
	"os"
	"fmt"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	client := CreateNewClient()

	DeployProcess(client)

	heros.InitHero(client, "t", PlayerSetting{NormalAttack:10, SpecialAttack: 30, AdditionalRange: 5})
	heros.InitHero(client, "b", PlayerSetting{NormalAttack:10, SpecialAttack: 40, AdditionalRange: 10})
	heros.InitHero(client, "h7", PlayerSetting{NormalAttack:0, SpecialAttack: 50, AdditionalRange: 20})

	play()
}

func play() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press <enter> to start..")

	for {
		reader.ReadString('\n')
		StartProcess(CreateNewClient())
	}
}