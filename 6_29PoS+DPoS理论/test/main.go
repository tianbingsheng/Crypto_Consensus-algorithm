package main

import (
	"fmt"
	"math/big"
)

func main() {

	var l = 256
	var b = big.NewInt(int64(l))
	b.Bytes()
	fmt.Print(b.Bytes())


	var c = big.NewInt(0)
	c.SetBytes(b.Bytes())
	fmt.Println(c.Int64())

}
