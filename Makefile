APP=xm-golang-task
PORT=8080

start: build run

build:
	docker build -t $(APP) .
run:
	docker run -p $(PORT):$(PORT) $(APP)
create-kafka-topic:
	docker exec -ti kafka sh -c "/opt/kafka/bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic xm-golang-task"
get-kafka-messages:
	docker exec -ti kafka sh -c "/opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic xm-golang-task --from-beginning"