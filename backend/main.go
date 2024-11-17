package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	xss "github.com/sahilchopra/gin-gonic-xss-middleware"
)

const kafkaTopic string = "message-events"

type messagePostRequest struct {
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
}

type messageDeleteRequest struct {
	PostIds []int64 `form:"postIds" json:"postIds" xml:"postIds" binding:"required"`
}

type Message struct {
	Content    string
	TimePosted time.Time
	PostId     uint64
}

func getRandomInt() uint64 {
	return rand.Uint64()
}

func asMessages(messages *sql.Rows) []Message {
	var res []Message
	for messages.Next() {
		var content string
		var timePosted int64
		var postId uint64
		if err := messages.Scan(&content, &timePosted, &postId); err == nil {
			res = append(res, Message{
				Content:    content,
				TimePosted: time.Unix(0, timePosted*int64(time.Millisecond)),
				PostId:     postId,
			})
		} else {
			fmt.Println(err)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].TimePosted.After(res[j].TimePosted)
	})
	return res
}

func consumeEvents(consumer func([]byte) bool) {
	groupId := string(getRandomInt())
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka",
		"acks":              "all",
		"auto.offset.reset": "earliest",
		"group.id":          groupId,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer kafkaConsumer.Close()

	err = kafkaConsumer.SubscribeTopics([]string{kafkaTopic}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		done := consumer(msg.Value)
		if done {
			return
		}
	}
}

func main() {
	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka",
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer kafkaProducer.Close()

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

	r.DELETE("/", func(c *gin.Context) {
		fmt.Println("received delete request")
		var deleteRequest messageDeleteRequest
		if err := c.BindJSON(&deleteRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, postId := range deleteRequest.PostIds {
			fmt.Printf("deleting message %d\n", postId)
		}

		func() {
			for _, p := range deleteRequest.PostIds {
				payload := gin.H{"Payload": p, "Name": "deleteMessage"}
				data, err := json.Marshal(payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				topic := kafkaTopic
				kafkaProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          []byte(data),
				}, nil)
			}
		}()

		c.Status(http.StatusOK)
	})

	r.POST("/", func(c *gin.Context) {
		fmt.Println("received write request")
		var messagePost messagePostRequest
		if err := c.BindJSON(&messagePost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		timePosted := time.Now()
		fmt.Printf("received write request for message of '%s'\n", messagePost.Content)

		message := Message{
			Content:    messagePost.Content,
			TimePosted: timePosted,
			PostId:     getRandomInt(),
		}

		payload := gin.H{"Payload": message, "Name": "newMessage"}
		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		topic := kafkaTopic
		resultChannel := make(chan kafka.Event)
		kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic},
			Value:          []byte(data),
		}, resultChannel)

		kafkaProducer.Flush(999999999)
		res := <-resultChannel
		c.JSON(http.StatusOK, res)
	})

	r.GET("/", func(c *gin.Context) {
		fmt.Println("received get request")
		c.Status(http.StatusOK)
	})

	r.GET("/watch", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Flush()
		fmt.Println("connected with client")

		consumeEvents(func(newEvent []byte) bool {
			select {
			case <-c.Request.Context().Done():
				// exit when done
				return true
			default:
				// no-op, keep going
			}
			c.Writer.Write([]byte(fmt.Sprintf("data: %s\n", newEvent)))
			c.Writer.Write([]byte("\n"))
			c.Writer.Flush()

			return false
		})
	})

	r.Run()
}
