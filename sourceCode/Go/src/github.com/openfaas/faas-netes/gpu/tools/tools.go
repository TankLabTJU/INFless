//File  : toos.go
//Author: Yanan Yang
//Date  : 2020/4/7
package tools

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomText(len int) (randText string) {
	//set random seed
	rand.Seed(time.Now().UnixNano())
	var captcha string
	for i := 0; i < len; i++ {
		//generate number form 0 to 9
		num := rand.Intn(10)
		//change number to string
		captcha += strconv.Itoa(num)
	}
	return captcha
}