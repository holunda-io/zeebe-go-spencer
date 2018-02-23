package zeebeutils

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/vmihailenco/msgpack"
	"fmt"
)

func StartProcess(client Client) {
	fmt.Println("Start process ", processId)

	payload := new(GameState)
	payload.BaddieHealth = 100

	instance := zbc.NewWorkflowInstance(processId, -1, make(map[string]interface{}))
	instance.Payload, _ = msgpack.Marshal(payload)
	msg, err := client.zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start Process response: ", msg.String())
}

func GetTask(subscriptionCh chan *zbc.SubscriptionEvent)  (GameState, *zbc.SubscriptionEvent) {
	message := <- subscriptionCh
	var payload GameState

	err := msgpack.Unmarshal(message.Task.Payload, &payload)
	if err != nil {
		panic(err)
	}

	return payload, message
}

func CompleteTask(client Client, state GameState, message *zbc.SubscriptionEvent) {
	p, _ := msgpack.Marshal(state)
	message.Task.Payload = p

	client.zbClient.CompleteTask(message)
}