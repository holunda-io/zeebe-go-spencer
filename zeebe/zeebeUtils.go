package zeebe

import (
	"errors"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"os"
	"os/signal"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
)

const defaultTopicName = "default-topic"
const brokerPort = "51015"
const processId = "fight"

type Client struct {
	zbClient            *zbc.Client
	subscriptionHandler chan *zbsubscribe.TaskSubscription
	topicName			string
}

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func NewClientWithDefaultTopic(brokerAddr string) Client {
	log.Println("Create new zeebe client")

	var client Client

	zbClient, err := zbc.NewClient(brokerAddr + ":" + brokerPort)
	if err != nil {
		log.Fatal("Zeebe Broker not running?")
		panic(errClientStartFailed)
	}

	client.zbClient = zbClient
	log.Println("Create Zeebe Subscription Handler")
	client.subscriptionHandler = createSubscriptionHandler()
	client.topicName = defaultTopicName

	return client
}

func (client Client) CreateTopicIfNotExists() {
	log.Printf("Create new topic '%s'", client.topicName)

	if topicExists(client, client.topicName) {
		log.Println("Topic does already exist")
		return
	}

	topic, err := client.zbClient.CreateTopic(client.topicName, 1)
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

func (client Client) CreateAndRegisterSubscription(task string, cb zbsubscribe.TaskSubscriptionCallback) {
	subscription, err := client.zbClient.TaskSubscription(client.topicName, "lockOwner", task, 32, cb)
	if err != nil {
		panic("Unable to open subscription")
	}
	log.Printf("Creating subscribtion to task %s", task)
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
