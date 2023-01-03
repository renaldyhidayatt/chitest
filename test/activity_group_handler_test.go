package test

import (
	"chigitaction/payload/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomActivityHandler(t *testing.T) request.ActivityCreateResponse {
	data := request.ActivityRequest{
		Title: "jabufaker.RandomString(20)",
		Email: "jabufaker.RandomEmail()",
	}

	dataBody := fmt.Sprintf(`{"title": "%s", "email": "%s"}`, data.Title, data.Email)
	requestBody := strings.NewReader(dataBody)

	requestt := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/activity-groups/", requestBody)
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
	require.Equal(t, data.Email, contextData["email"])

	newActivity := request.ActivityCreateResponse{
		ID:    uint64(contextData["id"].(float64)),
		Title: contextData["title"].(string),
		Email: contextData["email"].(string),
	}

	return newActivity
}

func TestActivityCreateHandler(t *testing.T) {
	t.Parallel()
	t.Run("Handler Create activity Group success", func(t *testing.T) {
		createRandomActivityHandler(t)
	})

	t.Run("Handler Create activity Group failed title blank", func(t *testing.T) {
		dataBody := fmt.Sprintf(`{"title": "%s", "email": "%s"}`, "", "jabufaker.RandomEmail()")
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:5000/api/activity-groups/", requestBody)
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
}

func TestGetAllActivityHandler(t *testing.T) {
	var mutex sync.Mutex
	t.Parallel()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityHandler(t)
			mutex.Unlock()
		}()
	}

	request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/activity-groups/", nil)
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

	var contextData = responseBody["data"].([]interface{})

	require.NotEqual(t, 0, len(contextData))
	// Data is not null

	for _, data := range contextData {
		list := data.(map[string]interface{})
		require.NotEmpty(t, list["id"])
		require.NotEmpty(t, list["title"])
		require.NotEmpty(t, list["email"])
		require.NotEmpty(t, list["created_at"])
		require.NotEmpty(t, list["updated_at"])
	}
}

func TestGetOneActivity(t *testing.T) {
	t.Parallel()
	newActivity := createRandomActivityHandler(t)

	t.Run("Success get one", func(t *testing.T) {
		id := fmt.Sprintf("%d", newActivity.ID)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/activity-groups/"+id, nil)
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
		require.Equal(t, newActivity.ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newActivity.Title, contextData["title"])
		require.Equal(t, newActivity.Email, contextData["email"])

		require.NotEmpty(t, contextData["created_at"])
		require.NotEmpty(t, contextData["updated_at"])
		require.Nil(t, contextData["deteled_at"])
	})

	t.Run("Id not found", func(t *testing.T) {
		wrongId := "999999"
		request := httptest.NewRequest(http.MethodGet, "http://localhost:5000/api/activity-groups/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}

func TestUpdateActivity(t *testing.T) {
	t.Parallel()
	newActivity := createRandomActivityHandler(t)

	t.Run("Success updated activity group", func(t *testing.T) {
		data := request.ActivityUpdateRequest{
			Title: "jabufaker.RandomString(20)",
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newActivity.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/activity-groups/"+id, requestBody)
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
		require.Equal(t, newActivity.ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newActivity.Email, contextData["email"])

		require.NotEqual(t, newActivity.UpdatedAt.String(), contextData["updated_at"])
		require.NotEqual(t, newActivity.Title, contextData["title"])

		require.NotEmpty(t, contextData["created_at"])

	})

	t.Run("Body blank", func(t *testing.T) {
		dataBody := fmt.Sprintf(`{"title": "%s"}`, "")
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newActivity.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/activity-groups/"+id, requestBody)
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

	t.Run("Id not found", func(t *testing.T) {
		data := request.ActivityUpdateRequest{
			Title: "jabufaker.RandomString(20)",
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		wrongId := "999999"
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:5000/api/activity-groups/"+wrongId, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}

func TestDeleteActivity(t *testing.T) {
	t.Parallel()
	newActivity := createRandomActivityHandler(t)

	t.Run("Deleted success", func(t *testing.T) {
		id := fmt.Sprintf("%d", newActivity.ID)
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/activity-groups/"+id, nil)
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
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:5000/api/activity-groups/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}
