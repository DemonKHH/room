package util

import (
	"fmt"
	"math/rand"
)

func GetRandomAvator() string {
	return fmt.Sprintf("%d", rand.Intn(100))
}
