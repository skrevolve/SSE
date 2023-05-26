// package main

// import (
// 	"log"

// 	"github.com/skrevolve/sse/routes"
// 	"github.com/skrevolve/sse/server"
// )

// func main() {

// 	app := server.Create()

// 	routes.Init(app)

// 	if err := server.Listen(app); err != nil {
// 		log.Panic(err)
// 	}
// }

// main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	name   string
	events chan *NoticeUrgent
}
type NoticeUrgent struct {
	Alert  bool
	Notice string
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Cache-Control",
		AllowCredentials: true,
	}))
	// app.Get("/sse", adaptor.HTTPHandler(handler(noticeHandler)))
	app.Get("/sse", noticeHandler);
	// app.Server().GetOpenConnectionsCount()
	app.Listen(":3000")
}

// func handler(f http.HandlerFunc) http.Handler {
// 	return http.HandlerFunc(f)
// }

func noticeHandler(c *fiber.Ctx) error {

	db, _ := DatabaseInit()
	client := &Client{name: c.Context().RemoteAddr().String(), events: make(chan *NoticeUrgent, 10)}
	go updateNoticeUrgent(client, db)

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
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
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

// func noticeHandler(w http.ResponseWriter, r *http.Request) {

// 	client := &Client{name: r.RemoteAddr, events: make(chan *NoticeUrgent, 10)}
// 	go updateNoticeUrgent(client, db)

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Transfer-Encoding", "chunked")

// 	timeout := time.After(1 * time.Second)
// 	select {
// 	case ev := <-client.events:
// 		var buf bytes.Buffer
// 		enc := json.NewEncoder(&buf)
// 		enc.Encode(ev)
// 		fmt.Fprintf(w, "data: %v\n\n", buf.String())
// 		fmt.Printf("data: %v\n", buf.String())
// 	case <-timeout:
// 		fmt.Fprintf(w, ": nothing to sent\n\n")
// 	}

// 	if f, ok := w.(http.Flusher); ok {
// 		f.Flush()
// 	}
// }

func updateNoticeUrgent(client *Client, db *gorm.DB) {

	type Alert struct {
		Idx	int
		Alert bool
		Notice string
	}

	result := []Alert{}

	for {
		db.Raw(`SELECT alert, notice FROM alert LIMIT 1`).Scan(&result)

		alert := false
		notice := ""

		if result[0].Alert {
			alert = true
			notice = result[0].Notice
		}

		db := &NoticeUrgent{
			Alert: alert,
			Notice: notice,
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
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}