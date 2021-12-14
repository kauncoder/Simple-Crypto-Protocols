package main

import (
    "fmt"
    "crypto/rand"
    "math/big"
    "time"
    )

var bits int = 256 //tested upto 2048

func main() {
    
    g:=FindGenerator()
    p:=new(big.Int)
    p,_=SafePrime()
    //channel to communicate b/w alice and bob
    c:= make (chan *big.Int)
    go IsAlice(p,g,c)
    go IsBob(p,g,c)
    time.Sleep(time.Millisecond*100)
}

func IsAlice (p *big.Int,g *big.Int, c chan *big.Int) {
    
    //generate private key 
    a,_:=GenerateSecret(p)
    //calcualte public key
    ga:=new(big.Int).Exp(g, a, p)
    //recevie bobs public key
    bpub:= <- c
    //send alices public key
    c <- ga
    //calculate gab
    gba:=new(big.Int).Exp(bpub, a, p)
    fmt.Println("alice shared secret",gba)
}

func IsBob (p *big.Int,g *big.Int,c chan *big.Int) {
    //generate private key 
    b,_:=GenerateSecret(p)
    //calcualte public key
    gb:=new(big.Int).Exp(g, b, p)
    c <- gb  //sends public key to alice
    apub:= <- c //waits for alices public key
    gab:=new(big.Int).Exp(apub, b, p)
    fmt.Println("bob shared secret",gab)    
}

func SafePrime() (*big.Int,error) {
//find prime p = 2q+1 st q is also a suitable prime   
//keep generating random primes till p=2q+1 is satisfied
    p:=new(big.Int)
    for {
        q,err:=rand.Prime(rand.Reader,bits)
        if err  != nil {
            return nil,err
        }
        one:=big.NewInt(1)
        p = p.Lsh(q, 1)
		p = p.Add(p, one)
        if p.ProbablyPrime(20){
            return p,nil
        }
    }
     return nil,nil
}

func FindGenerator() *big.Int {
    //can take g=2 for now and modify later
    two:=big.NewInt(2)
    return two
}

//generate private key
func GenerateSecret(p *big.Int) (*big.Int,error) {
    
    s,err:=rand.Int(rand.Reader,p)
        if err  != nil {
            return nil,err
        } else {
            return s,nil
        }
}

func KeyDeriviationFucn () {
    //not needed for this but can be added later if needed
}
