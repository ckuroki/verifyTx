package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	RPC_URL = "https://mainnet.infura.io/v3/50b6673dc48443e59047246df462902c"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verify_tc <txHash>")
		os.Exit(1)
	}
	txHash := common.HexToHash(os.Args[1])

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatalf("dial:", err.Error())
	}

	// get transaction by hash
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatalf("retrieving tx by hash:", err.Error())
	}
	if isPending {
		log.Fatalf("tx is in pending state: %s\n", txHash)
	}

	// get tx receipt
	txReceipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatalf("cannot get tx receipt: %s\n", txHash)
	}

	// Get target block
	block, err := client.BlockByNumber(context.Background(), txReceipt.BlockNumber)
	if err != nil {
		log.Fatal("cannot get block [%v] : %s ", txReceipt.BlockNumber, err.Error())
	}

	fmt.Printf("Found : %d transactions on block : %v\n", len(block.Transactions()), txReceipt.BlockNumber)

	// Now generate merkle trie
	db := rawdb.NewMemoryDatabase()
	rootHash := common.Hash{}
	merkleTrie, err := trie.New(rootHash, rootHash, trie.NewDatabase(db))
	if err != nil {
		log.Fatal("creating trie : %s ", err.Error())
	}

	txIndex := -1
	for ix, transaction := range block.Transactions() {
		k, err := rlp.EncodeToBytes(uint(ix))
		if err != nil {
			log.Fatal("encoding key [%d] : %s ", ix, err.Error())
		}
		v, err := rlp.EncodeToBytes(transaction)
		if err != nil {
			log.Fatal("encoding transaction [%d] - %v : %s ", ix, transaction, err.Error())
		}
		if txHash == transaction.Hash() {
			txIndex = ix
		}
		merkleTrie.Update(k, v)
	}
	// check txHash index on the trie
	if txIndex < 0 {
		log.Fatalf("txHash not in block transactions\n")
	}
	// get tx merkle proof
	proof := memorydb.New()
	proveKey, err := rlp.EncodeToBytes(uint(txIndex))
	if err != nil {
		log.Fatal("encoding prove key [%d] : %s ", txIndex, err.Error())
	}
	err = merkleTrie.Prove(proveKey, 0, proof)
	if err != nil {
		log.Fatalf("Failed to prove the node [%d] : %v", txIndex, err)
	}
	// Verify proof
	verifyValue, err := trie.VerifyProof(block.Header().TxHash, proveKey, proof)
	if err != nil {
		log.Fatalf("Failed to verify proof : %v", err)
	}
	// compare verified value with transaction
	transaction := block.Transaction(txHash)
	txValue, err := rlp.EncodeToBytes(transaction)
	if err != nil {
		log.Fatal("encoding transaction %v : %s ", transaction, err.Error())
	}
	res := bytes.Compare(verifyValue, txValue)
	if res != 0 {
		fmt.Printf("Verification failed")
		os.Exit(1)
	}
	if block.Header().TxHash == merkleTrie.Hash() {
		fmt.Printf("%v : %v \n", block.Header().TxHash, merkleTrie.Hash())
	} else {
		fmt.Printf("Hashes doesnt match")
	}

	fmt.Println("Verification successful")
}
