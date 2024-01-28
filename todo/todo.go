package todo

import (
	"html/template"
	"htmx/utils"
	"net/http"
	"strconv"
)

type TodoItem struct {
	Id        int
	Name      string
	Completed bool
}

type Toggle struct {
	Id    int
	Value bool
}

type TodoList []TodoItem

func (todoList *TodoList) CreateTodoAdder(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		*todoList = todoList.addTodo(name)
		println("Length:", len(*todoList))
		templates.ExecuteTemplate(w, "todos", todoList)
	}
}

func (todoList *TodoList) CreateTodoDeleter(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		id_no, err := strconv.Atoi(id)
		utils.HandleError(err)
		if id_no >= len(*todoList) {
			println("Index out of bounds", id_no, "for slice of length", len(*todoList), ".")
			return
		}
		*todoList = append((*todoList)[:id_no], (*todoList)[id_no+1:]...)
		templates.ExecuteTemplate(w, "todos", todoList)
	}
}

func (todoList *TodoList) CreateTodoToggler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		value := r.URL.Query().Get("value")

		id_no, err := strconv.Atoi(id)
		utils.HandleError(err)
		value_bool, err := strconv.ParseBool(value)
		utils.HandleError(err)
		(*todoList)[id_no].Completed = value_bool
		templates.ExecuteTemplate(w, "toggle", Toggle{id_no, value_bool})
	}

}

func (todoList *TodoList) addTodo(name string) TodoList {
	id := len(*todoList)
	println(id)
	return append(*todoList, *newTodoItem(id, name))
}

func newTodoItem(id int, name string) *TodoItem {
	return &TodoItem{
		Id:        id,
		Name:      name,
		Completed: false,
	}
}
