package main

import (
	. "zeebe-go-spencer/zeebeutils"
	"time"
	"math/rand"
	"bufio"
	"os"
	"zeebe-go-spencer/heros"
	"fmt"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	zbClient := CreateNewClient()

	DeployProcess(zbClient)

	heros.InitHero("t", PlayerSetting{NormalAttack:10, SpecialAttack: 30})
	heros.InitHero("b", PlayerSetting{NormalAttack:10, SpecialAttack: 40})
	heros.InitHero("h7", PlayerSetting{NormalAttack:0, SpecialAttack: 50})

	play()
}

func play() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press <enter> to start..")

	for {
		reader.ReadString('\n')
		StartProcess()
	}
}
