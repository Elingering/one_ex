package models

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GeneraOrderNo() string {
	str := time.Now().Format("20060102150405")
	rand.Seed(time.Now().UnixNano())
	body := fmt.Sprintf("%d", rand.Intn(10000000))
	return str + body
}