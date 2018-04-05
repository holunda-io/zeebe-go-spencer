package main

import (
	. "github.com/holunda-io/zeebe-go-spencer/zeebeutils"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"math/rand"
	"os"
	"time"
	"github.com/holunda-io/zeebe-go-spencer/common"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	zeebeHost := common.GetZeebeHost()
	client := NewClientWithDefaultTopic(zeebeHost)

	settings := map[string]PlayerSetting{
		"t":  {NormalAttack: 10, SpecialAttack: 30, AdditionalRange: 5},
		"b":  {NormalAttack: 10, SpecialAttack: 40, AdditionalRange: 10},
		"h7": {NormalAttack: 0, SpecialAttack: 50, AdditionalRange: 20},
	}

	hero := os.Getenv("HERO")

	InitHero(client, hero, settings[hero])
}

type handler func(GameState) GameState

func InitHero(client Client, prefix string, setting PlayerSetting) {
	normalSub := client.CreateSubscription(prefix+"-normal")
	specialSub := client.CreateSubscription(prefix+"-special")
	chooseSub := client.CreateSubscription(prefix+"-choose")

	for {
		select {
		case message := <-normalSub:
			handle(attack(prefix, setting.NormalAttack, setting.AdditionalRange), client, message)
		case message := <-specialSub:
			handle(attack(prefix, setting.SpecialAttack, setting.AdditionalRange), client, message)
		case message := <-chooseSub:
			handle(chooseAttack(prefix), client, message)
		}
	}
}

func handle(attackHandler handler, client Client, message *zbc.SubscriptionEvent) {
	payload := ExtractPayload(message)
	newPayload := attackHandler(payload)
	client.CompleteTask(newPayload, message)
}

func attack(prefix string, damage, additionalRange int) func(GameState) GameState {
	return func(payload GameState) GameState {
		doneDamage := damage + calculateAdditionalRange(additionalRange)
		printFormatted(prefix, "Attack with ", doneDamage, " damage")
		result := payload.BaddieHealth - doneDamage
		if result < 0 {
			payload.BaddieHealth = 0
		} else {
			payload.BaddieHealth = result
		}
		printFormatted(prefix, "==> New health status: ", payload.BaddieHealth)
		return payload
	}
}

func chooseAttack(prefix string) func(GameState) GameState {
	return func(payload GameState) GameState {
		switch rand.Intn(2) {
		case 1:
			payload.Decision = "special"
		default:
			payload.Decision = "normal"
		}

		printFormatted(prefix, "Chosen attack: ", payload.Decision)
		return payload
	}
}

func printFormatted(prefix string, msg ...interface{}) {
	log.Println("[", prefix, "\t] ", msg)
}

func calculateAdditionalRange(additionalRange int) int {
	return rand.Intn(additionalRange)
}
