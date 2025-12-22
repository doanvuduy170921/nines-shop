package rabbitmq

type RabbitMQService interface {
	Send()
	Receive()
}
