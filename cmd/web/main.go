package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/obanoff/basic-web-app/internals/config"
	"github.com/obanoff/basic-web-app/internals/driver"
	"github.com/obanoff/basic-web-app/internals/handlers"
	"github.com/obanoff/basic-web-app/internals/helpers"
	"github.com/obanoff/basic-web-app/internals/models"
	"github.com/obanoff/basic-web-app/internals/render"
	"github.com/obanoff/basic-web-app/internals/repository/dbrepo"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Staring application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what I'm going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=eugene password=")
	if err != nil {
		log.Fatal("Cannot connect to DB! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	// ONLY for development mode: rebuilds cache on every request instead of using already generated cache (usefull when making changes to templates and don't want to restart the server every time)
	app.UseCache = false

	// handlers.NewRepo(&app)
	// handlers.Repo.App = &app
	// handlers.Repo.DB = dbrepo.NewPostgresRepo(db.SQL, &app)
	handlers.Repo = &handlers.Repository{
		App: &app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, &app),
	}

	helpers.App = &app

	render.NewRenderer(&app)

	return db, nil
}
