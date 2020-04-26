package main

import (
	"fmt"
	"os"
	"strconv"
	"crypto/rand"
	"math/big"
)

var pub = big.NewInt(65537)
var one = big.NewInt(1)
var mod = new(big.Int)
var priv = new(big.Int)
//var pub = new(big.Int)

func init() {
	bitSize := 512
	if i, err := strconv.Atoi(os.Args[1]); err == nil && i <= 4096 && i >= 512 {
		bitSize = i
	}
	tmp, err := rand.Prime(rand.Reader, bitSize)
	if err != nil {
		fmt.Println("something bad happened")
		os.Exit(100)
	}
	tmp2, err := rand.Prime(rand.Reader, bitSize)
	if err != nil {
		fmt.Println("something bad happened")
		os.Exit(100)
	}
	p, err := rand.Prime(rand.Reader, bitSize)
	if err != nil {
		fmt.Println("something bad happened")
		os.Exit(100)
	}
	pub.Set(p)
	fmt.Println("pub:", pub)
	mod.Mul(tmp, tmp2)
	fmt.Println("mod:", mod)
	phi := new(big.Int)
	phi.Mul(tmp.Sub(tmp, one), tmp2.Sub(tmp2, one)) //.Mul(tmp2.Sub(one)) 
	fmt.Println("priv:", priv.ModInverse(pub,phi))
}

func main() {
	fmt.Printf("Testing new RSA keys...")
	b := big.NewInt(500)
	before := b.Exp(b, pub, mod)
	after := b.Exp(b, priv, mod)//, //b.ModExp(priv, mod), b)
	if before.String() == "500" && after.String() == "500" {
		fmt.Printf("Test passed!  You're good to go, place mod and pub in the client.")
	} else {
		fmt.Println("Test failed!  Something went wrong, before and after encryption and decryption do not match the control var")
	}
//	fmt.Println(big.NewInt(500).ModExp(pub, mod))
//	random, _ := rand.Prime(rand.Reader, 4096)
//	fmt.Println(b.Exp(65537), b.Exp(65537))
//	fmt.Println(b.Exp(65537))
//	for 
}
