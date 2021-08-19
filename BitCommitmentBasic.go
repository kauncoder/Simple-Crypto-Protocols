package main

import (
    "fmt"
    "crypto/rand"
    "math/big"
    "time"
    )

var bits int = 32 //tested upto 2048
var bigzero = big.NewInt(0)
var bigone = big.NewInt(1)
var bigtwo = big.NewInt(2)

func main() {
    
    g:=FindGenerator()
    p,_:=SafePrime()
    //channel to communicate b/w sender and receiver
    c:= make (chan *big.Int)
    go IsSender(p,g,c)
    go IsReceiver(p,g,c)
    time.Sleep(time.Millisecond*100) //main protocol waits for others to get over
}

//Sender function that commits and reveals the bit value
func IsSender(p *big.Int,g *big.Int, c chan *big.Int) {
    //value to be committed (can be zero or one)
    b:=bigzero
    //generate random value k
    k,_:=GenerateSecret(p)
    //calcualte f(k)
    f:=new(big.Int).Exp(g, k, p)
    //calcualte h(k)
    h:=new(big.Int).Mod(k,bigtwo)
    hb:=new(big.Int).Xor(h,b)
    fmt.Println("sender committed",b)
    Commit(f,hb,c)
    Reveal(k,c)
    
}

func Commit(f *big.Int,hb *big.Int, c chan *big.Int) {
    //commit values to Receiver
    c <- f
    c <- hb   
}

func Reveal(b *big.Int, c chan *big.Int) {
    c <- b
}


//Receiver function that checks and extracts committed value 
func IsReceiver (p *big.Int,g *big.Int,c chan *big.Int) {
    
    b:=bigtwo //initializing value with anything other than 0 o 1
    //receives the values of f and hb from sender
    f:= <- c 
    hb:= <-c
    //waits for reveal
    k:= <-c
    //check if value of OWF is right
    fc:=new(big.Int).Exp(g, k, p) 
    hc:=new(big.Int).Mod(k,bigtwo)
    
    if fc.Cmp(f)!=0 {
        fmt.Println("cheating sender detected")
    } else {    //calculate h(k) and extract b
        b=new(big.Int).Xor(hb,hc)
        fmt.Println("receiver revealed",b)
    }
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
    return bigtwo
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
