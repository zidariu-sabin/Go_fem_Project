package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuth = "authentification"
)

type Token struct {
	PlainText string `json:"token"`
	Hash      []byte `json:"-"`
	//the user id does not need to be known
	UserID int       `json:"-"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"-"`
}

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	emptyBytes := make([]byte, 32)

	//filles empty slice of bytes read with random characters
	//if this returns an error, there is an internal problem with the rand library for the system
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	//if a value is sent that does not fully meet what the encoder is, the rest is padded with "=" signs
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)

	//using the standard documented crypting lbrary of go to avoid creating a hashing algoritm ourselves
	//we pass the created token as a slice and encrypt it
	hash := sha256.Sum256([]byte(token.PlainText))

	token.Hash = hash[:]

	return token, nil
}
