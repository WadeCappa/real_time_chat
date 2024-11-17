package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	xss "github.com/sahilchopra/gin-gonic-xss-middleware"
)

type messagePostRequest struct {
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
}

type messageDeleteRequest struct {
	PostIds []int64 `form:"postIds" json:"postIds" xml:"postIds" binding:"required"`
}

type Message struct {
	Content    string
	TimePosted time.Time
	PostId     int64
}

func deleteMessage(c *gin.Context, writeChannel *amqp.Channel) {
	fmt.Println("received delete request")
	var deleteRequest messageDeleteRequest
	if err := c.BindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, postId := range deleteRequest.PostIds {
		fmt.Printf("deleting message %d\n", postId)
	}

	callWithDb(func(postgres *sql.DB) {
		if err := deletePost(postgres, deleteRequest); err != nil {
			c.JSON(http.StatusBadRequest, "Failed to delete post")
			return
		}
	})

	for _, p := range deleteRequest.PostIds {
		payload := gin.H{"Payload": p, "Name": "deleteMessage"}
		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err := writeToRabbit(writeChannel, data); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func createMessage(c *gin.Context, writeChannel *amqp.Channel) {
	fmt.Println("received write request")
	var messagePost messagePostRequest
	if err := c.BindJSON(&messagePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timePosted := time.Now()
	fmt.Printf("received write request for message of '%s'\n", messagePost.Content)

	res, err := runWithDb(func(postgres *sql.DB) (*Message, error) {
		messages, err := writePost(postgres, messagePost, timePosted)
		if err != nil {
			return nil, err
		}

		return messages, nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to write post")
		fmt.Println(err.Error())
		return
	}

	payload := gin.H{"Payload": *res, "Name": "newMessage"}
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := writeToRabbit(writeChannel, data); err != nil {
		fmt.Println(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func loadMessages(c *gin.Context) {
	fmt.Println("received get request")
	res, err := runWithDb(func(postgres *sql.DB) (*[]Message, error) {
		res, err := getMessages(postgres)
		if err != nil {
			return nil, err
		}

		return &res, nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get posts")
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, *res)
}

func watchEvents(c *gin.Context, eventSockets *EventSockets) {
	socketChannel := make(chan []byte)
	newId := eventSockets.AddChannel(socketChannel)
	defer eventSockets.RemoveChannel(newId)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Flush()
	fmt.Println("connected with client")

	for {
		newEvent := <-socketChannel
		select {
		case <-c.Request.Context().Done():
			// exit when done
			return
		default:
			// no-op, keep going
		}
		c.Writer.Write([]byte(fmt.Sprintf("data: %s\n", newEvent)))
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()
	}
}

func main() {
	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	var eventSockets EventSockets
	eventSockets.lastId = 0

	mode := os.Getenv("MODE")
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "DELETE"}
	config.AllowOriginFunc = func(origin string) bool {
		fmt.Println(origin)
		switch mode {
		case "production":
			return origin == "https://cantseewater.online"
		case "local":
			return origin == "http://localhost:3000"
		default:
			fmt.Println(fmt.Errorf("did not recognize the deployement mode: %s", mode))
			return false
		}
	}
	r.Use(cors.New(config))

	rabbit, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rabbit.Close()

	writeChannel, err := rabbit.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer writeChannel.Close()

	exchangeError := writeChannel.ExchangeDeclare(
		"events",
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if exchangeError != nil {
		fmt.Println(exchangeError.Error())
		return
	}

	killChannel := make(chan bool)
	listener, err := consumeRabbitEvents(rabbit, killChannel)
	if err != nil {
		killChannel <- true
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			newEvent := <-listener
			eventSockets.FanInMessage(newEvent.Body)
		}
	}()

	r.DELETE("/", func(c *gin.Context) {
		deleteMessage(c, writeChannel)
	})

	r.POST("/", func(c *gin.Context) {
		createMessage(c, writeChannel)
	})

	r.GET("/", loadMessages)

	r.GET("/watch", func(c *gin.Context) {
		watchEvents(c, &eventSockets)
	})

	r.Run()
}
