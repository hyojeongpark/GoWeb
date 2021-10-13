package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Cookie struct {
	Name  string
	Value string

	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

	MaxAge int
}

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"), //환경변수에 저장된ID가져오기 os.Getenv
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state) //googleOauthConfig는 유저를 어떤경로로보내야하는지 URL코드보냄
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) //쿠키만료시간 하루
	b := make([]byte, 16)                            //랜덤한 16byte짜리 Array만들기
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b) //string으로 바꿔주기
	//cookie struct만들어주기
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	//위에 저장한 쿠키 읽어오기
	oauthstate, _ := r.Cookie("oauthstate")

	//oathstate와 formvalue가 다르면 redirect해주기
	if r.FormValue("state" != oauthstate.Value) {
		log.Printf("invalid google oauth state cooke:%s state : %s\n", oauthstate.Value)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data)) //에러가 없을시 유저의 정보를 보낸다
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}
	return ioutil.ReadAll(resp.Body)
}

func main() {
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)
	n := negroni.Classic()
	n.Usehandler(mux)
	http.ListenAndServe(":3000", n)
}

//터미널에서 설치
//go get golang.org/x/oauth2
//go get cloud.google.com/go
