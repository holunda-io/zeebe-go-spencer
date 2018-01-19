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

	subscriptionCh, subscription := zeebeutils.CreateSubscription(zbClient, taskChoose)

	zeebeutils.StartGoRoutineToCloseSubscriptionOnExit(zbClient, subscription)

	waitForTaskAndComplete(subscriptionCh, zbClient)
}

func waitForTaskAndComplete(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for {
		fmt.Println("Wait for Task A")

		message := <-subscriptionCh
		var payload map[string]interface{}

		err := msgpack.Unmarshal(message.Task.Payload, &payload)
		if err != nil {
			panic(err)
		}

		fmt.Println("Current health status: ", payload["health"])
		payload["attack"] = "special"

		p, err := msgpack.Marshal(payload)
		message.Task.Payload = p

		// complete task after processing
		response, _ := zbClient.CompleteTask(message)
		fmt.Println("Complete Task Responce: ", response)
	}
}
