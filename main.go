package main

import (
	"fmt"
	"os"
	"strconv"
	"crypto/rand"
	"math/big"
)

var one = big.NewInt(1)
var tmp = new(big.Int)

const minBits = 512
const maxBits = 16384

func checkErr(bi *big.Int, err error) *big.Int {
	if err != nil {
		fmt.Println("Something bad happened during prime number generation:", err.Error())
		return nil
	}

	return bi
}

func main() {
	bitSize := minBits
	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil && i <= maxBits && i >= minBits {
			bitSize = i
		}
	}
	var q = checkErr(rand.Prime(rand.Reader, bitSize))
	var d = checkErr(rand.Prime(rand.Reader, bitSize))
	fmt.Printf("q:%v\nd:%v\n\n", q, d)
	fmt.Println("Deriving RSA exponents/modulus from q and d...")

	var pub *big.Int
	if len(os.Args) > 2 && os.Args[2] != "random" {
		if i, err := strconv.Atoi(os.Args[2]); err == nil {
			if pubExp := big.NewInt(int64(i)); pubExp.ProbablyPrime(bitSize) {
				pub = pubExp
			}
		}
	}
	if pub == nil {
		pub = checkErr(rand.Prime(rand.Reader, bitSize))
	}
	
	mod := new(big.Int).Mul(q, d)
	priv := new(big.Int).ModInverse(pub, tmp.Mul(q.Sub(q, one), d.Sub(d, one)))
	fmt.Println("publicExponent:", pub)
	fmt.Println("privateExponent:", priv)
	fmt.Println("modulus:", mod)

	fmt.Printf("Testing new RSA keys...")
	control := big.NewInt(500)
	encrypted := tmp.Exp(control, pub, mod)
	decrypted := tmp.Exp(encrypted, priv, mod)

	if control.String() == decrypted.String() {
		fmt.Printf("Test passed!\nThe generated primes work for RSA encryption.\nPlace the modulus and and the publicExponent in the client, and place the modulus and the privateExponent in the server, and everything should work.\n")
	} else {
		fmt.Printf("Test failed!\nSomething went wrong, before and after encryption did not match the control variable.\n")
	}
}
