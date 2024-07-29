package controller

import (
	"math/rand"
	"strconv"
)

func GenerateInt() string {
	num := strconv.Itoa(rand.Intn(900000) + 100000)
	return num
}
