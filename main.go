package main

import (
	"flag"
	"log"
	"redisservice/service"
)

func main() {

	getErrors := flag.Bool("getErrors", false, "get handlered errors from redis")
	flag.Parse()

	if *getErrors {
		log.Println("запустился обработчик ошибок")

		errHandler := service.NewErrorHandler()
		defer errHandler.CloseConnection()
		errHandler.Start()
	} else {
		publisher := service.NewPublisher()
		defer publisher.CloseConnection()

		go publisher.HandleState()
		go publisher.Start()

		subscriber := service.NewSubscriber()
		defer subscriber.CloseConnection()
		subscriber.Start()
	}

}
