package main

import (
	"fmt"
	"os"
	"strconv"
	"crypto/rand"
	"math/big"
)

var one = big.NewInt(1)

const minBits = 256
const maxBits = 16384

func checkErr(bi *big.Int, err error) *big.Int {
	if err != nil {
		fmt.Println("Something bad happened during prime number generation:", err.Error())
		return nil
	}

	return bi
}

func main() {
	bitSize := 512
	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil && i <= maxBits && i >= minBits {
			bitSize = i
		}
	}
	primeChan := make(chan *big.Int)
	fmt.Println("Finding some random " + strconv.Itoa(bitSize) + "-bit long prime numbers for use as q and d")
	for i := 0; i < 2; i++ {
		go func() {
			primeChan <- checkErr(rand.Prime(rand.Reader, bitSize))
		}()
	}
	q, d := <-primeChan, <-primeChan
	fmt.Printf("q:%v\nd:%v\n\n", q, d)
	fmt.Println("Doing math on q and d to arrive at RSA exponents/modulus...")

	var pub *big.Int
	go func() {
		var pubExp *big.Int
		if len(os.Args) > 2 && os.Args[2] != "random" {
			if i, err := strconv.Atoi(os.Args[2]); err == nil {
				pubExp = big.NewInt(int64(i))
			}
		}
		if pubExp == nil || !pubExp.ProbablyPrime(bitSize) {
			pubExp = checkErr(rand.Prime(rand.Reader, bitSize))
		}
		primeChan <- pubExp
	}()
	pub = <-primeChan
	var priv, mod *big.Int
	go func() {
		primeChan <- new(big.Int).ModInverse(pub, new(big.Int).Mul(q.Sub(q, one), d.Sub(d, one)))
	}()
	
	mod = new(big.Int).Mul(q, d)
	priv = <-primeChan
	fmt.Printf("modulus:%d\n\npublic:%d\n\nprivate:%d\n\n\n", mod, pub, priv)

	fmt.Printf("Testing new keys...")
	control := big.NewInt(500)
	encrypted := new(big.Int).Exp(control, pub, mod)
	decrypted := new(big.Int).Exp(encrypted, priv, mod)

	if control.String() == decrypted.String() {
		fmt.Printf("Test passed!\nThe generated primes work for RSA encryption.\nPlace the modulus and and the public exponent in the client, and place the modulus and the private exponent in the server, and everything should work.\n")
	} else {
		fmt.Printf("Test failed!\nSomething went wrong, before and after encryption did not match the control variable.\n")
		fmt.Println(control, decrypted, encrypted)
	}
}
