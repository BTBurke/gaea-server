package main

import "fmt"
import "os"
import "log"
import "github.com/elithrar/simple-scrypt"

func main() {

    pwd := os.Args[1]

    if len(pwd) == 0 {
        log.Fatal("Must specify a password to hash")
    }
    
    hash, _ := scrypt.GenerateFromPassword([]byte(pwd), scrypt.DefaultParams)
    
    fmt.Printf("%s\n", hash)
    
}