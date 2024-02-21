package kx509

import (
	"math/big"
)

/**
https://patents.google.com/patent/KR20040064780A/ko

P-KCDSASignatureValue ::= SEQUENCE {
	r BIT STRING,
	s INTEGER }
P-KCDSAParameters ::= SEQUENCE {
	p INTEGER, -- odd prime p = 2Jq+1
	q INTEGER, -- odd prime
	g INTEGER, -- generator of order q
	J INTEGER OPTIONAL, -- odd prime
	Seed OCTET STRING OPTIONAL
	Count INTEGER OPTIONAL }
P-KCDSAPublicKey ::= INTEGER -- Public key y
*/

type kcdsaParameters struct {
	P     *big.Int
	Q     *big.Int
	G     *big.Int
	J     *big.Int `asn1:"optional"`
	Seed  []byte   `asn1:"optional"`
	Count int      `asn1:"optional"`
}

/**
type kcdsaPrivateKey struct {
	Version    int
	PrivateKey []byte
	PublicKey  asn1.BitString `asn1:"optional,explicit,tag:0"`
}
*/
