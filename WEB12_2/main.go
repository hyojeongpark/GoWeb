package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
)

var rd *render.Render

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "hyojeong", Email: "123@gmail.com"}
	rd.JSON(w, http.StatusOK, user)
	// w.Header().Add("Content_type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, err)
		rd.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	user.CreatedAt = time.Now()
	rd.JSON(w, http.StatusOK, user)
	// w.Header().Add("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.New("Hello").ParseFiles("templates/hello.tmpl")
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprint(w, err)
	// 	return
	// }
	rd.HTML(w, http.StatusOK, "hello.tmpl", "HJP")
	// tmpl.ExecuteTemplate(w, "hello.tmpl", "HJP")
}
func main() { //github.com/gorilla/pat 설치
	rd = render.New(render.Options{
		Directory:  "template",
		Extensions: []string{".html", ".tmpl"}, //확장자명 불러오기
	})
	mux := pat.New() //gorilla/mux 패키지 자동 임포트

	mux.Get("/users", getUserInfoHandler) //하위경로 미정의시는 상위경로가 자동호출된다
	mux.Post("/users", addUserHandler).Methods("POST")
	mux.Get("/hello", helloHandler)

	http.ListenAndServe(":3000", mux)
}

//go get github.com/unrolled/render 설치
