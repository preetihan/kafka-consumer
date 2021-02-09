package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	r := setupRouter()

	topic := getEnv("TOPICS", Topics)
	partition, _ := strconv.Atoi(getEnv("PARTITION", Partition))
	bootstrapServers := getEnv("BOOTSTRAP_SERVERS", BootstrapServers)
	groupID := getEnv("GROUP_ID", GroupID)
	messageMinBytes, _ := strconv.Atoi(getEnv("MESSAGE_MIN_BYTES", MessageMinBytes))
	messageMaxBytes, _ := strconv.Atoi(getEnv("MESSAGE_MAX_BYTES", MessageMaxBytes))

	k := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{bootstrapServers},
		GroupID:   groupID,
		Topic:     topic,
		Partition: partition,
		MinBytes:  messageMinBytes,
		MaxBytes:  messageMaxBytes,
	})

	for {
		m, err := k.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("Consumed message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := k.Close(); err != nil {
		log.Println("failed to close reader:", err)
	}

	r.Run(":" + getEnv("PORT", Port))
}

func setupRouter() *gin.Engine {
	r := gin.New()

	environment := getEnv("ENV", Environment)
	ctxPath := getEnv("CTX_PATH", Ctxpath)
	bootstrapServers := getEnv("BOOTSTRAP_SERVERS", BootstrapServers)
	topics := getEnv("TOPICS", Topics)
	partition := getEnv("PARTITION", Partition)
	groupID := getEnv("GROUP_ID", GroupID)
	messageMinBytes := getEnv("MESSAGE_MIN_BYTES", MessageMinBytes)
	messageMaxBytes := getEnv("MESSAGE_MAX_BYTES", MessageMaxBytes)

	log.Println("ENV = " + environment)
	log.Println("CTX_PATH = " + ctxPath)
	log.Println("BOOTSTRAP_SERVERS = " + bootstrapServers)
	log.Println("TOPICS = " + topics)
	log.Println("PARTITION = " + partition)
	log.Println("GROUP_ID = " + groupID)
	log.Println("MESSAGE_MIN_BYTES = " + messageMinBytes)
	log.Println("MESSAGE_MAX_BYTES = " + messageMaxBytes)

	if environment != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, ctxPath+"/api/health"),
		gin.Recovery(),
	)

	api := r.Group(ctxPath + "/api")
	{
		api.GET("/health", checkHealth)
	}

	return r
}
