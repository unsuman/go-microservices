# Toll-Calculator
This project aims to enhance my understanding of microservices architecture in Go. The Toll-Calculator system is built using Go, gRPC, Kafka, Prometheus, gorilla/websocket, and sirupsen/logrus. This system is used to calculate toll tax for trucks by calculating amount of distance it travelled.
<img alt="Toll-Calculator Architecture" src="/assets/go_microservice.png">

## About
The Toll-Calculator is a microservices-based toll calculation system designed in Go. The system comprises multiple services, including:
- **Aggregator**: Combines data from various services.
- **Data Receiver**: Collects vehicle data from On-Board Units (OBUs).
- **Distance Calculator**: Computes travel distance.
- **Gateway**: Manages communication between services.

Key features:
- Implemented gRPC for inter-service communication.
- Used Kafka for message queuing.
- A mock-up OBU that sends random data to receiver.
- Integrated Prometheus for monitoring.
- Employed sirupsen/logrus for logging.
- Provides real-time toll calculation and billing for vehicles.
- Ensures scalable and fault-tolerant microservices architecture.

## Installing Kafka
```
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true bitnami/kafka:latest 
```

## Installing protobuf compiler (protoc compiler) 
For linux users or (WSL2) 
```
sudo apt install -y protobuf-compiler
```

For Mac users you can use Brew for this
```
brew install protobuff
```

## Installing GRPC and Protobuffer plugins for Golang.
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

2. GRPC 
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. NOTE that you need to set the /go/bin directory in your path
Just like this or whatever your go directly lives.
```
PATH="${PATH}:${HOME}/go/bin"
```

4. install the package dependencies
4.1 protobuffer package
```
go get google.golang.org/protobuf
```
4.2 grpc package
```
go get google.golang.org/grpc/
```

## Installing Prometheus
Install Prometheus in a Docker container
```
docker run -p 9090:9090 -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Installing prometheus golang client
```
go get github.com/prometheus/client_golang/prometheus
```

Installing Prometheus natively on your system
1. Clone the repository
```
git clone https://github.com/promtheus/prometheus.git
```

2. Install
```
cd prometheus
make build
```

3. Run the Prometheus deamon
```
./promtheus --config.file=<your_config_file>yml
```

4. In the projects case that would be (running from inside the project directory)
```
../prometheus/prometheus --config.file=.config/prometheus.yml
```
