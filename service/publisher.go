package service

import (
	"log"
	"time"

	"github.com/satori/go.uuid"

	"redisservice/cache"
)

// Publisher - генератор сообщений
type Publisher struct {
	ID    string
	cache *cache.Rediscache
}

// NewPublisher - конструктор Publisher
func NewPublisher() *Publisher {
	return &Publisher{
		ID:    uuid.NewV4().String(),
		cache: cache.NewRediscache(redisAddress, 0, 0, 0),
	}
}

// Start запускает генерацию сообщений
func (p *Publisher) Start() {

	log.Println("Publisher starts")

	for {
		time.Sleep(time.Millisecond * publishTimeMS)

		publisherValid, err := p.cache.Get(publisher)
		if err != nil && err.Error() != noNewMessages {
			log.Printf("can't get publisher's data from redis %s", err.Error())
		}

		// текущий генератор не имеет права отправлять сообщения
		if publisherValid != p.ID {
			continue
		}

		randomString := randomString(stringLength)
		err = p.cache.Add(messagesList, randomString)
		if err != nil {
			log.Printf("can't publish data to redis %s", err.Error())
		} else {
			log.Printf("publish %s", randomString)
		}
	}
}

// HandleState - проверяет состояние, возможно ли генерировать сообщения текущим генератором
// если валидного генератора нет, становится генератором.
// если есть - ждет.
func (p *Publisher) HandleState() {

	for {
		publisherExists, err := p.cache.Get(publisher)
		if err != nil && err.Error() != noNewMessages {
			continue
		}

		if len(publisherExists) == 0 || publisherExists == p.ID {
			// становится генератором
			// если он уже генератор, обновляется запись со свежим ttl,
			// чтобы следующие 500 мс он был действующим генератором
			if len(publisherExists) == 0 {
				log.Printf("current publisher %s is generator", p.ID)
			}
			p.cache.SetWithTTL(publisher, p.ID, publishTimeMS)
		}

		if publisherExists != p.ID {
			continue
		}

		// проверка статуса чаще, чем 500 мс
		time.Sleep(time.Millisecond * publishTimeMS / 2)
	}

}

// CloseConnection - закрывает подключение
func (p *Publisher) CloseConnection() {
	p.cache.Close()
}
