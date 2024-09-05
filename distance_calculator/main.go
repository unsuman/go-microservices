package main

const KafkaTopic = "obu_data"

func main() {
	service := NewCalculatorService()
	consumer := NewKafkaConsumer(KafkaTopic, service)
	consumer.ConsumeData()
}
