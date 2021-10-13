package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//테스트1
func TestIndex(t *testing.T) {
	//검증도구
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	//defer로 함수가 끝나기전에 닫아주기
	defer ts.Close()

	//에러
	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	//reponse로 나온 값을 data에 넣고
	data, _ := ioutil.ReadAll(resp.Body)
	//값이 같은지 확인
	assert.Equal("Hello World", string(data))
}

//테스트2
func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Get UserInfo")
}

//테스트3
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // app.go NewHandler()호출
	defer ts.Close()

	//--- Post로 보냈을때 users에 POST방식전송, 정보는 Json 방식으로 전송한다
	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"HyoJeong", "last_name":"Park", "email":"123@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//---위 POST방식 Json 정보를 서버가 user정보를 받아서 user정보를 리턴하는 부분
	user := new(User)                             //
	err = json.NewDecoder(resp.Body).Decode(user) //서버가 보낸 데이터를 읽어온다,Encoder/Decoder는 스트림기반 데이터를 다루고,Encoder는 go value를 JSON문자열로 반환, Decorder는 반대
	assert.NoError(err)                           //아무런 문제err가 없다.
	assert.NotEqual(0, user.ID)                   //user ID가 0이 아니다.   등록 기록됨
	// user 정보 등록(기록)
	id := user.ID                                               // := 인수 입력, 함수호출시 함수로 값을 전달해주는 값
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id)) // = 매개변수 입력, ID정보를 Get방식으로/users/ID를 넣어서 오도록 만든다
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	user2 := new(User)                             // := 인수 입력, 함수호출시 함수로 값을 전달해주는 값
	err = json.NewDecoder(resp.Body).Decode(user2) //정보 파싱, =매개변수 입력, 함수정의에서 전달받은 인수를 함수 내부로 전달하는 변수
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID) //user1 ID는 create한 user.ID이고 새로운정보를 받은user2 ID user2.ID
	assert.Equal(user.FirstName, user2.FirstName)
}

//테스트4(User type 생성 후)

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	//gorilla mux의 분배기능으로 여러개 설정해도 가능하다!
	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:89")

	resp, err = http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:56")

	resp, err = http.Get(ts.URL + "/users/12")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:12")
}
