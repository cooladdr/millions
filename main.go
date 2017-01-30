package main

import (
	"fmt"
	"github.com/cooladdr/millionps/dispatcher"
	"github.com/cooladdr/millionps/worker"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	MaxWorker = os.Getenv("MAX_WORKERS")
)

func main() {

	MaxWorker, error := strconv.Atoi(MaxWorker)
	if error != nil {
		MaxWorker = 5
	}
	fmt.Println("MaxWorker: ", MaxWorker)

	dispatcher := dispatcher.NewDispatcher(MaxWorker)
	dispatcher.Run()

	http.HandleFunc("/test", worker.PayloadHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	wait := make(chan int)
	<-wait
}
