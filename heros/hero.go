package heros

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"zeebe-go-spencer/zeebeutils"
	"fmt"
	"math/rand"
)

func InitHero(client zeebeutils.Client, prefix string, setting zeebeutils.PlayerSetting) {
	go attack(prefix, setting.NormalAttack, zeebeutils.CreateSubscription(client, prefix + "-normal"), client)
	go attack(prefix, setting.SpecialAttack, zeebeutils.CreateSubscription(client, prefix + "-special"), client)
	go chooseAttack(prefix, zeebeutils.CreateSubscription(client, prefix + "-choose"), client)
}

func attack(prefix string, damage int, subscriptionCh chan *zbc.SubscriptionEvent, client zeebeutils.Client) {
	for{
		payload, message := zeebeutils.GetTask(subscriptionCh)

		printFormatted(prefix, "attack with ", damage, " damage")
		payload.BaddieHealth = payload.BaddieHealth - damage
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
