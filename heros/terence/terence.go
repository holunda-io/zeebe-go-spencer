package terence

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"zeebe-go-spencer/zeebeutils"
)

const taskChoose = "t-choose"
const taskNormal = "t-normal"
const taskSpecial = "t-special"

func InitTerence() {

	zbClient := zeebeutils.CreateNewClient()

	subscriptionCh := zeebeutils.CreateSubscription(zbClient, taskChoose)
	subscriptionCh3 := zeebeutils.CreateSubscription(zbClient, taskSpecial)

	go attackSpecial(subscriptionCh3, zbClient)
	chooseAttack(subscriptionCh, zbClient)
}

func attackSpecial(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for{
		fmt.Println("Wait to attack")

		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}
		fmt.Println("attacking")

		payload.Health = payload.Health - 10
		fmt.Println("Current health status: ", payload.Health)

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		// complete task after processing
		response, _ := zbClient.CompleteTask(message)
		fmt.Println("Complete Task Response: ", response)
	}
}

func chooseAttack(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for {
		fmt.Println("Wait for Task A")

		message := <-subscriptionCh
		var payload zeebeutils.GameState

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		fmt.Println("Current health status: ", payload.Health)
		payload.Decision = "special"

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		// complete task after processing
		response, _ := zbClient.CompleteTask(message)
		fmt.Println("Complete Task Response: ", response)
	}
}
