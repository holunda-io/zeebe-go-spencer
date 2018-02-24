package heros

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"zeebe-go-spencer/zeebeutils"
	"fmt"
	"math/rand"
)

func InitHero(client zeebeutils.Client, prefix string, setting zeebeutils.PlayerSetting) {
	go attack(prefix, setting.NormalAttack, setting.AdditionalRange, zeebeutils.CreateSubscription(client, prefix + "-normal"), client)
	go attack(prefix, setting.SpecialAttack, setting.AdditionalRange, zeebeutils.CreateSubscription(client, prefix + "-special"), client)
	go chooseAttack(prefix, zeebeutils.CreateSubscription(client, prefix + "-choose"), client)
}

func attack(prefix string, damage, additionalRange int, subscriptionCh chan *zbc.SubscriptionEvent, client zeebeutils.Client) {
	for{
		payload, message := zeebeutils.GetTask(subscriptionCh)

		doneDamage := damage + calculateAdditionalRange(additionalRange)
		printFormatted(prefix, "attack with ", doneDamage, " damage")
		payload.BaddieHealth = payload.BaddieHealth - doneDamage
		printFormatted(prefix, "==> New health status: ", payload.BaddieHealth)

		zeebeutils.CompleteTask(client, payload, message)
	}
}

func chooseAttack(prefix string, subscriptionCh chan *zbc.SubscriptionEvent, client zeebeutils.Client) {
	for {
		payload, message := zeebeutils.GetTask(subscriptionCh)

		switch rand.Intn(2) {
			case 1:
				payload.Decision = "special"
			default:
				payload.Decision = "normal"
		}

		printFormatted(prefix, "Choosen attack ", payload.Decision)

		zeebeutils.CompleteTask(client, payload, message)
	}
}

func printFormatted(prefix string, msg ...interface{}) {
	fmt.Println("[",prefix,"\t] ", msg[:])
}

func calculateAdditionalRange(additionalRange int) int {
	return rand.Intn(additionalRange)
}
