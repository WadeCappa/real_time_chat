package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	xss "github.com/sahilchopra/gin-gonic-xss-middleware"
)

type messagePostRequest struct {
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
}

type messageDeleteRequest struct {
	PostIds []int64 `form:"postIds" json:"postIds" xml:"postIds" binding:"required"`
}

type message struct {
	Content    string
	TimePosted time.Time
	PostId     int64
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

func asMessages(messages *sql.Rows) []message {
	var res []message
	for messages.Next() {
		var content string
		var timePosted int64
		var postId int64
		if err := messages.Scan(&content, &timePosted, &postId); err == nil {
			res = append(res, message{
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

func getMessages(conn *sql.DB) ([]message, error) {
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

func writePost(conn *sql.DB, messagePost messagePostRequest, timePosted time.Time) error {
	post_id_conn, sequence_err := conn.Query("select nextval('post_id_sequence')")
	if sequence_err != nil {
		fmt.Println(sequence_err)
		return sequence_err
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
		return err
	}

	return nil
}

func main() {
	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	mode := os.Getenv("MODE")
	fmt.Println(mode)
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

		for _, p := range deleteRequest.PostIds {
			fmt.Printf("deleting message %d\n", p)
		}

		runWithDb(func(conn *sql.DB) {
			if err := deletePost(conn, deleteRequest); err != nil {
				c.JSON(http.StatusBadRequest, "Failed to delete post")
				return
			}
		})

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
		fmt.Println(messagePost.Content)

		runWithDb(func(conn *sql.DB) {
			err := writePost(conn, messagePost, timePosted)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "Failed to write post")
				return
			}
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

	r.Run()
}
