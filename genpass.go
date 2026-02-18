package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("TranGiaTien@8386"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}