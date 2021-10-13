package main

import (
	"WEB9/lzw"
	"fmt"
)

type Component interface {
	Operator(string)
}

//전역변수
var sentData string
var recvData string

//실제 기본기능
type SendComponent struct{}

//operator호출
func (self *SendComponent) Operator(data string) {
	sentData = data //data는 sentData로
}

//부가기능1(데코레이터)-압축
type ZipComponent struct {
	com Component
}

func (self *ZipComponent) Operator(data string) {
	zipData, err := lzw.Write([]byte(data)) //string->byte array로 압축
	if err != nil {
		panic(err)
	}
	//정상일 경우 decorate하고 있는 실제 component를 호출
	self.com.Operator(string(zipData)) //압축한 데이터를 호출
}

//부가기능2(데코레이터)-암호화
type EncryptComponent struct {
	key string
	com Component
}

func (self *EncryptComponent) Operator(data string) {
	//cipher패키지의 Encrypt기능을 써서 데이터와 key값을 넣어준다
	encryptData, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil { //암호화된 데이터가 오류가 있을 경우 panic
		panic(err)
	}
	self.com.Operator(string(encryptData))
}

//복호화단계
type DecryptComponent struct {
	key string
	com Component //데코레이터
}

func (self *DecryptComponent) Operator(data string) {
	decryptData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decryptData))
}

type UnzipComponent struct {
	com Component
}

func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data)) //압축이라 key값은 필요없음
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

//데이터 읽어주기
type ReadComponent struct{}

func (self *ReadComponent) Operator(data string) {
	recvData = data
}

func main() {

	sender := &EncryptComponent{ //암호화하는 컴포넌트는
		key: "abcde",
		com: &ZipComponent{ //압축하고
			com: &SendComponent{}, //보낸다
		},
	}
	sender.Operator("Hello World") //보내고 싶은 데이터

	fmt.Println(sentData)

	receiver := &UnzipComponent{ //압축부터 풀기
		com: &DecryptComponent{
			key: "abcde",
			com: &ReadComponent{},
		},
	}

	receiver.Operator(sentData)
	fmt.Println(recvData)
}
