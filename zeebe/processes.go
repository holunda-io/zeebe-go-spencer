package zeebe

import (
	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/holunda-io/zeebe-go-spencer/game"
)

func (client Client) DeployProcess(processFile string) {
	log.Printf("Deploy '%s' process '%s'\n", zbcommon.BpmnXml, processFile)

	response, err := client.zbClient.CreateWorkflowFromFile(client.topicName, zbcommon.BpmnXml, processFile)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}
	log.Println("Deployed Process response state ", response.State)
}

func (client Client) StartProcess(processId string) {
	log.Println("Start process ", processId)

	payload := game.NewGame()

	instance := zbc.NewWorkflowInstance(processId, -1, make(map[string]interface{}))
	instance.Payload, _ = msgpack.Marshal(payload)
	msg, err := client.zbClient.CreateWorkflowInstance(client.topicName, instance)

	if err != nil {
		panic(err)
	}

	log.Println("Start Process response: ", msg.String())
}

func ExtractPayload(event *zbsubscriptions.SubscriptionEvent) game.State {
	var payload game.State

	err := msgpack.Unmarshal(event.Task.Payload, &payload)
	if err != nil {
		panic(err)
	}

	return payload
}

func CompleteTask(clientApi zbsubscribe.ZeebeAPI, state game.State, event *zbsubscriptions.SubscriptionEvent) {
	p, _ := msgpack.Marshal(state)
	event.Task.Payload = p

	clientApi.CompleteTask(event)
}
