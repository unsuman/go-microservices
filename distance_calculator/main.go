package main

const KafkaTopic = "obu_data"
const clientEndPoint = "http://127.0.0.1:3300/aggregate"
const grpcEndpoint = "localhost:50051"

func main() {
	service := NewCalculatorService(grpcEndpoint)
	service = NewLoggingMiddleware(service)
	consumer := NewKafkaConsumer(KafkaTopic, service)
	consumer.ConsumeData()
}
