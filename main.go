package main

import (
	_ "database/sql"
	"fmt"
	"github.com/Dubrovsky18/tinkoff_cup/FileUpload"
	"github.com/Dubrovsky18/tinkoff_cup/Login"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)

}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/login", Login.LoginHandlePost).Methods("POST")
	router.HandleFunc("/login", Login.LoginHandleGet).Methods("GET")
	router.HandleFunc("/registration", Login.RegistrationHandlerGet).Methods("GET")
	router.HandleFunc("/registration", Login.RegistrationHandlerPost).Methods("POST")
	router.HandleFunc("/main", FileUpload.HandleMain).Methods("GET")
	router.HandleFunc("/main", FileUpload.FileUpload).Methods("POST")

	//http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":9090", nil)
}

func main() {
	handleFunc()
}
