/*
 * Pallier cryptosystem used for fast additive homomorphic encryption
 */

package main

import (
    "fmt"
    "crypto/rand"
    "math/big"
    )

var bitsize int = 256 //tested upto 2048
var bigzero = big.NewInt(0)
var bigone = big.NewInt(1)
var bigtwo = big.NewInt(2)

func main() {
    
    //keygen
    n,g,L,M:=KeyGen()
    msg,_:=GenerateRandom(n)
    //encrypt
    cipher:=Encrypt(msg,n,g)
    //decrypt
    messg:=Decrypt(cipher,n,L,M)
    fmt.Println("*****checking encrytion and decryption*****")
    fmt.Println("original plaintext:",msg,"\nciphertext:",cipher,"\ndecrypted plaintext:",messg)
    //test additive homomorphism
    TestAdditiveHomomorphism(n,g,L,M)
    //test Plaintext multiplication
    TestPlaintextMultiplication(n,g,L,M)
    
}

//Key generation function
func KeyGen() (*big.Int,*big.Int,*big.Int,*big.Int) {
    var n,g,L,M *big.Int
    p,_:=SafePrime()
    q,_:=SafePrime()
    a:=new(big.Int).Mul(p,q)
    pmo:=new(big.Int).Sub(p,bigone)
    qmo:=new(big.Int).Sub(q,bigone)
    b:=new(big.Int).Mul(pmo,qmo)
    gcd:=new(big.Int).GCD(bigzero,bigzero,a,b)
    if (gcd.Cmp(bigone)!=0) {
        //keep repeating till approprate values are generated
        fmt.Println("error one")
        n,g,L,M=KeyGen()
    } else {
        gcdpq:=new(big.Int).GCD(bigzero,bigzero,pmo,qmo)
        L=new(big.Int).Div(b,gcdpq)    //get lambda
        n=a    //get n
        g,M=FindGenerator(n,L) //get g and mu
    }
    return n,g,L,M
}

//Encryption function
func Encrypt(msg *big.Int,n *big.Int,g *big.Int) (*big.Int) {

    r,_:=GenerateRandom(n)
    nsq:=new(big.Int).Mul(n,n)
    gm:=new(big.Int).Exp(g,msg,nsq)
    rn:=new(big.Int).Exp(r,n,nsq)
    c:=new(big.Int).Mul(gm,rn)
    return new(big.Int).Mod(c,nsq)
}

//Decryption function
func Decrypt (cipher *big.Int,n *big.Int, L *big.Int, M *big.Int) *big.Int {

    nsq:=new(big.Int).Mul(n,n)
    x:=new(big.Int).Exp(cipher,L,nsq)
    xmo:=new(big.Int).Sub(x,bigone)
    xmomn:=new(big.Int).Div(xmo,n)
    msg:=new(big.Int).Mul(xmomn,M)
    return new(big.Int).Mod(msg,n)
}

//function for basic testing of additive homomorphism of Pallier cryptosystem
func TestAdditiveHomomorphism (n *big.Int, g *big.Int, L *big.Int, M *big.Int) {

    fmt.Println("*****running additive homomorphic testing*****")
    m1:=big.NewInt(100) //can use GenerateRandom(n) also
    m2:=big.NewInt(200)
    m_sum:=new(big.Int).Add(m1,m2)
    c1:=Encrypt(m1,n,g)
    c2:=Encrypt(m2,n,g)
    c_sum:=Decrypt(new(big.Int).Mul(c1,c2),n,L,M)
    fmt.Println("sum from plaintext: ",m_sum,"\nsum from ciphertext: ",c_sum)
}

func TestPlaintextMultiplication(n *big.Int, g *big.Int, L *big.Int, M *big.Int) {
    
    m:=big.NewInt(100)
    s:=big.NewInt(5)
    
    m_prod:=new(big.Int).Mul(m,s)
    c:=Encrypt(m,n,g)
    nsq:=new(big.Int).Mul(n,n)
    cp:=new(big.Int).Exp(c,s,nsq)
    c_prod:=Decrypt(cp,n,L,M)
    fmt.Println("scalar product from plaintext: ",m_prod,"\nscalar product from ciphertext: ",c_prod)
}

func SafePrime() (*big.Int,error) {
    p:=new(big.Int)
    for {
        q,err:=rand.Prime(rand.Reader,bitsize)
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


func FindGenerator(n *big.Int, L *big.Int) (*big.Int,*big.Int) {
    //can also take g = n+1 and M=new(big.Int).ModInverse(L,n) for simplicity
    g,_:=GenerateRandom(n)
    nsq:=new(big.Int).Mul(n,n)
    val:=new(big.Int).Exp(g,L,nsq)
    val2:=new(big.Int).Div(val,n)
    M:=new(big.Int).ModInverse(val2,n)
    if (M==nil) {
        fmt.Println("error in g")
        g,M=FindGenerator(n,L)
    }
    return g,M
}

//generate random values
func GenerateRandom(n *big.Int) (*big.Int,error) {
    s,err:=rand.Int(rand.Reader,n)
        if err  != nil {
            return nil,err
        } else {
            return s,nil
        }
}
