package main

// Some constant
const (
	Component        = "kafka-consumer"
	Environment      = "local"
	Ctxpath          = "/kafka"
	Port             = "8080"
	BootstrapServers = "kafka-1-service:9092"
	Topics           = "test-topic"
	Partition        = "0"
	GroupID          = "test_group"
	MessageMinBytes  = "10e3" // 10KB
	MessageMaxBytes  = "10e6" // 10MB
)
