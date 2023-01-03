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

type ActivityHandler = interfaces.IActivityHandler

type activityHandler struct {
	service services.ActivityService
}

func NewActivityHandler(service services.ActivityService) *activityHandler {
	return &activityHandler{service: service}
}

func (h *activityHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	activity, err := h.service.GetAll()

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	} else {
		response.ResponseMessage(w, "berhasil mendapatkan data", activity, http.StatusOK)
	}

}

func (h *activityHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	activityId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	todo, err := h.service.GetOne(uint64(activityId))

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	} else {
		response.ResponseMessage(w, "berhasil mendapatkan data", todo, http.StatusOK)
	}
}

func (h *activityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var activity request.ActivityRequest
	err := json.NewDecoder(r.Body).Decode(&activity)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	res, err := h.service.Create(activity)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return

	}

	response.ResponseMessage(w, "berhasil membuat Activity", res, http.StatusOK)
}

func (h *activityHandler) Update(w http.ResponseWriter, r *http.Request) {
	activityIDStr := r.URL.Query().Get("id")
	activityId, err := strconv.ParseInt(activityIDStr, 10, 64)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var activity request.ActivityUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&activity)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	res, err := h.service.Update(uint64(activityId), activity)

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	response.ResponseMessage(w, "Berhasil mengubah data", res, http.StatusOK)
}

func (h *activityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	activityIDStr := r.URL.Query().Get("id")
	activityId, err := strconv.ParseInt(activityIDStr, 10, 64)

	if err != nil {
		response.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	activity, err := h.service.Delete(uint64(activityId))

	if err != nil {
		response.ResponseError(w, http.StatusInternalServerError, err)
		return
	} else {
		response.ResponseMessage(w, "berhasil mendapatkan data", activity, http.StatusOK)
	}
}
