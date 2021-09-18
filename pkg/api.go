package api

import (
	"crypto/sha512"
	"encoding/base64"
	"sync"
	"time"
)

type Client interface {
	UpdatePasswordCount() int
	Save(id int, pwd string)
	HandlePassword(wg *sync.WaitGroup, password string, count int) <-chan int
	Get(id int) string
}

type client struct {
	hashes sync.Map
	count  int
}

func New() Client {
	return &client{
		count: 0,
	}
}

func (c *client) UpdatePasswordCount() int {
	c.count = c.count + 1
	return c.count
}

/**
 * Handler for hashing a password after the request
 */
func (c *client) HandlePassword(wg *sync.WaitGroup, password string, count int) <-chan int {
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

/*
 * Save hashed password to map
 *
 * @var		c	*clien
 * @global
 */
func (c *client) Save(id int, pwd string) {
	c.hashes.Store(id, pwd)
}

/*
 * Get the hashed password by id
 *
 * @var		c	*clien
 * @global
 * @return string
 */
func (c *client) Get(id int) string {
	val, ok := c.hashes.Load(id)
	if !ok {
		//better way to handle a bad load....
		return ""
	}
	return val.(string)
}

/**
 * Hash a password (sha512)
 *
 * @global
 * @param	password	byte
 * @return	string
 */
func hashPassword(password []byte) string {
	hash := sha512.Sum512(password)
	//convert 64byte to byte
	return base64.StdEncoding.EncodeToString(hash[:])
}
