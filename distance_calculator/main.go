package main

const KafkaTopic = "obu_data"

func main() {
	service := NewCalculatorService()
	service = NewLoggingMiddleware(service)
	consumer := NewKafkaConsumer(KafkaTopic, service)
	consumer.ConsumeData()
}
