package services

import (
	"chigitaction/interfaces"
	"chigitaction/models"
	"chigitaction/payload/request"
	"chigitaction/repository"
	"errors"
	"fmt"
	"time"
)

type ActivityService = interfaces.IActivityService

type activityService struct {
	repository repository.ActiveRepository
}

func NewActivityService(repository repository.ActiveRepository) *activityService {
	return &activityService{repository: repository}
}

func (s *activityService) GetAll() ([]models.Activity, error) {
	// Find all
	Activitys, err := s.repository.FindAll()

	if err != nil {
		return Activitys, err
	}

	return Activitys, nil
}

func (s *activityService) GetOne(id uint64) (models.Activity, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)

	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (s *activityService) Create(req request.ActivityRequest) (models.Activity, error) {
	Activity := models.Activity{
		Title: req.Title,
		Email: req.Email,
	}

	// Save
	newActivity, err := s.repository.Save(Activity)
	if err != nil {
		return newActivity, err
	}

	return newActivity, nil
}

func (s *activityService) Update(id uint64, req request.ActivityUpdateRequest) (models.Activity, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)
	// If activity group not found
	if Activity.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return Activity, errors.New(message)
	}

	if err != nil {
		return Activity, err
	}

	// Change field title to req update title
	Activity.Title = req.Title
	// Change time field updatUpdatedAted
	Activity.UpdatedAt = time.Now()

	// Update
	updatedActivity, err := s.repository.Update(Activity)
	if err != nil {
		return Activity, err
	}

	return updatedActivity, nil
}

func (s *activityService) Delete(id uint64) (bool, error) {
	// Find one
	Activity, err := s.repository.FindOne(id)
	// If activity group not found
	if Activity.ID == 0 {
		message := fmt.Sprintf("Activity with ID %d Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	ok, err := s.repository.Delete(Activity)
	if err != nil {
		return false, err
	}

	return ok, nil
}
