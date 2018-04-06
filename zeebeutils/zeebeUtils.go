package zeebeutils

import (
	"errors"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"os"
	"os/signal"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "process/fight.bpmn"
const processId = "fight"

type Client struct {
	zbClient            *zbc.Client
	subscriptionHandler chan *zbsubscribe.TaskSubscription
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
	client.subscriptionHandler = createSubscriptionHandler()

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
	topology, err := client.zbClient.RefreshTopology()
	if err != nil {
		log.Fatal("Error happens while loading topology")
		panic(err)
	}

	return topology.PartitionIDByTopicName[topicName] != nil
}

func DeployProcess(client Client) {
	log.Printf("Deploy '%s' process '%s'\n", zbcommon.BpmnXml, processFileBpmn)

	response, err := client.zbClient.CreateWorkflowFromFile(topicName, zbcommon.BpmnXml, processFileBpmn)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	log.Println("Deployed Process response state ", response.State)
}

func CreateAndRegisterSubscription(client Client, task string, cb zbsubscribe.TaskSubscriptionCallback)  {
	subscription, err := client.zbClient.TaskSubscription(topicName, "lockOwner", task, 32, cb)

	if err != nil {
		panic("Unable to open subscription")
	}

	registerSubscription(client.subscriptionHandler, subscription)

	subscription.Start()
}

func registerSubscription(closeSubscriptionHandler chan *zbsubscribe.TaskSubscription, subscription *zbsubscribe.TaskSubscription) {
	closeSubscriptionHandler <- subscription
}

func createSubscriptionHandler() chan *zbsubscribe.TaskSubscription {
	subscriptionChannel := make(chan *zbsubscribe.TaskSubscription)
	go subscriptionHandler(subscriptionChannel)
	return subscriptionChannel
}

func subscriptionHandler(subscriptionChannel chan *zbsubscribe.TaskSubscription) {
	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	var subscriptionList []*zbsubscribe.TaskSubscription

	for {
		select {
		case subscription := <-subscriptionChannel:
			log.Println("Adding subscription to handler")
			subscriptionList = append(subscriptionList, subscription)
		case <-osCh:
			log.Println("Closing subscriptions")
			for e := range subscriptionList {
				subscriptionList[e].Close()
			}
			os.Exit(0)
		}
	}
}
