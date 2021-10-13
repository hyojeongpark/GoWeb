package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//파일명_test 는 테스트모드가 됨 (컨벤션이라 부름)

//testing 패키지의 T 포인터를 인자로 받음.
func TestIndexPathHandler(t *testing.T) {
	//1. 전처리작업 세팅 assertion
	assert := assert.New(t)
	//2. res,req
	res := httptest.NewRecorder() //test용 리코더
	//메소드, 타겟, 바디
	req := httptest.NewRequest("GET", "/", nil)

	//barHandler(res,req) =>이렇게하면 /의 타겟인 indexhandler에 맞게 적용됨
	//3. mux로 serveHttp로 연결하여 분배
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	//응답코드 확인
	assert.Equal(http.StatusOK, res.Code) //유사성검사 assert
	//res.body를 바로 가져올 수 없어서 ioutil 사용해 버퍼값 전부 가져오게 함. 리턴의 두번쨰 인자는 err
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

//-------------2단계 테스팅--------------------//
func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World!", string(data))
}

//----------3단계 테스팅------------------------//
func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=Gryffindor", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello Gryffindor!", string(data))
}

//------4단계 테스팅---------------------------//
func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t) //검증,test tool

	res := httptest.NewRecorder()                  //응답
	req := httptest.NewRequest("GET", "/foo", nil) //nil (b0dy)값이 업기 때문에 build fail

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code) //res.code =>writeheader부분
}

//------5단계 테스팅---------------------------//
func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t) //유사성검사 도구

	res := httptest.NewRecorder() //네트워크를 사용하지 않고 local에서 test하기 위해서
	//실제로는 httptest가 아닌 httpreponse이다(test환경이 아닐 때)
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{"first_name" : "HyoJeong", "last_name" : "Park", "email" : "123@gmail.com"}`))
	//io.Reader 인터페이스 따르는 읽기 인스턴스 생성/ `annotation`

	mux := NewHttpHandler() //mux=분배 /서버로 요청이 왔을때 이 서버가 경로를 핸들링 해주는 기능
	mux.ServeHTTP(res, req) //실질적인 연결(인터페이스)
	//ServeHttp(요청과 응답을 전달해주는 매개체역할)그리고 그 안에 http양식에 맞추어 내용이 들어있다
	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)                             //app.go에서 가져오기
	err := json.NewDecoder(res.Body).Decode(user) //json 내장함수newdecoder
	assert.Nil(err)                               //에러같이 nil인지 아닌지
	assert.Equal("HyoJeong", user.FirstName)
	assert.Equal("Park", user.LastName)
}

//http=패키지
// 서버=요청할 주소에 있는 컴퓨터
// 로컬(client)=내 컴퓨터,핸드폰

// 로컬이 서버에 request요청(네이버의 index.html파일을 줘!)하면 response응답(naver.com)을 준다
// 이렇게 오갈때(ex>한미협약때 양식이 필요한것처럼)
//  보내는 양식(규약:프로토콜)에 맞춰 보내야하는데 이걸 HTTP라고함!
// hyper text를 주고받을 때 쓰는 포로토콜을 HTTP을 통해쓴다
