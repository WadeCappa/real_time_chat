package main

import (
	"fmt"
	"net/http"
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

	r := gin.Default()
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://cantseewater.online", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.SetTrustedProxies(nil)

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
