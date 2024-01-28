package main

import (
	"fmt"
	"html/template"
	"htmx/chat"
	"htmx/ticker"
	"htmx/todo"
	"htmx/utils"
	"net/http"
)

type PageData struct {
	Count int
	Todos todo.TodoList
}

func main() {
	var data PageData = PageData{
		0,
		todo.TodoList{},
	}

	var chat chat.Chat = make(chat.Chat)

	var tickers = ticker.CreateTickerListFromArray([]string{"AAAA", "BBBB", "CCCC"})
	tickers.StartAllTickers()

	templates, err := template.ParseFiles(
		"templates/index.html",
		"templates/count.html",
		"templates/todos.html",
		"templates/chat.html",
		"templates/ticker.html",
	)
	utils.HandleError(err)

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		data.Count++
		fmt.Fprintf(w, "%d", data.Count)
	})

	http.HandleFunc("/todo/new", data.Todos.CreateTodoAdder(templates))
	http.HandleFunc("/todo/delete", data.Todos.CreateTodoDeleter(templates))
	http.HandleFunc("/todo/toggle", data.Todos.CreateTodoToggler(templates))

	http.HandleFunc("/chat/send", chat.CreateSender())
	http.HandleFunc("/chat", chat.CreateListener(templates))

	http.HandleFunc("/ticker", ticker.CreateTickerListener(tickers, templates))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index", data)
	})

	fmt.Println("Hello, World!")
	http.ListenAndServe(":8080", nil)
}
