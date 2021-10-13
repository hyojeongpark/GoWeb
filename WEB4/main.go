package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

//-------2. os 자원을 쓰겠다! -uploadsHandler-----------
func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	//file읽는부분
	uploadfile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	defer uploadfile.Close()

	//-----3. 저장할 file 공간만들기(위치:uploads)---
	dirname := "./uploads"
	os.MkdirAll(dirname, 0777) //8진수형태 리눅스 명령어 read,write,execute
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	//파일만들기 file=비어있는파일
	file, err := os.Create(filepath)
	//os자원반납 함수종료전 close해주기
	defer file.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //문제가 생겼음을 표시
		fmt.Fprint(w, err)
		return //에러반환처리
	}
	//-----------4.복사copy--------------
	io.Copy(file, uploadfile) //Copy(어디에copy,무엇을copy)
	w.WriteHeader(http.StatusOK)
	//filepath는 upload경로를 알려준다
	fmt.Fprint(w, filepath)
}

//-----------1차 Step------------
func main() {
	http.HandleFunc("/uploads", uploadsHandler) //파일업로드 핸들러
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)
}
