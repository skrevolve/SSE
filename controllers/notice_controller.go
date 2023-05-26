package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skrevolve/sse/util"
	"github.com/valyala/fasthttp"
)

type NoticeController struct{}

type Client struct {
    name   string
    events chan *DashBoard
}

type DashBoard struct {
    Notice 	string
}

func (ct *NoticeController) UrgentNotice(c *fiber.Ctx) error {

	client := &Client{name: c.Context().RemoteIP().String(), events: make(chan *DashBoard, 20)}
    go func() {
		for {
			db := &DashBoard{
				Notice: util.MakeRamdomString(13),
			}
			client.events <- db
		}
	}()

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

		// var i int
		// for {
		// 	i++
		// 	event := client{
		// 		Name: "Notice",
		// 		Event: util.MakeRamdomString(13),
		// 	}
		// 	var buf bytes.Buffer
		// 	enc := json.NewEncoder(&buf)
		// 	enc.Encode(event)
		// 	fmt.Fprintf(w, "data: %v\n\n", buf.String())
		// 	fmt.Printf("data: %v\n", buf.String())
		// 	// msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
		// 	// fmt.Fprintf(w, "data: Message: %s\n\n", msg)
		// 	// fmt.Println(msg)

			err := w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				//break
			}

		// 	time.Sleep(1 * time.Second)
		// }
	}))

	return nil
}