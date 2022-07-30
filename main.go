package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blockchainAddress = blockchainAddress
	block := new(Block)
	blockchain.AddBlock(0, block.Hash())

	return blockchain
}

func (blockchain *Blockchain) AddBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockchain.transactionPool)
	blockchain.chain = append(blockchain.chain, block)
	blockchain.transactionPool = []*Transaction{}
	return block
}

func (blockchain *Blockchain) AddTransaction(from string, to string, amount float32) {
	transaction := NewTransaction(from, to, amount)
	blockchain.transactionPool = append(blockchain.transactionPool, transaction)
}

func (blockchain *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockchain.transactionPool {
		transactions = append(transactions, NewTransaction(transaction.from, transaction.to, transaction.amount))
	}

	return transactions
}

func (blockchain *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())

	return guessHashStr[:difficulty] == zeros
}

func (blockchain *Blockchain) ProofOfWork() int {
	transactions := blockchain.CopyTransactionPool()
	previousHash := blockchain.LastBlock().Hash()

	nonce := 0
	for !blockchain.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func (blockchain *Blockchain) Mining() bool {
	blockchain.AddTransaction(MINING_SENDER, blockchain.blockchainAddress, MINING_REWARD)
	nonce := blockchain.ProofOfWork()
	previousHash := blockchain.LastBlock().Hash()
	blockchain.AddBlock(nonce, previousHash)

	fmt.Println("action=mining, message=success")

	return true
}

func (blockchain *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0

	for _, block := range blockchain.chain {
		for _, transaction := range block.transactions {

			amount := transaction.amount

			if blockchainAddress == transaction.to {
				totalAmount += amount
			}

			if blockchainAddress == transaction.from {
				totalAmount -= amount
			}
		}
	}

	return totalAmount
}

func (blockchain *Blockchain) Show() {
	for index, block := range blockchain.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), index, strings.Repeat("=", 25))
		block.Show()
	}
}

func (blockchain *Blockchain) LastBlock() *Block {
	return blockchain.chain[len(blockchain.chain)-1]
}

type Block struct {
	timestamp    int64          `json:"timestamp"`
	nonce        int            `json:"nonce"`
	previousHash [32]byte       `json:"previous_hash"`
	transactions []*Transaction `json:"transactions"`
}

func NewBlock(nonce int, previousHash [32]byte, transaction []*Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transaction,
	}
}

func (block *Block) Hash() [32]byte {
	marshal, marshalErr := json.Marshal(block)

	if marshalErr != nil {
		fmt.Errorf("error generating hash block")
	}

	return sha256.Sum256(marshal)
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: block.previousHash,
		Transactions: block.transactions,
	})
}

func (block *Block) Show() {
	fmt.Printf("nonce: %d\n", block.nonce)
	fmt.Printf("previous Hash: %x\n", block.previousHash)
	fmt.Printf("timestamp: %d\n", block.timestamp)

	for index, transaction := range block.transactions {
		fmt.Printf("%s transaction %d %s\n", strings.Repeat("-", 22), index, strings.Repeat("-", 22))
		transaction.Show()
	}
}

type Transaction struct {
	from   string
	to     string
	amount float32
}

func NewTransaction(from string, to string, amount float32) *Transaction {
	return &Transaction{
		from:   from,
		to:     to,
		amount: amount,
	}
}

func (transaction *Transaction) Show() {
	fmt.Printf("	from: %s\n", transaction.from)
	fmt.Printf("	to: %s\n", transaction.to)
	fmt.Printf("	amount: %v\n", transaction.amount)
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float32 `json:"amount"`
	}{
		From:   transaction.from,
		To:     transaction.to,
		Amount: transaction.amount,
	})
}

func main() {
	myblockChainAddress := "my-blockchain-address"
	blockchain := NewBlockchain(myblockChainAddress)

	blockchain.AddTransaction("Alice", "Bob", 100)
	blockchain.AddTransaction("Bob", "Alice", 50)
	blockchain.Mining()

	blockchain.AddTransaction("Gabriel", "Jessica", 250)
	blockchain.Mining()

	blockchain.Show()

	fmt.Printf("Blockchain %.1f\n", blockchain.CalculateTotalAmount(myblockChainAddress))
	fmt.Printf("Gabriel %.1f\n", blockchain.CalculateTotalAmount("Gabriel"))
	fmt.Printf("Jessica %.1f\n", blockchain.CalculateTotalAmount("Jessica"))
}
