package RGBtrie

import (
	"fmt"
	"io/ioutil"
)

func main() {
	b, e := ioutil.ReadFile("c:/Users/85261/Desktop/1.txt")
	if e != nil {
		fmt.Println("read file error")
		return
	}
	fmt.Println(string(b))
}

