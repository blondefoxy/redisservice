package service

import (
	"log"
	"redisservice/cache"
)

// ErrorHandler - обработчик ошибок
type ErrorHandler struct {
	cache *cache.Rediscache
}

// NewErrorHandler конструктор ErrorHandler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		cache: cache.NewRediscache(redisAddress, 0, 0, 0),
	}
}

// Start - запусает ErrorHandler
func (eh *ErrorHandler) Start() {
	for {
		reply, err := eh.cache.Read(errorsList)
		if err != nil && err.Error() != noNewMessages {
			log.Fatal("cannot get errors", err)
		}

		if reply == "" {
			log.Println("no new errors")
			return
		}

		log.Printf("message with error: %s", reply)
	}
}

// CloseConnection - закрывает соединение с redis
func (eh *ErrorHandler) CloseConnection() {
	eh.cache.Close()
}
