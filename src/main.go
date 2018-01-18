package main

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
	"encoding/json"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "src/fight.bpmn"
const processId = "fight"

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func main() {
	zbClient := createNewClient()

	deployProcess(zbClient)

	startProcess(zbClient)
}

func createNewClient() (*zbc.Client) {
	fmt.Println("Create new zeebe client")

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	return zbClient
}

func loadTopologie(zbClient *zbc.Client) {
	fmt.Println("Load broker topologie")

	topology, err := zbClient.Topology()
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(topology, "", "    ")
	fmt.Println("Topologie: ", string(b))
}

func deployProcess(zbClient *zbc.Client) {
	fmt.Printf("Deploy '%s' process '%s'\n", zbc.BpmnXml, processFileBpmn)

	response, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, processFileBpmn)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	fmt.Println("Deployed Process Responce: ", response.String())
}

func startProcess(zbClient *zbc.Client) {
	fmt.Println("Start process ", processId)

	payload := make(map[string]interface{})
	payload["somePayload"] = "31243"
	payload["someOtherPayload"] = "lol"

	instance := zbc.NewWorkflowInstance(processId, -1, payload)
	msg, err := zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start Process responce: ", msg.String())
}
