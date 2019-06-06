package service

import "math/rand"

const (
	redisAddress = "127.0.0.1:6379"

	//stringLength длина случайной строки
	stringLength = 20

	// errorsList список для записи ошибок
	errorsList = "errors"

	// messagesList список для записи сообщений
	messagesList = "messages"

	// publisher хранит текущего генератора
	publisher = "currentPublisher"

	// PublishTimeMS интервал генерации сообщений
	publishTimeMS = 500

	noNewMessages = "redigo: nil returned"
)

// randomString генерирует случайные строки
func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}