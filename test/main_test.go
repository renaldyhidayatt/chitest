package test

import (
	"chigitaction/config"
	"chigitaction/router"
	"chigitaction/utils"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var ConnTest *gorm.DB
var Route *chi.Mux

func DropTable() {
	ConnTest.Raw("delete from todos")
}

func TestMain(m *testing.M) {

	if _, ok := os.LookupEnv("GO_ENV"); !ok {
		err := utils.Viper()
		if err != nil {
			log.Fatalf(".env file not load: %v", err)
		}
	}
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "chitesting")

	db, _ := config.Database()

	ConnTest = db

	Route = chi.NewRouter()

	router.NewTodoRoutes("/api/todo", ConnTest, Route).TodoRoute()
	router.NewActivityRoutes("/api/activity-groups", ConnTest, Route).ActivityRoute()

	http.ListenAndServe(":5000", Route)
}
