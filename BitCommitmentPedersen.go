/*
 * Pederdsen committment scheme uses DLOG assumption to provide a multiple-bits commitment scheme
 * Main advantage over old commitment scheme:
 * 1. Homomorphic
 * 2. Uses simpler assumptions (DLOG)
 */

package main

import (
    "fmt"
    "crypto/rand"
    "math/big"
    "time"
    )

var bits int = 256 //tested upto 2048
var bigzero = big.NewInt(0)
var bigone = big.NewInt(1)
var bigtwo = big.NewInt(2)

func main() {
    
    g:=bigtwo
    p,_:=SafePrime()
    s,_:=GenerateSecret(p)
    h:=new(big.Int).Exp(g,s,p) //choosing h as g^s
    //channel to communicate b/w sender and receiver
    c:= make (chan *big.Int)
    go IsSender(p,g,h,c)
    go IsReceiver(p,g,h,c)
    time.Sleep(time.Millisecond*100) //main protocol waits for others to get over
}

//Sender function that commits and reveals the bit value
func IsSender(p *big.Int,g *big.Int,h *big.Int, c chan *big.Int) {
    //value to be committed 
    m,_:=GenerateSecret(p)
    //generate random value k
    r,_:=GenerateSecret(p)
    //calcualte h^r 
    hr:=new(big.Int).Exp(h, r, p)
    //calcualte g^b
    gm:=new(big.Int).Exp(g, m, p)
    gb:=new(big.Int).Mul(hr,gm)
    commit:=new(big.Int).Mod(gb,p)
    fmt.Println("sender committed via Pedersen",m)
    Commit(commit,c)
    Reveal(m,r,c)
    
}

func Commit(commit *big.Int, c chan *big.Int) {
    //commit values to Receiver
    c <- commit
    //c <- hxorb  
}

func Reveal(m *big.Int, r *big.Int, c chan *big.Int) {
    c <- m
    c <- r
}


//Receiver function that checks and extracts committed value 
func IsReceiver (p *big.Int,g *big.Int,h *big.Int,c chan *big.Int) {

    //receives the values of f and hb from sender
    commit:= <- c 
    //waits for reveal
    msgreveal:= <- c
    randomreveal:= <- c
    //check values
    hrcalc:= new(big.Int).Exp(h, randomreveal, p) //caclucalte h^r
    gmcalc:= new(big.Int).Exp(g,msgreveal,p)    //calcualte g^m
    commitcalc:= new(big.Int).Mul(hrcalc,gmcalc)   //getting h^r.g^m
    commitcalc= new(big.Int).Mod(commitcalc,p)  //modding with p
    
    if commitcalc.Cmp(commit)!=0 {
        fmt.Println("cheating sender detected")
    } else {
        fmt.Println("receiver was revealed",msgreveal)
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
