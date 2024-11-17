package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

type Event struct {
	Name    string
	Payload []byte
}

func runWithDb(consumer func(*sql.DB)) {
	host := "db"
	port := 5432
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	consumer(db)
}

func asMessages(messages *sql.Rows) []Message {
	var res []Message
	for messages.Next() {
		var content string
		var timePosted int64
		var postId int64
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

func getMessages(postgres *sql.DB) ([]Message, error) {
	res, err := postgres.Query("select content, time_posted, post_id from user_post")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Close()

	return asMessages(res), nil
}

func deletePost(postgres *sql.DB, deleteRequest messageDeleteRequest) error {
	_, err := postgres.Query("delete from user_post where user_post.post_id = any($1)", pq.Array(deleteRequest.PostIds))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func writePost(postgres *sql.DB, messagePost messagePostRequest, timePosted time.Time) ([]Message, error) {
	post_id_conn, sequence_err := postgres.Query("select nextval('post_id_sequence')")
	if sequence_err != nil {
		fmt.Println(sequence_err)
		return nil, sequence_err
	}
	var post_id int64
	for post_id_conn.Next() {
		if err := post_id_conn.Scan(&post_id); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("got new id from database: %d\n", post_id)

	_, err := postgres.Query("insert into user_post (time_posted, content, post_id) values ($1, $2, $3)", timePosted.UnixMilli(), messagePost.Content, post_id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return []Message{{
		Content:    messagePost.Content,
		TimePosted: timePosted,
		PostId:     post_id,
	}}, nil
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

	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println("received delete request")
		var deleteRequest messageDeleteRequest
		if err := c.BindJSON(&deleteRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, postId := range deleteRequest.PostIds {
			fmt.Printf("deleting message %d\n", postId)
		}

		runWithDb(func(postgres *sql.DB) {
			if err := deletePost(postgres, deleteRequest); err != nil {
				c.JSON(http.StatusBadRequest, "Failed to delete post")
				return
			}
		})

		func() {
			for _, p := range deleteRequest.PostIds {
				payload := gin.H{"Payload": p, "Name": "deleteMessage"}
				data, err := json.Marshal(payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				err = writeChannel.Publish(
					"events",
					"",
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        data,
					})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
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

		runWithDb(func(postgres *sql.DB) {
			messages, err := writePost(postgres, messagePost, timePosted)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "Failed to write post")
				return
			}
			func() {
				for _, m := range messages {
					payload := gin.H{"Payload": m, "Name": "newMessage"}
					data, err := json.Marshal(payload)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					err = writeChannel.Publish(
						"events",
						"",
						false,
						false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        data,
						})
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}
			}()
		})

		c.Status(http.StatusOK)
	})

	r.GET("/", func(c *gin.Context) {
		fmt.Println("received get request")
		runWithDb(func(postgres *sql.DB) {
			res, err := getMessages(postgres)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "Failed to get posts")
				return
			}
			c.JSON(http.StatusOK, res)
		})
	})

	r.GET("/watch", func(c *gin.Context) {
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
	})

	r.Run()
}
