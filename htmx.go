package main

import (
	"fmt"
	"html/template"
	"htmx/chat"
	"htmx/ticker"
	"htmx/todo"
	"htmx/utils"
	"net/http"
	"slices"
)

type SymbolActivity struct {
	Active   []string
	Inactive []string
}

type PageData struct {
	Count   int
	Todos   todo.TodoList
	Tickers SymbolActivity
}

func allTickers() []string {
	return []string{"AAAA", "BBBB", "CCCC"}
}

func main() {
	var data PageData = PageData{
		0,
		todo.TodoList{},
		SymbolActivity{
			Active:   []string{},
			Inactive: allTickers(),
		},
	}

	var chat chat.Chat = make(chat.Chat)

	var tickers = ticker.CreateTickerListFromArray(allTickers())
	tickers.StartAllTickers()
	defer tickers.StopAllTickers()

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
	http.HandleFunc("/ticker/add", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		symbol_len := len(r.PostForm)
		active := make([]string, symbol_len)
		if symbol_len == 0 {
			fmt.Println("No symbols to add")
			templates.ExecuteTemplate(w, "ticker-form", data.Tickers)
			return
		}
		n := 0
		for _, v := range r.PostForm {
			if len(v) == 0 || v[0] == "" {
				println("Empty value")
				continue
			}
			fmt.Println("Value:", v)
			active[n] = v[0]
			n++
		}
		active = active[:n]
		slices.Sort(active)
		activity := SymbolActivity{
			Active: active,
			Inactive: utils.Filter(allTickers(), func(e string) bool {
				return !slices.Contains(active, e)
			}),
		}
		templates.ExecuteTemplate(w, "ticker-form", activity)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index", data)
	})

	fmt.Println("Hello, World!")
	http.ListenAndServe(":8080", nil)
}
