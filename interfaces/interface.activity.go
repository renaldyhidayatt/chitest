package interfaces

import (
	"chigitaction/models"
	"chigitaction/payload/request"
	"net/http"
)

type IActivityRepository interface {
	Save(Activity models.Activity) (models.Activity, error)
	FindAll() ([]models.Activity, error)
	FindOne(id uint64) (models.Activity, error)
	Update(Activity models.Activity) (models.Activity, error)
	Delete(Activity models.Activity) (bool, error)
}

type IActivityService interface {
	Create(req request.ActivityRequest) (models.Activity, error)
	GetAll() ([]models.Activity, error)
	GetOne(id uint64) (models.Activity, error)
	Update(id uint64, req request.ActivityUpdateRequest) (models.Activity, error)
	Delete(id uint64) (bool, error)
}

type IActivityHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
