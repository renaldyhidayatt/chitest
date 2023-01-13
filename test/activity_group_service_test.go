package test

import (
	"chigitaction/models"
	"chigitaction/payload/request"
	"chigitaction/repository"
	"chigitaction/services"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomActivityService(t *testing.T) models.Activity {
	repository := repository.NewRepositoryActivity(ConnTest)
	service := services.NewActivityService(repository)

	data := request.ActivityRequest{
		Title: "dotitle",
		Email: "doemail",
	}

	// Create
	newActivity, err := service.Create(data)
	if err != nil {
		panic(err)
	}

	// Test pass
	require.Equal(t, data.Title, newActivity.Title)
	require.Equal(t, data.Email, newActivity.Email)
	require.NotEmpty(t, newActivity.ID)
	require.NotEmpty(t, newActivity.CreatedAt)
	require.NotEmpty(t, newActivity.UpdatedAt)
	require.Nil(t, newActivity.DeletedAt)

	return newActivity
}

func TestCreateActivityServices(t *testing.T) {
	defer DropTable()
	t.Parallel()
	createRandomActivityService(t)
}

func TestGetAllServices(t *testing.T) {
	var mutex sync.Mutex
	defer DropTable()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityService(t)
			mutex.Unlock()
		}()
	}

	t.Parallel()

	repository := repository.NewRepositoryActivity(ConnTest)
	service := services.NewActivityService(repository)

	// Get activity groups
	Activitys, err := service.GetAll()
	if err != nil {
		panic(err)
	}

	for _, data := range Activitys {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.Email)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Nil(t, data.DeletedAt)
	}

}

func TestGetOneService(t *testing.T) {
	defer DropTable()
	// Create random data
	newActivity := createRandomActivityService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivity(ConnTest)
	service := services.NewActivityService(repository)

	// Find all
	Activity, err := service.GetOne(newActivity.ID)
	if err != nil {
		panic(err)
	}

	require.Equal(t, newActivity.ID, Activity.ID)
	require.Equal(t, newActivity.Title, Activity.Title)
	require.Equal(t, newActivity.Email, Activity.Email)
	require.NotEmpty(t, Activity.CreatedAt)
	require.NotEmpty(t, Activity.UpdatedAt)
	require.Nil(t, Activity.DeletedAt)
}

func TestUpdateActivityService(t *testing.T) {
	defer DropTable()
	// Create random data
	newActivity := createRandomActivityService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivity(ConnTest)
	service := services.NewActivityService(repository)

	dataUpdated := request.ActivityUpdateRequest{
		Title: "dotass",
	}

	t.Run("Update success", func(t *testing.T) {

		updatedActivity, err := service.Update(newActivity.ID, dataUpdated)

		if err != nil {
			panic(err)
		}

		require.Equal(t, newActivity.ID, updatedActivity.ID)
		require.Equal(t, newActivity.Email, updatedActivity.Email)

		require.NotEqual(t, newActivity.Title, updatedActivity.Title)
		require.NotEqual(t, newActivity.UpdatedAt, updatedActivity.UpdatedAt)

		require.NotEmpty(t, updatedActivity.CreatedAt)
		require.Nil(t, updatedActivity.DeletedAt)

	})

	t.Run("Update failed activity group not found", func(t *testing.T) {
		_, err := service.Update(7329323, dataUpdated)
		require.Error(t, err)

		message := fmt.Sprintf("Activity with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}

func TestDeleteActivityService(t *testing.T) {
	defer DropTable()
	// Create random data
	newActivity := createRandomActivityService(t)

	t.Parallel()
	repository := repository.NewRepositoryActivity(ConnTest)
	service := services.NewActivityService(repository)

	t.Run("Delete success", func(t *testing.T) {

		ok, err := service.Delete(newActivity.ID)
		if err != nil {
			panic(err)
		}

		require.True(t, ok)

	})

	t.Run("Delete failed activity group not found", func(t *testing.T) {
		ok, err := service.Delete(7329323)
		require.Error(t, err)
		require.False(t, ok)

		message := fmt.Sprintf("Activity with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}
