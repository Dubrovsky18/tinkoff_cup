package FileLoad

import (
	"fmt"
	"github.com/Dubrovsky18/tinkoff_cup/RunTester"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	Name  string
	port1 int
	port2 int
	port3 int
}

var user User

var filePathOut string

var ports = make([]int, 0)

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user.Name = session.Value
	t, err := template.ParseFiles("templates/upload.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "upload", nil)
}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	link := r.FormValue("link")
	fmt.Println(link)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	for i := 0; i < 3; {
		port := rand.Intn(30000-20000+1) + 20000
		if !contains(ports, port) {
			ports = append(ports, port)
			if user.port1 == 0 {
				user.port1 = port
			} else if user.port2 == 0 {
				user.port2 = port
			} else if user.port3 == 0 {
				user.port3 = port
			}
			i++
		}
	}
	fmt.Println(user.port1, user.port2, user.port3)

	fileName := fileHeader.Filename
	filePath := fmt.Sprintf("Test/%s/%s", user.Name, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RunTester.RunTester(link, user.Name, fmt.Sprintf("Test/%s", user.Name), strconv.Itoa(user.port1), strconv.Itoa(user.port2), strconv.Itoa(user.port3))
	fmt.Println("Selenium test completed successfully")

	// Redirect the user to the download page
	http.Redirect(w, r, "/download", http.StatusSeeOther)
}
