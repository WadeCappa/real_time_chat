package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/WadeCappa/real_time_chat/auth"
	"github.com/WadeCappa/real_time_chat/chat-kafka-manager/publisher"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	xss "github.com/sahilchopra/gin-gonic-xss-middleware"
)

type messagePostRequest struct {
	Content   string `form:"content" json:"content" xml:"content" binding:"required"`
	ChannelId int64  `form:"channelId" json:"channelId" xml:"channelId" binding:"required"`
}

const (
	DEFAULT_KAFKA_HOSTNAME = "localhost:9092"
	DEFAULT_AUTH_HOSTNAME  = "localhost:50051"
)

var (
	authHostname  = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	kafkaHostname = flag.String("kafka-hostname", DEFAULT_KAFKA_HOSTNAME, "the hostname for kafka")
)

func createMessage(c *gin.Context) {
	userId := c.GetInt64(auth.CURRENT_USER_KEY)
	log.Printf("looking at userid of %d\n", userId)

	log.Println("received write request")
	var messagePost messagePostRequest
	if err := c.BindJSON(&messagePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf("received write request for message of '%s'\n", messagePost.Content)

	err := publisher.PublishChatMessageToChannel([]string{*kafkaHostname}, userId, messagePost.Content, messagePost.ChannelId)
	if err != nil {
		log.Printf("Failed to write message: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("Successfully wrote message")
	c.Status(http.StatusOK)
}

func main() {
	log.Println("Starting backend...")
	frontendUrl := os.Getenv("FRONTEND_URL")
	log.Printf("Looking for connections from %s\n", frontendUrl)

	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(auth.Build(*authHostname))
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "DELETE"}
	config.AllowOriginFunc = func(origin string) bool {
		log.Println(origin)
		return origin == frontendUrl
	}
	r.Use(cors.New(config))

	r.POST("/", func(c *gin.Context) {
		createMessage(c)
	})

	r.Run()
}
