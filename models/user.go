package models

import (
	"errors"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("Invalid login")
)

func RegisterUser(username string, password string) error {
	cost := bcrypt.DefaultCost
	//convert the password from string to byte
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	//Add the user to redis
	//the zero tells the set method that the key should not expire
	//the return will return either an error or nil
	return client.Set("user:"+username, hash, 0).Err()
}

func AuthenticateUser(username, password string) error {
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}
	//check if the hash the user entered matched with the one stored:
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return ErrInvalidLogin
	}
	return nil
}
