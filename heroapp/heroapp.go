package main

import (
	. "github.com/holunda-io/zeebe-go-spencer/zeebe"
	"log"
	"math/rand"
	"os"
	"github.com/holunda-io/zeebe-go-spencer/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"github.com/holunda-io/zeebe-go-spencer/game"
)

var settings = map[string]game.PlayerSetting{
"t":  {NormalAttack: 10, SpecialAttack: 30, AdditionalRange: 5},
"b":  {NormalAttack: 10, SpecialAttack: 40, AdditionalRange: 10},
"h7": {NormalAttack: 0, SpecialAttack: 50, AdditionalRange: 20},
}

func main() {

	zeebeHost := common.GetZeebeHost()
	client := NewClientWithDefaultTopic(zeebeHost)

	hero := os.Getenv("HERO")
	initHero(client, hero, settings[hero])
}

type handler func(game.State) game.State

func initHero(client Client, prefix string, setting game.PlayerSetting) {
	go client.CreateAndRegisterSubscription(prefix+"-normal", handle(attack(prefix, setting.NormalAttack, setting.AdditionalRange)))
	go client.CreateAndRegisterSubscription(prefix+"-special", handle(attack(prefix, setting.SpecialAttack, setting.AdditionalRange)))
	go client.CreateAndRegisterSubscription(prefix+"-choose", handle(chooseAttack(prefix)))
}

func handle(attackHandler handler) zbsubscribe.TaskSubscriptionCallback {
	return func(clientApi zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
		log.Printf("Incoming event: %s", event)
		payload := ExtractPayload(event)
		newPayload := attackHandler(payload)
		CompleteTask(clientApi, newPayload, event)
	}
}

func attack(prefix string, damage, additionalRange int) func(game.State) game.State {
	return func(payload game.State) game.State {
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

func chooseAttack(prefix string) func(game.State) game.State {
	return func(payload game.State) game.State {
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
