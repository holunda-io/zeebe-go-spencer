package zeebeutils

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/vmihailenco/msgpack"
	"log"
)

func StartProcess(client Client) {
	log.Println("Start process ",processId)

	payload := new(GameState)
	payload.BaddieHealth = 100

	instance := zbc.NewWorkflowInstance(processId, -1, make(map[string]interface{}))
	instance.Payload, _ = msgpack.Marshal(payload)
	msg, err := client.zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	log.Println("Start Process response: ", msg.String())
}

func ExtractPayload(message *zbc.SubscriptionEvent)  GameState {
	var payload GameState

	err := msgpack.Unmarshal(message.Task.Payload, &payload)
	if err != nil {
		panic(err)
	}

	return payload
}

func CompleteTask(client Client, state GameState, message *zbc.SubscriptionEvent) {
	p, _ := msgpack.Marshal(state)
	message.Task.Payload = p

	client.zbClient.CompleteTask(message)
}