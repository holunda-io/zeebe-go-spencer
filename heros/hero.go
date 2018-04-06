package heros

import (
	"github.com/holunda-io/zeebe-go-spencer/zeebeutils"
	"log"
	"math/rand"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
)

type handler func(zeebeutils.GameState) zeebeutils.GameState

func InitHero(client zeebeutils.Client, prefix string, setting zeebeutils.PlayerSetting) {
	go zeebeutils.CreateAndRegisterSubscription(client, prefix+"-normal",
		handle(attack(prefix, setting.NormalAttack, setting.AdditionalRange)))

	go zeebeutils.CreateAndRegisterSubscription(client, prefix+"-special",
		handle(attack(prefix, setting.SpecialAttack, setting.AdditionalRange)))

	go zeebeutils.CreateAndRegisterSubscription(client, prefix+"-choose",
		handle(chooseAttack(prefix)))
}

func handle(attackHandler handler) zbsubscribe.TaskSubscriptionCallback {
	return func(clientApi zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
		payload := zeebeutils.ExtractPayload(event)
		newPayload := attackHandler(payload)
		zeebeutils.CompleteTask(clientApi, newPayload, event)
	}
}

func attack(prefix string, damage, additionalRange int) func(zeebeutils.GameState) zeebeutils.GameState {
	return func(payload zeebeutils.GameState) zeebeutils.GameState {
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

func chooseAttack(prefix string) func(zeebeutils.GameState) zeebeutils.GameState {
	return func(payload zeebeutils.GameState) zeebeutils.GameState {
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
