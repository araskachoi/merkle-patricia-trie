package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	fmt.Println("RUNNING TEST: TestBytes")
	fmt.Println("IN FILE: bytes_test.go")
	require.Equal(t, "313233", fmt.Sprintf("%x", []byte("123"))) // convert to hex
	require.Equal(t, []byte{49, 50, 51}, []byte("123"))          // convert to decimal
}
