package main

import (
	"encoding/hex"
	"testing"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestNull(t *testing.T) {
	fmt.Println("RUNNING TEST: TestNull")
	fmt.Println("IN FILE: crypto_test.go")
	null := []byte{}
	require.Equal(t, "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		hex.EncodeToString(Keccak256(null)))
}

func TestNullRLP(t *testing.T) {
	fmt.Println("RUNNING TEST: TestNullRLP")
	fmt.Println("IN FILE: crypto_test.go")
	emptyRLP, err := rlp.EncodeToBytes([]byte{})
	require.NoError(t, err)

	require.Equal(t, "56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		hex.EncodeToString(Keccak256(emptyRLP)))
}

func TestEmptyArrayRLP(t *testing.T) {
	fmt.Println("RUNNING TEST: TestEmptyArrayRLP")
	fmt.Println("IN FILE: crypto_test.go")
	emptyArray := make([][]byte, 0)
	emptyArrayRLP, err := rlp.EncodeToBytes(emptyArray)
	require.NoError(t, err)
	require.Equal(t, "1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		hex.EncodeToString(Keccak256(emptyArrayRLP)))
}
