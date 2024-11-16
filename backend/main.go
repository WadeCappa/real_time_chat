package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
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

type SocketData struct {
	inUse    bool
	uniqueId uint64
	channel  chan []byte
}

type EventSockets struct {
	sockets    []SocketData
	lastId     uint64
	globalLock sync.Mutex
}

func (sockets *EventSockets) FanInMessage(event []byte) {
	eventCopy := event
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if c.inUse {
			fmt.Printf("talking to socket %d\n", c.uniqueId)
			c.channel <- eventCopy
		}
	}
}

func (sockets *EventSockets) AddChannel(newChannel chan []byte) uint64 {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	sockets.lastId++
	defer fmt.Printf("created channel of id %d\n", sockets.lastId)
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if !c.inUse {
			c.channel = newChannel
			c.inUse = true
			c.uniqueId = sockets.lastId
			return sockets.lastId
		}
	}

	// Only if we can't re-use a spot in our slice, do we append a new channel
	sockets.sockets = append(sockets.sockets, SocketData{
		channel:  newChannel,
		uniqueId: sockets.lastId,
		inUse:    true,
	})
	return sockets.lastId
}

func (sockets *EventSockets) RemoveChannel(id uint64) {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if c.uniqueId == id {
			close(c.channel)
			c.channel = nil
			c.inUse = false
			fmt.Printf("removed socket %d\n", id)
			return
		}
	}
	fmt.Printf("Could not find socket %d\n", id)
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

func getMessages(conn *sql.DB) ([]Message, error) {
	res, err := conn.Query("select content, time_posted, post_id from user_post")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Close()

	return asMessages(res), nil
}

func deletePost(conn *sql.DB, deleteRequest messageDeleteRequest) error {
	_, err := conn.Query("delete from user_post where user_post.post_id = any($1)", pq.Array(deleteRequest.PostIds))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func writePost(conn *sql.DB, messagePost messagePostRequest, timePosted time.Time) ([]Message, error) {
	post_id_conn, sequence_err := conn.Query("select nextval('post_id_sequence')")
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

	_, err := conn.Query("insert into user_post (time_posted, content, post_id) values ($1, $2, $3)", timePosted.UnixMilli(), messagePost.Content, post_id)
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

	// Lots of coordination overhead with this, also this won't scale
	// past one instance. Better to use pub/sub
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

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ch.Close()

	exchangeError := ch.ExchangeDeclare(
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

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bindError := ch.QueueBind(queue.Name, "", "events", false, nil)
	if bindError != nil {
		fmt.Println(bindError.Error())
		return
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			rabbitEvent := <-msgs
			eventSockets.FanInMessage(rabbitEvent.Body)
		}
	}()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

		runWithDb(func(conn *sql.DB) {
			if err := deletePost(conn, deleteRequest); err != nil {
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
				err = ch.Publish(
					"events",
					queue.Name,
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

		runWithDb(func(conn *sql.DB) {
			messages, err := writePost(conn, messagePost, timePosted)
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
					err = ch.Publish(
						"events",
						queue.Name,
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
		runWithDb(func(conn *sql.DB) {
			res, err := getMessages(conn)
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
