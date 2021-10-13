package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//JSON담을 struct
type User struct {
	FirstName string `json:"first_name"` //go는 _가 들어간 이름을 싫어함
	LastName  string `json:"last_name"`  //annotation으로 json에서는
	Email     string `json:"email"`      //어떤형태로 표시할지 설정
	CreatedAt time.Time
}

type fooHandler struct{}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)                           //USer struct에 인스턴스만들기
	err := json.NewDecoder(r.Body).Decode(user) //request body는 newdecoder의 reader
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user.CreatedAt = time.Now()   //go structure형태->
	data, _ := json.Marshal(user) //json형태로/ _는 error
	//response의 포맷이 text:plain타입이라
	//json모양처럼 예쁘게 출력하도록w.Header에서 설정
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
	//fmt.Fprint(w, "foo!!!!!!!!!")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" { //name이 없으면
		name = "World" //name은 world다
	}
	fmt.Fprintf(w, "Hello %s!", name) // localhost:3000/bar?name="hyojeong"
}

//main
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World") //writer에 Print(response)하라
	}) // "/" = 절대경로(handler)

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", mux) //listenandserver= 웹서버 구동실행
}
