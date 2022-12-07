package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

// reading the HTML files
func init() {
	tmpl = template.Must(template.ParseFiles("StudentInfo.html"))

}

type studentInfo struct {
	Sid    string
	Name   string
	Course string
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	student := studentInfo{
		Sid:    r.FormValue("sid"),
		Name:   r.FormValue("name"),
		Course: r.FormValue("course"),
	}
	tmpl.Execute(w, struct {
		Success bool
		Student studentInfo
	}{true, student})
}
func main() {
	http.HandleFunc("/", StudentHandler)
	http.ListenAndServe(":8080", nil)
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
