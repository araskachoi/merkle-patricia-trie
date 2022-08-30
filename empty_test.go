package main

import (
	"testing"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestEmptyNodeHash(t *testing.T) {
	fmt.Println("RUNNING TEST: TestEmptyNodeHash")
	fmt.Println("IN FILE: empty_test.go")

	emptyRLP, err := rlp.EncodeToBytes(EmptyNodeRaw)
	require.NoError(t, err)
	require.Equal(t, EmptyNodeHash, Keccak256(emptyRLP))
}
