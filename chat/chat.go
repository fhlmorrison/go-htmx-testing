package chat

import (
	"fmt"
	"html/template"
	"htmx/utils"
	"net/http"
	"strings"
	"time"
)

type Chat chan string

func (channel Chat) CreateSender() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("chat-input")
		if msg != "" {
			fmt.Printf("Got message:\t %s\n", msg)
			channel <- msg
		}
	}
}

func (channel Chat) CreateListener(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Println("Client connected to chat")
		for {
			select {
			case <-r.Context().Done():
				fmt.Println("Client disconnected from chat")
				return
			case msg := <-channel:
				var response_text string = ""
				for _, v := range strings.Split(msg, " ") {
					fmt.Fprint(w, "event: ChatUpdate\n")
					word := utils.ExecuteTemplateToString(templates, "word", v+" ")

					fmt.Fprintf(w, "data: %s\n\n", word)

					response_text += fmt.Sprintf("%s ", v)
					w.(http.Flusher).Flush()
					time.Sleep(200 * time.Millisecond)
				}
				element := utils.ExecuteTemplateToString(templates, "chat-end", response_text)
				fmt.Fprintf(w, "event: ChatEnd\ndata: %s\n\n", element)
				w.(http.Flusher).Flush()
				fmt.Printf("Sent message:\t %s\n", response_text)
			}
		}
	}
}
