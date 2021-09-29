package api

import (
	"crypto/sha512"
	"encoding/base64"
	"sync"
	"time"
)

//Client struct to hold the password map
type Client struct {
	hashes sync.Map
}

// Get an instance of the password handler client
func New() *Client {
	return &Client{}
}

//Handler for hashing a password after the request
func (c *Client) HandlePassword(wg *sync.WaitGroup, password string, count int) <-chan int {
	ch := make(chan int)

	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 5)
		hashedPwd := hashPassword([]byte(password))
		c.Save(count, hashedPwd)
		wg.Done()
	}()

	return ch
}

//Save hashed password to map
func (c *Client) Save(id int, pwd string) {
	c.hashes.Store(id, pwd)
}

//Get the hashed password by id
func (c *Client) Get(id int) string {
	val, ok := c.hashes.Load(id)
	if !ok {
		//better way to handle a bad load....
		return ""
	}
	return val.(string)
}

//Hash a password (sha512)
func hashPassword(password []byte) string {
	hash := sha512.Sum512(password)
	//convert 64byte to byte
	return base64.StdEncoding.EncodeToString(hash[:])
}
