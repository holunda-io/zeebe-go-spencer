package terence

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"zeebe-go-spencer/zeebeutils"
	"math/rand"
)

const taskChoose = "t-choose"
const taskNormal = "t-normal"
const taskSpecial = "t-special"

func InitTerence() {

	zbClient := zeebeutils.CreateNewClient()

	subscriptionCh := zeebeutils.CreateSubscription(zbClient, taskChoose)
	subscriptionChNormal := zeebeutils.CreateSubscription(zbClient, taskNormal)
	subscriptionChSpecial := zeebeutils.CreateSubscription(zbClient, taskSpecial)

	go attackNormal(subscriptionChNormal, zbClient)
	go attackSpecial(subscriptionChSpecial, zbClient)
	chooseAttack(subscriptionCh, zbClient)
}

func attackSpecial(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for{
		fmt.Println("Wait for special attack")

		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		fmt.Println("Current health status: ", payload.Health)
		fmt.Println("Do special attack")
		payload.Health = payload.Health - 50
		fmt.Println("==> New health status: ", payload.Health)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		zbClient.CompleteTask(message)
	}
}

func attackNormal(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for{
		fmt.Println("Wait for normal attack")

		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		fmt.Println("Current health status: ", payload.Health)
		fmt.Println("Do normal attack")
		payload.Health = payload.Health - 10
		fmt.Println("==> New health status: ", payload.Health)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		zbClient.CompleteTask(message)
	}
}

func chooseAttack(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for {
		fmt.Println("Wait for Choose attack")

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

		fmt.Println("Choosen attack ", payload.Decision)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		zbClient.CompleteTask(message)
	}
}
