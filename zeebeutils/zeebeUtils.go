package zeebeutils

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"errors"
	"fmt"
	"os"
	"os/signal"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "process/fight.bpmn"
const processId = "fight"

type Client struct {
	zbClient            *zbc.Client
	subscriptionHandler chan *zbmsgpack.TaskSubscriptionInfo
}

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func CreateNewClient() (Client) {
	fmt.Println("Create new zeebe client")

	var client Client

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	client.zbClient = zbClient
	client.subscriptionHandler = createSubscriptionHandler(zbClient)

	return client
}


func DeployProcess(client Client) {
	fmt.Printf("Deploy '%s' process '%s'\n", zbc.BpmnXml, processFileBpmn)

	response, err := client.zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, processFileBpmn)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	fmt.Println("Deployed Process Response: ", response.String())
}

func CreateSubscription(client Client, task string) (chan *zbc.SubscriptionEvent) {
	subscriptionCh, subscription, _ := client.zbClient.TaskConsumer(topicName, "lockOwner", task)

	registerSubscription(client.subscriptionHandler, subscription)
	return subscriptionCh
}

func registerSubscription(closeSubscriptionHandler chan *zbmsgpack.TaskSubscriptionInfo, subscription *zbmsgpack.TaskSubscriptionInfo) {
	closeSubscriptionHandler <- subscription
}


func createSubscriptionHandler(zbClient *zbc.Client) (chan *zbmsgpack.TaskSubscriptionInfo) {
	subscriptionChannel := make(chan *zbmsgpack.TaskSubscriptionInfo)
	go subscriptionHandler(zbClient, subscriptionChannel)
	return subscriptionChannel
}

func subscriptionHandler(zbClient *zbc.Client, subscriptionChannel chan *zbmsgpack.TaskSubscriptionInfo) {
	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	var subscriptionList []*zbmsgpack.TaskSubscriptionInfo

	for {
		select {
			case subscription := <- subscriptionChannel:
				fmt.Println("Adding subscription to handler")
				subscriptionList = append(subscriptionList, subscription)
			case <-osCh:
				fmt.Println("Closing subscriptions")
				for e := range subscriptionList {
					zbClient.CloseTaskSubscription(subscriptionList[e])
				}
				os.Exit(0)
		}
	}
}