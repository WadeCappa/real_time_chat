package main

import (
	"database/sql"
	"fmt"
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
	consumer(db)
}

func toJson(messages *sql.Rows) gin.H {
	var res []gin.H
	for messages.Next() {
		var content string
		var timePosted int
		if err := messages.Scan(&content, &timePosted); err != nil {
			res = append(res, gin.H{
				"content":    content,
				"timePosted": timePosted,
			})
		}
	}
	return gin.H{
		"messages": res,
	}
}

func getRows(conn *sql.DB) gin.H {
	res, err := conn.Query("select * from user_post")
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	return toJson(res)
}

func main() {
	mode := os.Getenv("GIN_MODE")

	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())
	r.SetTrustedProxies(nil)

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowOriginFunc = func(origin string) bool {
		switch mode {
		case "release":
			return origin == "https://cantseewater.online"
		default:
			return origin == "localhost:3000"
		}
	}
	r.Use(cors.New(config))

	r.POST("/write", func(c *gin.Context) {
		var messagePost messagePostRequest
		if err := c.ShouldBindJSON(&messagePost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		timePosted := time.Now()

		fmt.Printf("received new message: %s; posted at %s", messagePost.Content, timePosted.Format("ANSIC"))

		runWithDb(func(conn *sql.DB) {
			query := fmt.Sprintf("insert into user_post (time_posted, content) values (%d, '%s')", timePosted.UnixMilli(), messagePost.Content)
			fmt.Println(query)
			res, err := conn.Exec(query)
			fmt.Println(res)
			fmt.Println(err)
			c.JSON(http.StatusOK, getRows(conn))
		})
	})

	r.GET("/get", func(c *gin.Context) {
		runWithDb(func(conn *sql.DB) {
			c.JSON(http.StatusOK, getRows(conn))
		})
	})

	r.Run()
}
