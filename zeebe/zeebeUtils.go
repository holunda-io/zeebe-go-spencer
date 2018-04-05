package zeebe

import (
	"errors"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/zbmsgpack"
	"log"
	"os"
	"os/signal"
)

const defaultTopicName = "default-topic"
const brokerPort = "51015"
const processId = "fight"

type Client struct {
	zbClient            *zbc.Client
	subscriptionHandler chan *zbmsgpack.TaskSubscriptionInfo
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
	client.subscriptionHandler = createSubscriptionHandler(zbClient)
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
	topology, err := client.zbClient.Topology()
	if err != nil {
		log.Fatal("Error happens while loading topology")
		panic(err)
	}

	return topology.PartitionIDByTopicName[topicName] != nil
}

func (client Client) CreateSubscription(task string) chan *zbc.SubscriptionEvent {
	subscriptionCh, subscription, _ := client.zbClient.TaskConsumer(client.topicName, "lockOwner", task)

	client.subscriptionHandler <- subscription
	return subscriptionCh
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
