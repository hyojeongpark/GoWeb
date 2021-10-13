package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//User정보 등록 실행후 No User정보 error발생시 해결책은user[map]을 만들고 cereate한 user정보를 등록
var userMap map[int]*User //NewHandler 시점 초기화	선언과 초기화
var lastID int            //마지막 ID등록

// user struct
type User struct { //json을 읽을수 있는 struct를 만들어 준다. ID 정수형 추가
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

//실행1> 테스트1이 잘 작동하면 진행
func indexHandler(w http.ResponseWriter, r *http.Request) { //
	fmt.Fprint(w, "Hello World") // Fprint은 w값에 "Hello World" 값을 주어서 프린팅해라는 의미
}

//실행2
func usersHandler(w http.ResponseWriter, r *http.Request) { //
	fmt.Fprint(w, "Get UserInfo by /Users/{id}") // Fprint은 w값에 "Hello World" 값을 주어서 프린팅해라는 의미
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//mux중에 Vars가 있는데 (id의) request를 알아서 파싱해주는 기능
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) //vars는 string이므로 Atoi로 int정수형으로 바꾸면 첫번째 인티져형id와, 두번째 err가 나온다
	//vars := mux.Vars(r)                   //r은 request vars가 Id 파싱을 자동으로 해준다
	//fmt.Fprint(w, "User Id:", vars["id"]) // vars[id] 형식으로 작성해야 파싱값이 User Id:로 들어가 w로 출력된다
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //변환과정에 에러발생시 StatusBadRequest출력
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id] //해당 id인티져 유저맵에 있으면 행당유저의 정보를 보여주고
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id) // 없으면 No User ID 메세지를 보여준다
		return
	}
	//user := new(User)
	//user.ID = 2
	//user.FirstName = "Hyunjun"
	//user.LastName = "Park"
	//user.Email = "123@gmail.com"

	// err가 아니라면 실행
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(user) //[]슬라이스 바이트나 문자열 사용기반, Go 밸류를 JSON 문자열로 변환	<-> upMarchal() json을 go로 변환?
	fmt.Fprint(w, string(data))
}

//	POST 전송방식 이기에 선언과 초기화 변환 작업의 번거로움
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) //	Decode는 JSON문자열을 go Value 반환, Encoder는 반대 Go Value를 JSON문자열로(json.Marshal=)
	//	같은 json 패키지의 Decode() Marshal()을 사용
	if err != nil { //에러가 있다면 잘못된 JSON형식
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// Created User
	lastID++                    //userID가 만들어질때마다 하나씩 증가
	user.ID = lastID            //매번 바뀐다. 하나생겨날때마다 매번 증가
	user.CreatedAt = time.Now() //현재시간 설정
	userMap[user.ID] = user

	w.WriteHeader(http.StatusCreated) //잘 생성되었다고 알려준다	// 201
	data, _ := json.Marshal(user)     //struct 나 int 타입을 []byte Slice로 변환
	fmt.Fprint(w, string(data))
}

//실행후 No User정보 error발생시 해결책은user[map]을 만들고 cereate한 user정보를 등록한후 user정보가 있으면 그 정보를 리턴하고 user정보가 없으면 No User정보

//main과 연결되는 (handler를 연결해주는) NewHandler
func NewHandler() http.Handler { //	interface	외부통신(공개)기능
	userMap = make(map[int]*User) //맵 초기화후 언제등록? Create User때 해야한다
	lastID = 0                    //사용하기 전에 초기화 해준다 마지막 값을 0
	mux := mux.NewRouter()        //gorilla/mux 패키지 자동 임포트

	mux.HandleFunc("/", indexHandler)                           //하위경로 미정의시는 상위경로가 자동호출된다
	mux.HandleFunc("/users", usersHandler).Methods("GET")       // gorilla Mux가 지원기능 Method 경로 설정해주면 GET 핸들어 호출
	mux.HandleFunc("/users", createUserHandler).Methods("POST") // 메소드 POST방식은 CreateUserHandler가 호출
	//	gorilla/mux의	*mux.Route	-	.Methods("GET")("POST") 전송방식,기능
	//이해하려고 하지말자 깃허브에 나와있음 [0-9] 이 부분
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)
	return mux
}
