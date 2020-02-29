package business

import (
	"gola/app/common/def"
	"sync"
	"time"
)

var todos = map[int64]*def.Todo{}
var mux = sync.RWMutex{}
var count int64

// GetTodo 取Todo
func GetTodo(id int64) *def.Todo {
	mux.RLock()
	defer mux.RUnlock()

	todo, ok := todos[id]
	if ok {
		return todo
	}
	return nil
}

// GetTodos 取Todos
func GetTodos(done *bool) []*def.Todo {
	mux.RLock()
	defer mux.RUnlock()

	tmp := []*def.Todo{}
	for i := range todos {
		todo := todos[i]
		if done != nil {
			if *done && !todo.Done {
				continue
			}
			if !*done && todo.Done {
				continue
			}
		}
		tmp = append(tmp, todo)
	}
	return tmp
}

// AddTodo 新增Todo
func AddTodo(text string, done bool, expiredAt *time.Time, labels map[string]interface{}) (todo *def.Todo) {
	mux.Lock()
	defer mux.Unlock()

	count++
	id := count
	todo = &def.Todo{
		ID:        id,
		Text:      text,
		Done:      done,
		ExpiredAt: expiredAt,
		Labels:    labels,
	}
	todos[id] = todo
	return todo
}

// ToggleTodoDone 更改狀態
func ToggleTodoDone(id int64) *def.Todo {
	mux.Lock()
	defer mux.Unlock()

	t := GetTodo(id)
	if t != nil {
		t.Done = !t.Done
	}
	return t
}

// RemoveTodo 刪除Todo
func RemoveTodo(id int64) *def.Todo {
	mux.Lock()
	defer mux.Unlock()

	todo, ok := todos[id]
	if ok {
		delete(todos, id)
	}

	return todo
}
