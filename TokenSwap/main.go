package main

import (
	"fmt"
)

var money float64 = 100
var coinNumber float64 = 100
var coinPrice float64 = 1

//전제 : 코인은 한번의 입력에 하나씩만 빠짐.
//1. 코인 가격을 입력 받는다 -> scanf
//2. 함수 내에서 현재 코인 가격과 입력 받은 가격 비교후 코드를 분기한다.
//2-1 코인가격이 1이다 (코드 실행 후 처음 가격 입력받음)
//갯수 전역변수에서 1 빼주고 가격전역변수 변경.
//2-2 코인가격 = 입력받은 가격
//2-3 2-2가 아닐떄

func rebalace(enteredP float64) {

	if coinPrice == 1 {
		fmt.Println("코인 첫 가격")
		fmt.Println("코인 가격 : ", coinPrice, "입력받은 가격 : ", enteredP)
		coinNumber -= 1
		money += coinPrice

		coinPrice = money / coinNumber
		//coinPrice := fmt.Sprintf("%.2f", money/coinNumber)
		fmt.Println("남은 돈 : ", money, "남은 코인 : ", coinNumber, "코인가격 : ", coinPrice)
		fmt.Println("-----------------------------------------")
		fmt.Println("-----------------------------------------")

	} else if coinPrice == enteredP {
		fmt.Println("코인 가격 : ", coinPrice, "입력받은 가격 : ", enteredP)
		coinNumber -= 1
		money += coinPrice

		coinPrice = money / coinNumber
		//coinPrice := fmt.Sprintf("%.2f", money/coinNumber)
		fmt.Println("남은 돈 : ", money, "남은 코인 : ", coinNumber, "코인가격 : ", coinPrice)
		fmt.Println("-----------------------------------------")
		fmt.Println("-----------------------------------------")
	} else {
		fmt.Println("코인 가격 : ", coinPrice, "입력받은 가격 : ", enteredP)
		fmt.Println("입력한 가격이 코인 가격과 맞지 않습니다")
		fmt.Println("-----------------------------------------")
		fmt.Println("-----------------------------------------")
	}

}

func main() {
	for {
		var enteredPrice float64
		fmt.Println("현재 돈 : ", money, "현재 코인수 : ", coinNumber, "코인 가격 : ", coinPrice)
		fmt.Printf("코인 가격을 입력해 주세요 : %f", coinPrice)
		fmt.Scanln(&enteredPrice)
		fmt.Println("입력받은 가격 : ", enteredPrice)

		rebalace(enteredPrice)
	}
}
