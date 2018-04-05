package zeebeutils

import (
	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
)

func (client Client) DeployProcess(processFile string) {
	log.Printf("Deploy '%s' process '%s'\n", zbc.BpmnXml, processFile)

	response, err := client.zbClient.CreateWorkflowFromFile(client.topicName, zbc.BpmnXml, processFile)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}
	log.Println("Deployed Process response state ", response.State)
}

func (client Client) StartProcess() {
	log.Println("Start process ", processId)

	payload := new(GameState)
	payload.BaddieHealth = 100

	instance := zbc.NewWorkflowInstance(processId, -1, make(map[string]interface{}))
	instance.Payload, _ = msgpack.Marshal(payload)
	msg, err := client.zbClient.CreateWorkflowInstance(client.topicName, instance)

	if err != nil {
		panic(err)
	}

	log.Println("Start Process response: ", msg.String())
}

func ExtractPayload(message *zbc.SubscriptionEvent) GameState {
	var payload GameState

	err := msgpack.Unmarshal(message.Task.Payload, &payload)
	if err != nil {
		panic(err)
	}

	return payload
}

func (client Client) CompleteTask(state GameState, message *zbc.SubscriptionEvent) {
	p, _ := msgpack.Marshal(state)
	message.Task.Payload = p

	client.zbClient.CompleteTask(message)
}
