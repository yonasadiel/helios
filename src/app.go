package helios

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // use sqlite dialect
)

// Helios is the core of the apps
type Helios struct {
	models []interface{}
	DB     *gorm.DB
	DBErr  error
	store  *sessions.CookieStore
}

// App will be the core app that has all the models
// and be the core of the server
var App Helios

// Initialize the database to production database
func (app *Helios) Initialize() {
	app.DB, app.DBErr = gorm.Open("sqlite3", "db.sqlite3")
	if app.DBErr != nil {
		panic(app.DBErr)
	}
	key := []byte(os.Getenv("HELIOS_SECRET"))
	app.store = sessions.NewCookieStore(key)
}

// RegisterModel so the database will be migrated
func (app *Helios) RegisterModel(model interface{}) {
	app.models = append(app.models, model)
}

// CloseDB close the database connection
func (app *Helios) CloseDB() {
	app.DB.Close()
}

// Migrate migrate all the models
func (app *Helios) Migrate() {
	for _, model := range app.models {
		app.DB.AutoMigrate(model)
	}
}

// BeforeTest has to be called everytime a test is run
// It will reset the database
func (app *Helios) BeforeTest() {
	if app.DB == nil {
		app.DB, app.DBErr = gorm.Open("sqlite3", ":memory:")
		if app.DBErr != nil {
			panic(app.DBErr)
		}
		app.Migrate()
	} else {
		for _, model := range app.models {
			app.DB.Unscoped().Delete(model, "true")
		}
	}
}

func (app *Helios) getSession(r *http.Request) *sessions.Session {
	session, _ := app.store.Get(r, os.Getenv("SESSION_NAME"))
	return session
}
