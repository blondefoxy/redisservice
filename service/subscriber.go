package service

import (
	"log"
	"redisservice/cache"
)

// Subscriber - читатель сообщений
type Subscriber struct {
	cache *cache.Rediscache
}

// NewSubscriber - конструктор Subscriber
func NewSubscriber() *Subscriber {
	return &Subscriber{
		cache: cache.NewRediscache(redisAddress, 0, 0, 0),
	}
}

// CloseConnection - закрывает подключение к redis
func (s *Subscriber) CloseConnection() {
	s.cache.Close()
}

// Start - запускает получение сообщений
func (s *Subscriber) Start() {

	log.Println("Subscriber starts")

	for {
		message, err := s.cache.Read(messagesList)
		if err != nil && err.Error() != noNewMessages {
			log.Printf("error while reading message: %v", err)
		}
		if len(message) == 0 {
			continue
		}
		err = s.validate(message)
		if err != nil {
			log.Printf("error while validation message: %v", err)
		}
	}
}

var k = 0

// validate валидирует сообщения, выбраковывает 5% как ошибки
func (s *Subscriber) validate(message string) error {
	k++
	if k%20 == 0 {
		return s.cache.Add(errorsList, message)
	} else {
		// как-то обрабатыввается
		log.Printf("message received: %s", message)
	}
	return nil
}
