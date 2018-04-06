package zeebeutils

import (
	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
)

func StartProcess(client Client) {
	log.Println("Start process ", processId)

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

func ExtractPayload(event *zbsubscriptions.SubscriptionEvent) GameState {
	var payload GameState

	err := msgpack.Unmarshal(event.Task.Payload, &payload)
	if err != nil {
		panic(err)
	}

	return payload
}

func CompleteTask(client zbsubscribe.ZeebeAPI, state GameState, event *zbsubscriptions.SubscriptionEvent) {
	p, _ := msgpack.Marshal(state)
	event.Task.Payload = p

	client.CompleteTask(event)
}
