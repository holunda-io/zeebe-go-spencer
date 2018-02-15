package heros

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"zeebe-go-spencer/zeebeutils"
	"math/rand"
)

func InitHero(prefix string, setting zeebeutils.PlayerSetting) {
	zbClient := zeebeutils.CreateNewClient()
	go attack(prefix, setting.NormalAttack, zeebeutils.CreateSubscription(zbClient, prefix + "-normal"), zbClient)
	go attack(prefix, setting.SpecialAttack, zeebeutils.CreateSubscription(zbClient, prefix + "-special"), zbClient)
	go chooseAttack(prefix, zeebeutils.CreateSubscription(zbClient, prefix + "-choose"), zbClient)
}

func attack(prefix string, damage int, subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for{
		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		printFormatted(prefix, "attack with ", damage, " damage")
		payload.BaddieHealth = payload.BaddieHealth - damage
		printFormatted(prefix, "==> New health status: ", payload.BaddieHealth)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		zbClient.CompleteTask(message)
	}
}

func chooseAttack(prefix string, subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for {
		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		switch rand.Intn(2) {
			case 1:
				payload.Decision = "special"
			default:
				payload.Decision = "normal"
		}

		printFormatted(prefix, "Choosen attack ", payload.Decision)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		zbClient.CompleteTask(message)
	}
}

func printFormatted(prefix string, msg ...interface{}) {
	fmt.Println("[",prefix,"\t] ", msg[:])
}
