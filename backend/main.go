package main

import (
	"backend/channels"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

func deleteMessage(c *gin.Context) {
	fmt.Println("received delete request")
	var deleteRequest messageDeleteRequest
	if err := c.BindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, postId := range deleteRequest.PostIds {
		fmt.Printf("deleting message %d\n", postId)
	}

	publisher, err := StartPublisher()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, p := range deleteRequest.PostIds {
		payload := gin.H{"Payload": p, "Name": "deleteMessage"}
		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		publisher <- data
	}

	c.Status(http.StatusOK)
}

func createMessage(c *gin.Context) {
	fmt.Println("received write request")
	var messagePost messagePostRequest
	if err := c.BindJSON(&messagePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timePosted := time.Now()
	fmt.Printf("received write request for message of '%s'\n", messagePost.Content)

	publisher, err := StartPublisher()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	newMessage := Message{
		Content:    messagePost.Content,
		TimePosted: timePosted,
		PostId:     rand.Int64N(math.MaxInt64),
	}

	payload := gin.H{"Payload": newMessage, "Name": "newMessage"}
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	publisher <- data

	c.Status(http.StatusOK)
}

func watchChat(c *gin.Context, eventSockets *channels.EventSockets) {
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
	fmt.Println("Starting backend...")
	frontendUrl := os.Getenv("FRONTEND_URL")
	fmt.Printf("Looking for connections from %s\n", frontendUrl)

	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	eventSockets := channels.New()

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "DELETE"}
	config.AllowOriginFunc = func(origin string) bool {
		fmt.Println(origin)
		return origin == frontendUrl
	}
	r.Use(cors.New(config))

	var subscriber, err = StartSubscriber()
	for err != nil {
		subscriber, err = StartSubscriber()
		time.Sleep(time.Second * 5)
	}

	go func() {
		for {
			newEvent := <-subscriber
			eventSockets.FanInMessage(newEvent)
		}
	}()

	r.DELETE("/", func(c *gin.Context) {
		deleteMessage(c)
	})

	r.POST("/", func(c *gin.Context) {
		createMessage(c)
	})

	r.GET("/watch", func(c *gin.Context) {
		watchChat(c, eventSockets)
	})

	r.Run()
}
