package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"stupidauth/repos"
)

func SetKey(username string) error {
	k, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	return repos.SetKey(username, k)
}
