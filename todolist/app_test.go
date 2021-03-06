package todolist

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	year := strconv.Itoa(time.Now().Year())

	app.AddTodo("a do some stuff due may 23")

	todo := app.TodoList.FindById(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal(fmt.Sprintf("%s-05-23", year), todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}

func TestAddTodoWithEuropeanDates(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}

	app.AddTodo("a do some stuff due 23 may")

	todo := app.TodoList.FindById(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal("2017-05-23", todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}

func TestListbyProject(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	app.Load()

	// create three todos w/wo a project
	app.AddTodo("this is a test +testme")
	app.AddTodo("this is a test +testmetoo @work")
	app.AddTodo("this is a test with no projects")
	app.CompleteTodo("c 1")

	// simulate listTodos
	input := "l by p"
	filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	grouped := app.getGroups(input, filtered)

	assert.Equal(3, len(grouped.Groups))

	// testme project has 1 todo and its completed
	assert.Equal(1, len(grouped.Groups["testme"]))
	assert.Equal(true, grouped.Groups["testme"][0].Completed)

	// testmetoo project has 1 todo and it has a context
	assert.Equal(1, len(grouped.Groups["testmetoo"]))
	assert.Equal(1, len(grouped.Groups["testmetoo"][0].Contexts))
	assert.Equal("work", grouped.Groups["testmetoo"][0].Contexts[0])
}

func TestListbyContext(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	app.Load()

	// create three todos w/wo a context
	app.AddTodo("this is a test +testme")
	app.AddTodo("this is a test +testmetoo @work")
	app.AddTodo("this is a test with no projects")
	app.CompleteTodo("c 1")

	// simulate listTodos
	input := "l by c"
	filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	grouped := app.getGroups(input, filtered)

	assert.Equal(2, len(grouped.Groups))

	// work context has 1 todo and it has a project of testmetoo
	assert.Equal(1, len(grouped.Groups["work"]))
	assert.Equal(1, len(grouped.Groups["work"][0].Projects))
	assert.Equal("testmetoo", grouped.Groups["work"][0].Projects[0])

	// There are two todos with no context
	assert.Equal(2, len(grouped.Groups["No contexts"]))

	// check to see if the a todos with no context contain a
	// completed todo
	var hasACompletedTodo bool
	for _, todo := range grouped.Groups["No contexts"] {
		if todo.Completed {
			hasACompletedTodo = true
		}
	}
	assert.Equal(true, hasACompletedTodo)
}
