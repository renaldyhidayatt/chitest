package repository

import (
	"chigitaction/interfaces"
	"chigitaction/models"

	"gorm.io/gorm"
)

type ActiveRepository = interfaces.IActivityRepository

type activityRepository struct {
	db *gorm.DB
}

func NewRepositoryActivity(db *gorm.DB) *activityRepository {
	return &activityRepository{db: db}

}

func (r *activityRepository) FindAll() ([]models.Activity, error) {
	var Activitys []models.Activity

	err := r.db.Find(&Activitys).Error
	if err != nil {
		return Activitys, nil
	}

	return Activitys, nil
}

func (r *activityRepository) FindOne(id uint64) (models.Activity, error) {
	var Activity models.Activity

	err := r.db.Where("id = ?", id).Find(&Activity).Error
	if err != nil {
		return Activity, nil
	}

	return Activity, nil
}

func (r *activityRepository) Save(Activity models.Activity) (models.Activity, error) {
	err := r.db.Create(&Activity).Error
	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (r *activityRepository) Update(Activity models.Activity) (models.Activity, error) {
	err := r.db.Save(&Activity).Error
	if err != nil {
		return Activity, err
	}

	return Activity, nil
}

func (r *activityRepository) Delete(Activity models.Activity) (bool, error) {
	err := r.db.Delete(&Activity).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
