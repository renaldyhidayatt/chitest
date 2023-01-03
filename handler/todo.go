package handler

import (
	"chigitaction/interfaces"
	"chigitaction/payload/request"
	"chigitaction/payload/response"
	"chigitaction/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type TodoHandler = interfaces.ITodoHandler

type todoHandler struct {
	service services.TodoService
}

func NewTodoHandler(service services.TodoService) *todoHandler {
	return &todoHandler{service: service}
}

func (h *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var todo request.TodoCreateRequest
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	newTodo, err := h.service.Create(todo)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return

	}

	response.ResponseMessage(w, "berhasil membuat todo", newTodo, http.StatusOK)
}

func (h *todoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	todo, err := h.service.GetAll(uint64(todoID))

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)

	} else {
		response.ResponseMessage(w, "Berhasil mendapatkan data", todo, http.StatusOK)
	}

}

func (h *todoHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	todo, err := h.service.GetOne(uint64(todoID))

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	} else {
		response.ResponseMessage(w, "berhasil mendapatkan data", todo, http.StatusOK)
	}

}

func (h *todoHandler) Update(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	todoID, _ := strconv.ParseInt(todoIDStr, 10, 64)

	var todo request.TodoUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	res, err := h.service.Update(uint64(todoID), todo)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	response.ResponseMessage(w, "Berhasil mengubah data", res, http.StatusOK)
}

func (h *todoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	todo, err := h.service.Delete(uint64(todoID))

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	} else {
		response.ResponseMessage(w, "berhasil mendapatkan data", todo, http.StatusOK)
	}
}
