package test

import (
	"chigitaction/payload/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTodoHandler(t *testing.T) request.TodoResponse {
	newActivityGroup := createRandomActivityHandler(t)
	data := request.TodoResponse{
		Title:      "jabufaker.RandomString(20)",
		ActivityID: newActivityGroup.ID,
	}

	dataBody := fmt.Sprintf(`{"title": "%s", "activity_group_id": %d}`, data.Title, data.ActivityID)
	requestBody := strings.NewReader(dataBody)

	requestt := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
	requestt.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	Route.ServeHTTP(recorder, requestt)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	require.Equal(t, 201, response.StatusCode)
	require.Equal(t, "Success", responseBody["status"])
	require.Equal(t, "Success", responseBody["message"])

	require.NotEmpty(t, responseBody["data"])

	var contextData = responseBody["data"].(map[string]interface{})
	require.NotEmpty(t, contextData["id"])
	require.NotEmpty(t, contextData["created_at"])
	require.NotEmpty(t, contextData["updated_at"])

	require.Equal(t, data.Title, contextData["title"])
	require.Equal(t, data.ActivityID, uint64(contextData["activity_group_id"].(float64)))

	require.Equal(t, true, contextData["is_active"])

	isActiveString := strconv.FormatBool(contextData["is_active"].(bool))
	newtodo := request.TodoResponse{
		ID:         uint64(contextData["id"].(float64)),
		Title:      contextData["title"].(string),
		ActivityID: uint64(contextData["activity_group_id"].(float64)),
		IsActive:   isActiveString,
	}

	return newtodo
}

func TestCreateTodoHandler(t *testing.T) {
	t.Parallel()

	t.Run("create new todo success", func(t *testing.T) {
		createRandomTodoHandler(t)
	})

	t.Run("create new todo without title", func(t *testing.T) {
		newActivityGroup := createRandomActivityHandler(t)
		data := request.TodoResponse{
			ActivityID: newActivityGroup.ID,
		}

		dataBody := fmt.Sprintf(`{"activity_group_id": "%d"}`, data.ActivityID)
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 400, response.StatusCode)
		require.Equal(t, "Bad Request", responseBody["status"])
		require.Equal(t, "title cannot be null", responseBody["message"])

		require.Empty(t, responseBody["data"])
	})

	t.Run("create new todo withoud activity group id", func(t *testing.T) {
		data := request.TodoResponse{
			Title: "jabufaker.RandomString(20)",
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 400, response.StatusCode)
		require.Equal(t, "Bad Request", responseBody["status"])
		require.Equal(t, "activity_group_id cannot be null", responseBody["message"])

		require.Empty(t, responseBody["data"])
	})
}

func TestGetAllTodoHandler(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	var newTodos []request.TodoResponse

	// Create channel for store new todos created
	channel := make(chan request.TodoResponse)
	defer close(channel)
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodo := createRandomTodoHandler(t)
			channel <- newTodo
			mutex.Unlock()
		}()
		newTodos = append(newTodos, <-channel)
	}

	t.Run("Get all todo without query activity group id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/todo", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		require.NotEmpty(t, responseBody["data"])

		contextBody := responseBody["data"].([]interface{})
		// Length todos must be greater than 0
		require.NotEqual(t, 0, len(contextBody))

		for _, data := range contextBody {
			list := data.(map[string]interface{})
			require.NotEmpty(t, list["id"])
			require.NotEmpty(t, list["title"])
			require.NotEmpty(t, list["is_active"])
			require.NotEmpty(t, list["created_at"])
			require.NotEmpty(t, list["updated_at"])
			require.Nil(t, list["deleted_at"])
		}
	})

	t.Run("Get all todo with query activity group id", func(t *testing.T) {
		ActivityID := fmt.Sprintf("%d", newTodos[0].ActivityID)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/todo?activity_group_id="+ActivityID, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		require.NotEmpty(t, responseBody["data"])

		contextBody := responseBody["data"].([]interface{})
		// Length todos must be 1
		require.Equal(t, 1, len(contextBody))

		for _, data := range contextBody {
			list := data.(map[string]interface{})
			require.Equal(t, newTodos[0].ID, uint64(list["id"].(float64)))
			require.Equal(t, newTodos[0].Title, list["title"])

			var isActive string
			// If isActive is false
			if newTodos[0].IsActive == "false" {
				isActive = "0"
			} else {
				isActive = "1"
			}
			require.Equal(t, isActive, list["is_active"])

			require.NotEmpty(t, list["created_at"])
			require.NotEmpty(t, list["updated_at"])

			require.Nil(t, list["deleted_at"])
		}
	})

}

func TestGetOneTodoHandler(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	var newTodos []request.TodoResponse

	// Create channel for store new todos created
	channel := make(chan request.TodoResponse)
	defer close(channel)
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodo := createRandomTodoHandler(t)
			channel <- newTodo
			mutex.Unlock()
		}()
		newTodos = append(newTodos, <-channel)
	}

	t.Run("Get one todo success", func(t *testing.T) {
		todoId := fmt.Sprintf("%d", newTodos[0].ID)
		requestt := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/todo/todo-items/"+todoId, nil)
		requestt.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, requestt)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		var contextData = responseBody["data"].(map[string]interface{})
		require.NotEmpty(t, contextData)

		require.Equal(t, newTodos[0].ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newTodos[0].Title, contextData["title"])

		var isActive string
		// If isActive is false
		if newTodos[0].IsActive == "false" {
			isActive = "0"
		} else {
			isActive = "1"
		}
		require.Equal(t, isActive, contextData["is_active"])

		require.NotEmpty(t, contextData["created_at"])
		require.NotEmpty(t, contextData["updated_at"])

		require.Nil(t, contextData["deleted_at"])
	})

	t.Run("Get one ID not found", func(t *testing.T) {
		wrongID := fmt.Sprintf("%d", 9999999)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/todo/todo-items/"+wrongID, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongID)
		require.Equal(t, message, responseBody["message"])

		require.Empty(t, responseBody["data"])
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Parallel()

	t.Run("Success update todo with field title", func(t *testing.T) {
		newTodo := createRandomTodoHandler(t)
		data := request.TodoUpdateRequest{
			Title: "jabufaker.RandomString(20)",
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newTodo.ID)
		requestt := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/todo/todo-items/"+id, requestBody)
		requestt.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, requestt)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])
		require.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		require.Equal(t, newTodo.ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newTodo.ActivityID, uint64(contextData["activity_group_id"].(float64)))

		var isActive string
		// If isActive is false
		if newTodo.IsActive == "false" {
			isActive = "0"
		} else {
			isActive = "1"
		}
		require.Equal(t, isActive, contextData["is_active"])

		require.NotEmpty(t, contextData["created_at"])

		require.NotEqual(t, newTodo.UpdatedAt.String(), contextData["updated_at"])
		require.NotEqual(t, newTodo.Title, contextData["title"])

		require.Nil(t, newTodo.DeletetAt)
	})

	t.Run("Success update todo with field is_active", func(t *testing.T) {
		newTodo := createRandomTodoHandler(t)
		data := request.TodoUpdateRequest{
			IsActive: false,
		}

		dataBody := fmt.Sprintf(`{"is_active": %t}`, data.IsActive)
		requestBody := strings.NewReader(dataBody)
		id := fmt.Sprintf("%d", newTodo.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/todo/todo-items/"+id, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])
		require.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		require.Equal(t, newTodo.ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newTodo.ActivityID, uint64(contextData["activity_group_id"].(float64)))
		require.Equal(t, newTodo.Title, contextData["title"])

		require.NotEmpty(t, contextData["created_at"])

		require.NotEqual(t, newTodo.UpdatedAt.String(), contextData["updated_at"])
		require.NotEqual(t, newTodo.IsActive, contextData["is_active"])

		require.Nil(t, newTodo.DeletetAt)

	})

	t.Run("Id not found", func(t *testing.T) {
		data := request.ActivityUpdateRequest{
			Title: "jabufaker.RandomString(20)",
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		wrongId := "999999"
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/todo/"+wrongId, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoHandler(t)

	t.Run("Deleted success", func(t *testing.T) {
		id := fmt.Sprintf("%d", newTodo.ID)
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/todo/"+id, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])
		require.Empty(t, responseBody["data"])
	})

	t.Run("Id not found", func(t *testing.T) {
		wrongId := "999999"
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/todo/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}
