package test

import (
	"chigitaction/models"
	"chigitaction/repository"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomActivityRepository(t *testing.T) models.Activity {
	ActivityRepository := repository.NewRepositoryActivity(ConnTest)

	Activity := models.Activity{
		Title: "dotasss",
		Email: "jabufaker.RandomEmail()",
	}

	// Save to db
	newActivity, err := ActivityRepository.Save(Activity)
	if err != nil {
		panic(err)
	}

	// Test pass
	require.Equal(t, Activity.Title, newActivity.Title)
	require.Equal(t, Activity.Email, newActivity.Email)
	require.NotEmpty(t, newActivity.ID)
	require.NotEmpty(t, newActivity.CreatedAt)
	require.NotEmpty(t, newActivity.UpdatedAt)
	require.Empty(t, newActivity.DeletedAt)

	return newActivity
}

func TestCreateActivity(t *testing.T) {
	defer DropTable()
	t.Parallel()
	createRandomActivityRepository(t)
}

func TestFindAllActivity(t *testing.T) {
	var mutex sync.Mutex
	defer DropTable()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityRepository(t)
			mutex.Unlock()
		}()
	}

	t.Parallel()
	ActivityRepository := repository.NewRepositoryActivity(ConnTest)

	// Find all
	Activity, err := ActivityRepository.FindAll()
	if err != nil {
		panic(err)
	}

	for _, data := range Activity {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.Email)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Empty(t, data.DeletedAt)
	}
}

func TestFindOneActivity(t *testing.T) {
	defer DropTable()
	// Create random data
	newActivity := createRandomActivityRepository(t)

	t.Parallel()
	ActivityRepository := repository.NewRepositoryActivity(ConnTest)

	// Find all
	Activity, err := ActivityRepository.FindOne(newActivity.ID)
	if err != nil {
		panic(err)
	}

	require.Equal(t, newActivity.ID, Activity.ID)
	require.Equal(t, newActivity.Title, Activity.Title)
	require.Equal(t, newActivity.Email, Activity.Email)
	require.NotEmpty(t, Activity.CreatedAt)
	require.NotEmpty(t, Activity.UpdatedAt)
	require.Empty(t, Activity.DeletedAt)
}

func TestUpdateActivityRepository(t *testing.T) {
	defer DropTable()
	newActivity := createRandomActivityRepository(t)
	t.Parallel()
	ActivityRepository := repository.NewRepositoryActivity(ConnTest)

	dataUpdate := models.Activity{
		ID:        newActivity.ID,
		Title:     "jabufaker.RandomString(20",
		Email:     "jabufaker.RandomEmail()",
		CreatedAt: newActivity.CreatedAt,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// update
	updateActivity, err := ActivityRepository.Update(dataUpdate)
	if err != nil {
		panic(err)
	}

	require.Equal(t, newActivity.ID, updateActivity.ID)
	require.Equal(t, newActivity.CreatedAt, updateActivity.CreatedAt)
	// require.Equal(t, newActivity.DeletedAt, updateActivity.DeletedAt)
	require.NotEqual(t, newActivity.Title, updateActivity.Title)
	require.NotEqual(t, newActivity.Email, updateActivity.Email)
}

func TestDeleteActivityRepository(t *testing.T) {
	DropTable()
	newActivity := createRandomActivityRepository(t)
	t.Parallel()

	ActivityRepository := repository.NewRepositoryActivity(ConnTest)

	ok, err := ActivityRepository.Delete(newActivity)
	if err != nil {
		panic(err)
	}
	require.True(t, ok)

	Activity, err := ActivityRepository.FindOne(newActivity.ID)
	if err != nil {
		panic(err)
	}
	require.Equal(t, 0, int(Activity.ID))
}
