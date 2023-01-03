package test

import (
	"chigitaction/models"
	"chigitaction/repository"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTodoRepository(t *testing.T) models.Todo {
	todoRepository := repository.NewTodoRepository(ConnTest)

	newActivity := createRandomActivityRepository(t)

	todo := models.Todo{
		ActivityGroupID: newActivity.ID,
		Title:           "jabufaker.RandomString(20)",
		IsActive:        true,
	}

	// Save to db
	newTodo, err := todoRepository.Save(todo)
	if err != nil {
		panic(err.Error())
	}

	// Test pas
	require.NoError(t, err)

	require.NotEmpty(t, newTodo.ID)
	require.NotEmpty(t, newActivity.CreatedAt)
	require.NotEmpty(t, newActivity.UpdatedAt)
	require.Empty(t, newActivity.DeletedAt)

	require.Equal(t, todo.ActivityGroupID, newTodo.ActivityGroupID)
	require.Equal(t, todo.Title, newTodo.Title)
	require.Equal(t, todo.IsActive, newTodo.IsActive)

	return newTodo
}

func TestFindAllTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomTodoRepository(t)
			mutex.Unlock()
		}()
	}

	todoRepository := repository.NewTodoRepository(ConnTest)

	// Find all
	todos, err := todoRepository.FindAll()
	if err != nil {
		panic(err.Error())
	}

	require.NoError(t, err)
	require.NotEqual(t, 0, len(todos))

	for _, data := range todos {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.ActivityGroupID)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)

		require.NotNil(t, data.IsActive)

		require.Empty(t, data.DeletedAt)
	}

}

func TestFindByActivityTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// todos for store `new data todos` from channel
	var todos []models.Todo

	// channel for store data `new data todos` from process create random todo
	channel := make(chan models.Todo)
	defer close(channel)

	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodos := createRandomTodoRepository(t)
			channel <- newTodos
			mutex.Unlock()
		}()
		todos = append(todos, <-channel)
	}

	todoRepository := repository.NewTodoRepository(ConnTest)

	// Find by actiivity group
	todos, err := todoRepository.FindByActivityID(todos[0].ActivityGroupID)
	if err != nil {
		panic(err.Error())
	}

	require.NoError(t, err)
	require.Equal(t, 1, len(todos))

	for _, data := range todos {
		// Equal
		require.Equal(t, todos[0].ID, data.ID)
		require.Equal(t, todos[0].Title, data.Title)
		require.Equal(t, todos[0].ActivityGroupID, data.ActivityGroupID)
		require.Equal(t, todos[0].IsActive, data.IsActive)

		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Empty(t, data.DeletedAt)
	}

}

func TestFindOneTodoRepository(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoRepository(t)

	todoRepository := repository.NewTodoRepository(ConnTest)

	// Find One
	todo, err := todoRepository.FindOne(newTodo.ID)
	if err != nil {
		panic(err.Error())
	}

	require.NoError(t, err)

	// Equal
	require.Equal(t, todo.ID, newTodo.ID)
	require.Equal(t, todo.Title, newTodo.Title)
	require.Equal(t, todo.ActivityGroupID, newTodo.ActivityGroupID)
	require.Equal(t, todo.IsActive, newTodo.IsActive)

	require.NotEmpty(t, newTodo.CreatedAt)
	require.NotEmpty(t, newTodo.UpdatedAt)
	require.Empty(t, newTodo.DeletedAt)

}

func TestCreateTodoRepository(t *testing.T) {
	t.Parallel()
	createRandomTodoRepository(t)
}

func TestUpdateTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// todos for store `new data todos` from channel
	var todos []models.Todo

	// channel for store data `new data todos` from process create random todo
	channel := make(chan models.Todo)
	defer close(channel)

	// Create some random data
	for i := 0; i < 2; i++ {
		go func() {
			mutex.Lock()
			newTodos := createRandomTodoRepository(t)
			channel <- newTodos
			mutex.Unlock()
		}()
		todos = append(todos, <-channel)
	}

	todoRepository := repository.NewTodoRepository(ConnTest)

	dataUpdate := models.Todo{
		ID:              todos[0].ID,
		ActivityGroupID: todos[1].ActivityGroupID,
		Title:           "jabufaker.RandomString(20)",
		IsActive:        false,
		CreatedAt:       todos[0].CreatedAt,
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	// Update
	todo, err := todoRepository.Update(dataUpdate)

	if err != nil {
		panic(err.Error())
	}

	require.NoError(t, err)

	// Test
	require.Equal(t, todo.ID, todos[0].ID)

	require.NotEqual(t, todo.Title, todos[0].Title)
	require.NotEqual(t, todo.ActivityGroupID, todos[0].ActivityGroupID)
	require.NotEqual(t, todo.IsActive, todos[0].IsActive)

	require.NotEmpty(t, todo.CreatedAt)
	require.NotEmpty(t, todo.UpdatedAt)

	require.Empty(t, todos[0].DeletedAt)
}

func TestDeleteTodoRepository(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoRepository(t)

	todoRepository := repository.NewTodoRepository(ConnTest)

	// Update
	ok, err := todoRepository.Delete(newTodo)
	if err != nil {
		panic(err.Error())
	}

	require.NoError(t, err)
	require.True(t, ok)

	todo, err := todoRepository.FindOne(newTodo.ID)
	if err != nil {
		panic(err.Error())
	}
	nullId := uint64(0)
	require.Equal(t, nullId, todo.ID)
}
