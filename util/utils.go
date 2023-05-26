package util

import (
	"math/rand"
	"time"
	"unsafe"
)

func GetYmd() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func GetYmdHms() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

func MakeRamdomString(n int) string {
	const (
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}