package main

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
	"encoding/json"
	"os/signal"
	"os"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "src/simpleProcess.bpmn"
const processId = "simpleProcess"
const taskA = "task-a"

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func main() {

	zbClient := createNewClient()

	deployProcessBpmn(zbClient)
	//deployProcessYaml(zbClient) <- is working but than you don't see a BPMN in Monitor

	startProcess(zbClient)

	subscriptionCh, subscription := createSubscriptionForTaskA(zbClient)

	startGoRoutineToCloseSubscriptionOnExit(zbClient, subscription)

	waitForTaskAndComplete(subscriptionCh, zbClient)
}

func createNewClient() (*zbc.Client) {
	fmt.Println("Create new zeebe client")

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	return zbClient
}

func deployProcessBpmn(zbClient *zbc.Client) {
	deployProcess(zbClient, zbc.BpmnXml, processFileBpmn)
}

func deployProcess(zbClient *zbc.Client, resourceType, path string) {
	fmt.Printf("Deploy '%s' process '%s'\n", resourceType, processFileYaml)

	response, err := zbClient.CreateWorkflowFromFile(topicName, resourceType, path)
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

func createSubscriptionForTaskA(zbClient *zbc.Client) (chan *zbc.SubscriptionEvent, *zbmsgpack.TaskSubscription) {
	fmt.Println("Open task subscription for Task A")

	subscriptionCh, subscription, _ := zbClient.TaskConsumer(topicName, "lockOwner", taskA)
	return subscriptionCh, subscription
}

func waitForTaskAndComplete(subscriptionCh chan *zbc.SubscriptionEvent, zbClient *zbc.Client) {
	for {
		fmt.Println("Wait for Task A")

		message := <-subscriptionCh
		fmt.Println("Message of task A subscription: ", message.String())

		// complete task after processing
		response, _ := zbClient.CompleteTask(message)
		fmt.Println("Complete Task Responce: ", response)
	}
}

func startGoRoutineToCloseSubscriptionOnExit(zbClient *zbc.Client, subscription *zbmsgpack.TaskSubscription) {
	fmt.Println("Create go routine which waits for app interrrupt")

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	go func() {
		<-osCh
		fmt.Println("Closing subscription.")
		_, err := zbClient.CloseTaskSubscription(subscription)
		if err != nil {
			fmt.Println("failed to close subscription: ", err)
		} else {
			fmt.Println("Subscription closed.")
		}
		os.Exit(0)
	}()
}
