package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	name   string
	events chan *NoticeUrgent
}

type NoticeUrgent struct {
	Status 		bool
	Description string
}

type Row struct {
	Status 		bool	`redis:"status"`
	Description string	`redis:"description"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Cache-Control",
		AllowCredentials: true,
	}))
	app.Get("/sse", noticeHandler);
	// app.Server().GetOpenConnectionsCount()
	app.Listen(":3000")
}

func noticeHandler(c *fiber.Ctx) error {

	//db, _ := DatabaseInit()
	rdb, _ := RedisInit()

	client := &Client{name: c.Context().RemoteAddr().String(), events: make(chan *NoticeUrgent, 10)}
	// go updateNoticeUrgentBySql(client, db)
	go updateNoticeUrgentByRedis(client, rdb)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		timeout := time.After(1 * time.Second)
		select {
		case ev := <-client.events:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.Encode(ev)
			fmt.Fprintf(w, "event: %v\n", "notice")
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
			fmt.Printf("event: %v\n", "notice")
			fmt.Printf("data: %v\n", buf.String())
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}

		err := w.Flush()
		if err != nil {
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
		}
	}))

	return nil
}

func updateNoticeUrgentByRedis(client *Client, rdb *redis.Client) {

	var row Row

	for {
		data, _ := rdb.HGet(context.Background(), "urgent", "notice").Result()
		json.Unmarshal([]byte(data), &row)

		status := false
		description := ""
		if row.Status {
			status = row.Status
			description = row.Description
		}

		db := &NoticeUrgent{
			Status: status,
			Description: description,
		}

		client.events <- db
	}
}

func updateNoticeUrgentBySql(client *Client, db *gorm.DB) {

	row := []Row{}

	for {
		db.Raw(`SELECT status, description FROM alert LIMIT 1`).Scan(&row)

		status := false
		description := ""
		if row[0].Status {
			status = row[0].Status
			description = row[0].Description
		}

		db := &NoticeUrgent{
			Status: status,
			Description: description,
		}

		client.events <- db
	}
}

func DatabaseInit() (*gorm.DB, error) {

	dsn := "root:password@tcp(localhost:3306)/app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Second * 10)

	return db, nil
}

func RedisInit() (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr: "192.168.0.172:6379",
		Password: "",
		DB: 0,
	})

	return rdb, nil
}