package utils

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	c := getConfig("../config.yaml")
	fmt.Println(c)
}
