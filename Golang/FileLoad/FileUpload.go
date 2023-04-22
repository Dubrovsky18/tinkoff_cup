package FileLoad

import (
	"fmt"
	"github.com/Dubrovsky18/tinkoff_cup/Tests"
	"html/template"
	"io"
	"net/http"
	"os"
)

type User struct {
	Name  string
	port1 string
	port2 string
	port3 string
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

	//for i := 0; i < 3; {
	//	port := rand.Intn(30000-20000+1) + 20000
	//	if !contains(ports, port) {
	//		ports = append(ports, port)
	//		if user.port1 == "" {
	//			user.port1 = string(port)
	//		} else if user.port2 == "" {
	//			user.port2 = string(port)
	//		} else if user.port3 == "" {
	//			user.port3 = string(port)
	//		}
	//		i++
	//	}
	//}
	//fmt.Println(user.port1, user.port2, user.port3)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

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
	//Tests.RunTester(userFolder, fileName, link, user.Name, user.port1, user.port2, user.port3)
	filePathOut = Tests.RunTester(filePath, fileName, link, user.Name)
	fmt.Println(filePathOut)

	fmt.Println("Selenium test completed successfully")

	// Redirect the user to the download page
	http.Redirect(w, r, "/download", http.StatusSeeOther)
}
