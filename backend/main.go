package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	xss "github.com/sahilchopra/gin-gonic-xss-middleware"
)

type message struct {
	content    string
	timePosted time.Time
}

type messagePostRequest struct {
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
}

func toJson(messages []message) gin.H {
	var res []gin.H
	for _, m := range messages {
		res = append(res, gin.H{
			"content":    m.content,
			"timePosted": m.timePosted.UnixMilli(),
		})
	}
	return gin.H{
		"messages": res,
	}
}

func main() {
	var messages []message

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
		newMessage := message{messagePost.Content, timePosted}
		messages = append(messages, newMessage)
		c.JSON(http.StatusOK, toJson([]message{newMessage}))
	})

	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, toJson(messages))
	})

	r.Run()
}
