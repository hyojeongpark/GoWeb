package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

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

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	//gorilla mux의 분배기능으로 여러개 설정해도 가능하다!
	resp, err := http.Get(ts.URL + "/users/1")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:1")

	resp, err = http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "User Id:56")
}

//---4단계 TestCreateUser Post 테스팅 완료검사---//
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

//---5단계 TestDeleteUser 테스팅
func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil) //Delete는 기본메소드가 아니다. 다만, Do(req) request에 메소드정의, Id는 임의로/Users/1, Body 값은 없으므로 nil로 한다.
	resp, err := http.DefaultClient.Do(req)                     //Do(req): Do(Method)="DELETE" 메소드를 정의
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)          // 아래 log로 찍어보기위해 ReadAll(Body)로 읽어온다, delete할것이 없었다
	log.Print(string(data))                       //log로 찍어본다.
	assert.Contains(string(data), "No User Id:1") //"지울게 없었다 즉 유저가 없었다"는 메세지를 포함해야 한다, 위 임의등록시킨 1번ID를 등록시킨 적이 없다

	//[메모]터미널-> cd WEB6/myapp으로 이동-> '$go test' 실행 -> goconvey로 시간확인
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"HyoJeong", "last_name":"Park", "email":"123@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//---위 POST방식 Json 정보를 서버가 user정보를 받아서 user정보를 리턴하는 부분
	user := new(User)                             //
	err = json.NewDecoder(resp.Body).Decode(user) //서버가 보낸 데이터를 읽어온다,Encoder/Decoder는 스트림기반 데이터를 다루고,Encoder는 go value를 JSON문자열로 반환, Decorder는 반대
	assert.NoError(err)                           //아무런 문제err가 없다.
	assert.NotEqual(0, user.ID)                   //user ID가 0이 아니다.   등록 기록됨

	//[삭제]
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil) //Delete는 기본메소드가 아니다. 다만, Do(req) request에 메소드정의, Id는 임의로/Users/1, Body 값은 없으므로 nil로 한다.
	resp, err = http.DefaultClient.Do(req)                     //Do(req): Do(Method)="DELETE" 메소드를 정의
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)                // 아래 log로 찍어보기위해 ReadAll(Body)로 읽어온다, delete할것이 없었다
	log.Print(string(data))                            //log로 찍어본다.
	assert.Contains(string(data), "Deleted User Id:1") //"지울게 없었다 즉 유저가 없었다"는 메세지를 포함해야 한다, 위 임의등록시킨 1번ID를 등록시킨 적이 없다

}

//6단계 updtae test
func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"id":1,"first_name":"updated","last_name":"updated","email":"updated@gmail.com"}`))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:1")

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"HyoJeong","last_name":"Park","email":"123@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d,"first_name":"updated"}`, user.ID)
	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.NotEqual(http.StatusBadRequest, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)
	assert.Equal("updated", updateUser.FirstName)
	assert.Equal(user.LastName, updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)
}
