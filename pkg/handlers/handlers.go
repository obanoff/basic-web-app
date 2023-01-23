package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/obanoff/basic-web-app/pkg/config"
	"github.com/obanoff/basic-web-app/pkg/models"
	"github.com/obanoff/basic-web-app/pkg/render"
)

// Repo the repository used by the handlers
var Repo = &Repository{}

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
// func NewRepo(a *config.AppConfig) {
// 	Repo = &Repository{
// 		App: a,
// 	}
// }

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{}, r)
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{}, r)
}

// Reservation is the handler for the make-reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{}, r)
}

// Generals is the handler for the genenerals-quoters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.gohtml", &models.TemplateData{}, r)
}

// Majors is the handler for the majors-suite page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.gohtml", &models.TemplateData{}, r)
}

// Availability is the handler for the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.gohtml", &models.TemplateData{}, r)
}

// PostAvailability is the handler for the search-availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// retrieve values from form inputs
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("start date is %s, and end date is %s", start, end)))
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.gohtml", &models.TemplateData{}, r)
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles requests for availability and sends back JSON repsonse
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
}
