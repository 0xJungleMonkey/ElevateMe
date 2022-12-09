package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

// reading the HTML files
func init() {
	tmpl = template.Must(template.ParseFiles("login.html"))

}

// customer login credential
type credentialInfo struct {
	Username string
	Password string
}

// Collect input info and show on html
func StudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	log_in := credentialInfo{

		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	tmpl.Execute(w, struct {
		Success bool
		Log_in  credentialInfo
	}{true, log_in})
}
func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	// mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/", StudentHandler)
	http.ListenAndServe(":8080", mux)
}

func Gin() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
