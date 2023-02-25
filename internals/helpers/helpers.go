package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/obanoff/basic-web-app/internals/config"
)

var App *config.AppConfig

func ClientError(w http.ResponseWriter, status int) {
	App.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	App.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
