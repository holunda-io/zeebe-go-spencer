package main

import (
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/holunda-io/zeebe-go-spencer/zeebeutils"
	"log"
	"net/http"
)

func main() {
	client := CreateNewClient()

	CreateNewTopicIfNotExists(client)

	DeployProcess(client)

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

	fmt.Fprintln(w, "Zeebe with Go - Fight of Bud Spencer, Terrence Hill and H7-25")
	fmt.Fprintln(w, "-------------------------------------------------------------")
	fmt.Fprintln(w, "/fight ... start a fight")
}

func fight(w http.ResponseWriter, r *http.Request) {
	log.Println("Start fight")

	StartProcess(CreateNewClient())
	fmt.Fprint(w, "Started fight")
}
