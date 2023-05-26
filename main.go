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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
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
	app.Get("/sse", adaptor.HTTPHandler(handler(noticeHandler)))
	app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func noticeHandler(w http.ResponseWriter, r *http.Request) {

	dsn := "root:password@tcp(localhost:3306)/app?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	client := &Client{name: r.RemoteAddr, events: make(chan *NoticeUrgent, 10)}
	go updateNoticeUrgent(client, db)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

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

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

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