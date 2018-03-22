package zeebeutils

import (
	"errors"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"log"
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

func CreateNewClient() Client {
	log.Println("Create new zeebe client")

	var client Client

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		log.Fatal("Zeebe Broker not running?")
		panic(errClientStartFailed)
	}

	client.zbClient = zbClient
	client.subscriptionHandler = createSubscriptionHandler(zbClient)

	return client
}

func CreateNewTopicIfNotExists(client Client) {
	log.Printf("Create new topic '%s'", topicName)

	if topicExists(client, topicName) {
		log.Println("Topic does already exist")
		return
	}

	topic, err := client.zbClient.CreateTopic(topicName, 1)
	if err != nil {
		log.Fatal("Could not create topic")
		panic(err)
	}

	log.Println("Created topic: ", topic)
}

func topicExists(client Client, topicName string) bool {
	topology, err := client.zbClient.Topology()
	if err != nil {
		log.Fatal("Error happens while loading topology")
		panic(err)
	}

	return topology.PartitionIDByTopicName[topicName] != nil
}

func DeployProcess(client Client) {
	log.Printf("Deploy '%s' process '%s'\n", zbc.BpmnXml, processFileBpmn)

	response, err := client.zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, processFileBpmn)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	log.Println("Deployed Process response state ", response.State)
}

func CreateSubscription(client Client, task string) chan *zbc.SubscriptionEvent {
	subscriptionCh, subscription, _ := client.zbClient.TaskConsumer(topicName, "lockOwner", task)

	registerSubscription(client.subscriptionHandler, subscription)
	return subscriptionCh
}

func registerSubscription(closeSubscriptionHandler chan *zbmsgpack.TaskSubscriptionInfo, subscription *zbmsgpack.TaskSubscriptionInfo) {
	closeSubscriptionHandler <- subscription
}

func createSubscriptionHandler(zbClient *zbc.Client) chan *zbmsgpack.TaskSubscriptionInfo {
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
		case subscription := <-subscriptionChannel:
			log.Println("Adding subscription to handler")
			subscriptionList = append(subscriptionList, subscription)
		case <-osCh:
			log.Println("Closing subscriptions")
			for e := range subscriptionList {
				zbClient.CloseTaskSubscription(subscriptionList[e])
			}
			os.Exit(0)
		}
	}
}
