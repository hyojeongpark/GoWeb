package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type fooHandler struct{}

//인스턴스 형태로 등록.
//1. responsewriter = response에 무언가를 쓸 수 있게, request=요청검토
func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)//새로운 객체를 user에 넣어준다
	//요청의 바디에 있는 데이터를 user 형태로 변경시켜줌.
	err := json.NewDecoder(r.Body).Decode(user)//json형태로 decode해주고
	if err != nil { //에러가 발생하면
		w.WriteHeader(http.StatusBadRequest) 
		fmt.Fprint(w, "Bad Request: ", err)
		//더 이상 에러가 나지 않도록 리턴시켜준다고 함.
		return
	}

	user.CreatedAt = time.Now()
	//JSON데이터로 바꾸려고 하는데, 현재는 user 라는 구조체이기 때문에 json으로 인코딩해주는 Marshal 사용.
	//리턴의 첫번째는 byte 타입의 data, 두번째는 error.
	data, _ := json.Marshal(user)
	w.Header().Add("content-type", "application/json")//header에 다음내용을 추가해서 json이라는 것을 알려줌
	w.WriteHeader(http.StatusCreated)
	//data변수가 byte array이기 때문에 스트링으로 변환.
	fmt.Fprint(w, string(data))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World") // Fprint은 w값에 "Hello World" 값을 주어서 프린팅해라는 의미
}
func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func NewHttpHandler() http.Handler {
	//경로에 따라 분배해주는 라우터. mux라는 인스턴스를 만둘어주면, 경로를 거기에 등록하고 해당 경로를 그 인스턴스에 넘겨주는 방식.
	mux := http.NewServeMux()
	//HandleFunc은 Handler 등록하는 함수.   핸들러를 함수로 직접 등록.
	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/bar", barHandler)

	//인터페이스를 구조체의 포인터를 호출함으로써 호출하는건가?
	mux.Handle("/foo", &fooHandler{})
	return mux
}
