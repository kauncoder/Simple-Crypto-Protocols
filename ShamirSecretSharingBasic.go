package main

import (
    "fmt"
    "crypto/rand"
    "math/big"
    )

var bits int = 32
var n int = 5  // (n is the total shares to be generated)
var t int = 3  // (t is threshold for reconstructing secret)

/*
 * Dealer-based secret sharing
 * generate private key using RSA
 * generate coefficients
 * generate shares for individuals
 * create a function to get secret from shares (Lagrange interpolation)
 */

func main() {
    //generate the secret (using random number for now but can also use RSA/ECDSA key)
    p,_:=SafePrime()
    secret,_:= GenerateSecret(p)
    //for threshold t (here t=3) players generate t-1 coefficients to create the polynomial
    coeffs := CreatePolynomial (secret, p, t)
    //generate shares for each of the n (here n=5) participant; returns (x,y) value for each
    xshares,yshares:=GenerateShares(t, n, p ,coeffs)
    
    //Lagrange calculator takes t+ shares and returns our secret
    participants:=3  //number of participants available to reconsruct secret (has to be >=t)
    //we are taking first t values for our example but it can be any shares as long as threshold is met
    reconX:=xshares[0:participants]
    reconY:=yshares[0:participants]
    reconstructedSecret:=ReconstructSecret(reconX,reconY,participants,p)
    if secret.Cmp(reconstructedSecret)==0 {
        fmt.Println("secret reconstructed successfully")
    } else {
        fmt.Println("failed at reconstruction")
    }
        
}

//generate a Safe Prime for the field we will use
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

//generate our secret (using random number here for now)
func GenerateSecret(p *big.Int) (*big.Int,error) {
    
    //privatekey, err := rsa.GenerateKey(rand.Reader, 32)
    s,err:=rand.Int(rand.Reader,p)
        if err  != nil {
            return nil,err
        } else {
            return s,nil
        }
}

//generate coefficients for the polynomial
func CreatePolynomial (secretkey *big.Int, p *big.Int, k int) ([]*big.Int){
   
    coeffs := make([]*big.Int, t) 
    //the zeroth coefficient is our secret
    coeffs[0]=secretkey 
     //generate t-1 values for the rest of the coefficients
    for i:=1;i<t;i++ {
        coeffs[i],_=GenerateSecret(p)
    }
    return (coeffs)
}

//generate shares for each of the n parties
func GenerateShares(t int, n int, p *big.Int,coeffs []*big.Int) ([]*big.Int, []*big.Int) {
    
    xshares := make([]*big.Int, n) //stores the x value in (x,y) for each participant
    yshares := make([]*big.Int, n) //stores the y value in (x,y) for each participant
    for i:=0;i<n;i++ {
        xshare:=big.NewInt(int64(i+1))  //assigning easily calculatable x shares
        yshare:=big.NewInt(int64(0))
        xshares[i]=xshare
        for j:=0;j<t;j++ {
            xsharepow:=new(big.Int).Exp(xshare,big.NewInt(int64(j)),p)
            yshare=new(big.Int).Add(yshare,new(big.Int).Mul(xsharepow,coeffs[j]))
        }
        yshares[i]=new(big.Int).Mod(yshare,p)
    }
    return xshares,yshares
}

//Lagrange interpolation to reconstruct the secret
func ReconstructSecret(x []*big.Int,y []*big.Int,numparticipants int, p *big.Int) *big.Int {
    
    yval:=big.NewInt(0)
    b:=make([]*big.Int,numparticipants)
    
    //calculate share contribution from each participant's x value (note that this doesn't require interaction with the participants since these values need not be secret)
    for i:=0;i<numparticipants;i++ {
        b[i]=big.NewInt(1)
        for j:=0;j<numparticipants;j++ {
            if i==j {
                continue
            }
            mi:=new(big.Int).ModInverse(new(big.Int).Sub(x[j],x[i]),p)
            bt:=new(big.Int).Mul(x[j],mi)
            bt=new(big.Int).Mod(bt,p)
            b[i]=new(big.Int).Mul(b[i],bt)
            b[i]=new(big.Int).Mod(b[i],p)
       }
    }
    
    //calculate secret using the individual shares
    for i:=0;i<numparticipants;i++ {
        b[i]=new(big.Int).Mod(b[i],p)
        yval=new(big.Int).Add(yval,new(big.Int).Mul(b[i],y[i]))
        yval=new(big.Int).Mod(yval,p)
    }
    secret:=yval
    return secret
}
