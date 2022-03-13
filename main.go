package main

import(
	"example.com/m/2_1/tempconv"
	"fmt"
)

func main(){
	var k tempconv.Kelvin = 0
	fmt.Print(tempconv.KtoC(k))
}