package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestStorageTrie(t *testing.T) {
	fmt.Println("RUNNING TEST: TestStorageTrie")
	fmt.Println("IN FILE: storage_proof_test.go")
	
	// slot indexes
	// slot0 := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000") 
		// 0x4e46545475746f7269616c000000000000000000000000000000000000000016
	// slot1 := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001")
		// 0x4e46540000000000000000000000000000000000000000000000000000000006

	nameSlot := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	symbolSlot := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001")
	counterSlot := common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000006")
	ownersSlot := common.FromHex("0x679795a0195a1b76cdebb7c51d74e058aee92919b8c3389af86ef24535e8a28c")
	balancesSlot := common.FromHex("0x211f3d93987b5218a32eac3af2b87ceafa2dad2bbbdbe3688f9c11e352c27cd8")

	// encode values to be stored
	// ownerAddress, err := rlp.EncodeToBytes(common.FromHex("0xde74da73d5102a796559933296c73e7d1c6f37fb"))
	// require.NoError(t, err)

	// lastCompletedMigration, err := rlp.EncodeToBytes(common.FromHex("0x02"))
	// require.NoError(t, err)

	name, err := rlp.EncodeToBytes([]byte("NFTTutorial")) // common.FromHex("0x4e46545475746f7269616c000000000000000000000000000000000000000016")
	require.NoError(t, err)
	symbol, err := rlp.EncodeToBytes([]byte("NFT")) // common.FromHex("0x4e46540000000000000000000000000000000000000000000000000000000006")
	require.NoError(t, err)
	ownerAddress, err := rlp.EncodeToBytes(common.FromHex("0x2813736e6204ee248e79c26de69d49bddbe0f7d0"))
	require.NoError(t, err)
	balances, err := rlp.EncodeToBytes(common.FromHex("0x02")) // should use 0x02 instead of padded hex (identical to the one provided in the "Slot")
	require.NoError(t, err)
	counter, err := rlp.EncodeToBytes(common.FromHex("0x02"))


	fmt.Println("nameSlot: ", nameSlot)
	fmt.Println("symbolSlot: ", symbolSlot)
	fmt.Println("ownerSlot: ", ownersSlot)
	fmt.Println("balanceSlot: ", balancesSlot)
	fmt.Println("counterSlot: ", counterSlot)

	fmt.Println("name: ", name)
	fmt.Println("symbol: ", symbol)
	fmt.Println("ownerAddress ", ownerAddress)
	fmt.Println("balances: ", balances)
	fmt.Println("counter: ", counter)

	fmt.Println("name to hex: ", common.Bytes2Hex(name))
	fmt.Println("symbol to hex: ", common.Bytes2Hex(symbol))
	fmt.Println("ownerAddress to hex: ", common.Bytes2Hex(ownerAddress))
	fmt.Println("balances to hex: ", common.Bytes2Hex(balances))
	fmt.Println("counter to hex: ", common.Bytes2Hex(counter))

	// create a trie and store the key-value pairs, the key needs to be hashed
	trie := NewTrie()
	trie.Put(crypto.Keccak256(nameSlot), name)
	trie.Put(crypto.Keccak256(symbolSlot), symbol)
	trie.Put(crypto.Keccak256(counterSlot), counter)
	trie.Put(crypto.Keccak256(ownersSlot), ownerAddress)
	trie.Put(crypto.Keccak256(balancesSlot), balances)

	// compute the root hash and check if consistent with the storage hash of contract 0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b
	rootHash := trie.Hash()
	storageHash := common.FromHex("0xcf1e4b90f815964e5f79b713232d0cfb7bb54617e7775bededbc4bd9d96c0fad")

	fmt.Println("storageHash: ", fmt.Sprintf("%+x", storageHash))
	fmt.Println("rootHash:", fmt.Sprintf("%+x", rootHash))

	require.Equal(t, storageHash, rootHash)
}

func TestContractStateProof(t *testing.T) {
	fmt.Println("RUNNING TEST: TestContractStateProof")
	fmt.Println("IN FILE: storage_proof_test.go")

	// curl https://eth-mainnet.g.alchemy.com/v2/<API_KEY> \
	//       -X POST \
	//       -H "Content-Type: application/json" \
	//       -d '{"jsonrpc":"2.0","method":"eth_getProof","params":["0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b",["0x0"], "0xA8894B"],"id":1}'

	jsonFile, err := os.Open("storage_proof_slot_0.json")
	require.NoError(t, err)

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	// load into the struct
	var response EthGetProofResponse
	err = json.Unmarshal(byteValue, &response)
	require.NoError(t, err)

	result := response.Result

	account := common.HexToAddress("0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b")
	fmt.Println(fmt.Sprintf("decoded account state data from untrusted source for address %x: balance is %x, nonce is %x, codeHash: %x, storageHash: %x",
		account, result.Balance, result.Nonce, result.CodeHash, result.StorageHash))

	// get the state root hash from etherscan: https://etherscan.io/block/11045195
	stateRootHash := common.HexToHash("0x8c571da4c95e212e508c98a50c2640214d23f66e9a591523df6140fd8d113f29")

	// create a proof trie, and add each node from the account proof
	proofTrie := NewProofDB()
	for _, node := range result.AccountProof {
		proofTrie.Put(crypto.Keccak256(node), node)
	}

	// verify the proof against the stateRootHash
	validAccountState, err := VerifyProof(
		stateRootHash.Bytes(), crypto.Keccak256(account.Bytes()), proofTrie)
	require.NoError(t, err)

	// double check the account state is identical with the account state in the result.
	accountState, err := rlp.EncodeToBytes([]interface{}{
		result.Nonce,
		result.Balance.ToInt(),
		result.StorageHash,
		result.CodeHash,
	})
	require.NoError(t, err)
	require.True(t, bytes.Equal(validAccountState, accountState), fmt.Sprintf("%x!=%x", validAccountState, accountState))

	// now we can trust the data in StorageStateResult
}

func TestContractStorageProofSlot0(t *testing.T) {
	fmt.Println("RUNNING TEST: TestContractStorageProofSlot0")
	fmt.Println("IN FILE: storage_proof_test.go")

	// curl https://eth-mainnet.g.alchemy.com/v2/<API_KEY> \
	//       -X POST \
	//       -H "Content-Type: application/json" \
	//       -d '{"jsonrpc":"2.0","method":"eth_getProof","params":["0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b",["0x0"], "0xA8894B"],"id":1}'

	// Read storage proof
	jsonFile, err := os.Open("storage_proof_slot_0.json")
	require.NoError(t, err)

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	// parse the proof
	var response EthGetProofResponse
	err = json.Unmarshal(byteValue, &response)
	require.NoError(t, err)

	result := response.Result

	// the storage hash and the proof is the data to be verified
	storageHash := result.StorageHash
	storageProof := result.StorageProof[0]

	// encode the key-value pair
	key := common.LeftPadBytes(storageProof.Key, 32)
	value, err := rlp.EncodeToBytes(storageProof.Value)
	require.NoError(t, err)

	// build a trie with the nodes in the proof
	proofTrie := NewProofDB()
	for _, node := range storageProof.Proof {
		proofTrie.Put(crypto.Keccak256(node), node)
	}

	// verify the proof
	verified, err := VerifyProof(
		storageHash.Bytes(), crypto.Keccak256(key), proofTrie)
	require.NoError(t, err)

	// confirm the value from the proof is consistent with the reported value
	require.True(t, bytes.Equal(verified, value), fmt.Sprintf("%x != %x", verified, value))
}

func TestContractStorageProofSlot1(t *testing.T) {
	fmt.Println("RUNNING TEST: TestContractStorageProofSlot1")
	fmt.Println("IN FILE: storage_proof_test.go")

	// curl https://eth-mainnet.g.alchemy.com/v2/<API_KEY> \
	//       -X POST \
	//       -H "Content-Type: application/json" \
	//       -d '{"jsonrpc":"2.0","method":"eth_getProof","params":["0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b",["0x1"], "0xA8894B"],"id":1}'

	jsonFile, err := os.Open("storage_proof_slot_1.json")
	require.NoError(t, err)

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	fmt.Println("loaded eip1186_proof")

	// load into the struct
	var response EthGetProofResponse
	err = json.Unmarshal(byteValue, &response)
	require.NoError(t, err)

	result := response.Result

	storageHash := result.StorageHash
	storageProof := result.StorageProof[0]
	value, err := rlp.EncodeToBytes(storageProof.Value)
	require.NoError(t, err)
	// 0x0000000000000000000000000000000000000000000000000000000000000000
	key := common.LeftPadBytes(storageProof.Key, 32)

	proofTrie := NewProofDB()
	for _, node := range storageProof.Proof {
		proofTrie.Put(crypto.Keccak256(node), node)
	}

	verified, err := VerifyProof(
		storageHash.Bytes(), crypto.Keccak256(key), proofTrie)

	require.NoError(t, err)
	require.True(t, bytes.Equal(verified, value), fmt.Sprintf("%x != %x", verified, value))
}
