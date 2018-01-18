package zeebeutils

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"os"
	"os/signal"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "process/fight.bpmn"
const processId = "fight"

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func CreateNewClient() (*zbc.Client) {
	fmt.Println("Create new zeebe client")

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	return zbClient
}

func DeployProcess(zbClient *zbc.Client) {
	fmt.Printf("Deploy '%s' process '%s'\n", zbc.BpmnXml, processFileBpmn)

	response, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, processFileBpmn)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	fmt.Println("Deployed Process Responce: ", response.String())
}

func StartProcess(zbClient *zbc.Client) {
	fmt.Println("Start process ", processId)

	payload := make(map[string]interface{})
	payload["health"] = "100"

	instance := zbc.NewWorkflowInstance(processId, -1, payload)
	msg, err := zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start Process responce: ", msg.String())
}

func CreateSubscription(zbClient *zbc.Client, task string) (chan *zbc.SubscriptionEvent, *zbmsgpack.TaskSubscription) {
	fmt.Println("Open task subscription for ", task)

	subscriptionCh, subscription, _ := zbClient.TaskConsumer(topicName, "lockOwner", task)
	return subscriptionCh, subscription
}

func StartGoRoutineToCloseSubscriptionOnExit(zbClient *zbc.Client, subscription *zbmsgpack.TaskSubscription) {
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
