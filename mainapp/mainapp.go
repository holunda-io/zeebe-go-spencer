package main

import (
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/holunda-io/zeebe-go-spencer/zeebe"
	"log"
	"net/http"
	"github.com/holunda-io/zeebe-go-spencer/common"
)

var client Client

func main() {
	log.Println("##### Starting Mainapp #####")

	zeebeHost := common.GetZeebeHost()

	client = NewClientWithDefaultTopic(zeebeHost)
	client.CreateTopicIfNotExists()
	client.DeployProcess("process/fight.bpmn")

	startHttpServer()
}

func startHttpServer() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)
	router.HandleFunc("/fight", fight)

	log.Println("Start http server on 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Show index")
	w.Header().Add("content-type", "text/html")

	fmt.Fprintln(w, "<pre>Zeebe with Go - Fight of Bud Spencer, Terrence Hill and H7-25")
	fmt.Fprintln(w, "-------------------------------------------------------------")
	fmt.Fprintln(w, "<a href=\"/fight\">/fight</a> ... start a fight</pre>")
}

func fight(w http.ResponseWriter, r *http.Request) {
	log.Println("Start fight")

	client := NewClientWithDefaultTopic(common.GetZeebeHost())
	client.StartProcess("fight")
	fmt.Fprint(w, "Started fight")
}
