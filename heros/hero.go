package heros

import (
	"github.com/holunda-io/zeebe-go-spencer/zeebeutils"
	"math/rand"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
)

type handler func(zeebeutils.GameState) zeebeutils.GameState

func InitHero(client zeebeutils.Client, prefix string, setting zeebeutils.PlayerSetting) {
	normalSub := zeebeutils.CreateSubscription(client, prefix + "-normal")
	specialSub := zeebeutils.CreateSubscription(client, prefix + "-special")
	chooseSub := zeebeutils.CreateSubscription(client, prefix + "-choose")

	for {
		select {
		case message := <- normalSub:
			handle(attack(prefix, setting.NormalAttack, setting.AdditionalRange), client, message)
		case message := <- specialSub:
			handle(attack(prefix, setting.SpecialAttack, setting.AdditionalRange), client, message)
		case message := <- chooseSub:
			handle(chooseAttack(prefix), client, message)
		}
	}
}

func handle(attackHandler handler, client zeebeutils.Client, message *zbc.SubscriptionEvent) {
	payload := zeebeutils.ExtractPayload(message)
	newPayload := attackHandler(payload)
	zeebeutils.CompleteTask(client, newPayload, message)
}

func attack(prefix string, damage, additionalRange int)  func(zeebeutils.GameState) zeebeutils.GameState {
	return func(payload zeebeutils.GameState) zeebeutils.GameState {
		doneDamage := damage + calculateAdditionalRange(additionalRange)
		printFormatted(prefix, "Attack with ",doneDamage," damage")
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
	log.Println("[",prefix,"\t] ",msg)
}

func calculateAdditionalRange(additionalRange int) int {
	return rand.Intn(additionalRange)
}
