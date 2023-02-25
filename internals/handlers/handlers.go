package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/obanoff/basic-web-app/internals/config"
	"github.com/obanoff/basic-web-app/internals/forms"
	"github.com/obanoff/basic-web-app/internals/helpers"
	"github.com/obanoff/basic-web-app/internals/models"
	"github.com/obanoff/basic-web-app/internals/render"
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// err = errors.New("this is an error message")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.FormValue("phone_number"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required(
		"first_name",
		"last_name",
		"email",
	)
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)

		return
	}

	// put data into session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// redirection
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
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
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles requests for availability and sends back JSON repsonse
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	data := r.Form
	resp := jsonResponse{
		StartDate: data.Get("start_date"),
		EndDate:   data.Get("end_date"),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//get data from session and cast it to model.Reservation type
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		log.Println("cannot get item from session")
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	}, r)
}
