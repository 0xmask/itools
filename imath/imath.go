package imath

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

// RandStr 获取随机字符串
func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Decimal 保留n位小数
func Decimal(value float64, length int) (res float64) {
	switch length {
	case 1:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", value), 64)
	case 2:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	case 3:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	case 4:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	case 5:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", value), 64)
	case 6:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", value), 64)
	default:
		res = 0
	}
	return res
}
