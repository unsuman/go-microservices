package main

const KafkaTopic = "obu_data"
const clientEndPoint = "http://127.0.0.1:3300/aggregate"

func main() {
	service := NewCalculatorService(clientEndPoint)
	service = NewLoggingMiddleware(service)
	consumer := NewKafkaConsumer(KafkaTopic, service)
	consumer.ConsumeData()
}
