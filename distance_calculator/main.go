package main

const KafkaTopic = "obu_data"

func main() {
	consumer := NewKafkaConsumer(KafkaTopic)
	consumer.ConsumeData()
}
