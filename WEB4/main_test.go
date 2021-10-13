package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	path := `C:\Users\hyo\Pictures\블로그사진\jamescorden.png`
	//파일열기
	file, _ := os.Open(path)
	//함수종료전 닫아주기
	defer file.Close()
	//1.os에 파일핸들자원요청 2. os가 프로그램에 파일핸들제공 3. 핸들반환
	//3번 반환이 되기전 끝나면 자원을 먹어버리기 때문에 defer로 지연시켜줌

	//png파일명이 같은경우에 다시 올려보기 위함이다
	//os.RemoveAll("./uploads")

	//데이터를 버퍼에 실어보자
	buf := &bytes.Buffer{} //데이터가 저장되는곳
	writer := multipart.NewWriter(buf)
	//Base(path)=>file이름만 잘라내기
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	//multi에 데이터를 넣어준다
	assert.NoError(err)
	//복사
	io.Copy(multi, file)
	//닫아주기
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-type", writer.FormDataContentType())
	uploadsHandler(res, req)
	//데이터가 보내짐
	assert.Equal(http.StatusOK, res.Code)
	//여기까지하면 uploads에 파일이 생긴다

	//올린파일과 localhost:8080에 업로드되는 것이 같은지 확인
	uploadFilePath := "./uploads/" + filepath.Base(path)
	//os.Stat 는 file의 정보를 보여줌
	_, err = os.Stat(uploadFilePath)
	//에러가 없어야 한다 (이걸 통과하면 파일이 있음을 뜻한다)
	assert.NoError(err)
	//비교해보기
	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(uploadData, originData)
}
