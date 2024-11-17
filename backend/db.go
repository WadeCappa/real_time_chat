package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/lib/pq"
)

func getDb() (*sql.DB, error) {
	host := "db"
	port := 5432
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return sql.Open("postgres", psqlInfo)
}

func callWithDb(consumer func(*sql.DB)) {
	db, err := getDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	consumer(db)
}

func runWithDb[T any](consumer func(*sql.DB) (*T, error)) (*T, error) {
	db, err := getDb()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()
	return consumer(db)
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

func writePost(postgres *sql.DB, messagePost messagePostRequest, timePosted time.Time) (*Message, error) {
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

	return &Message{
		Content:    messagePost.Content,
		TimePosted: timePosted,
		PostId:     post_id,
	}, nil
}
