package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"math/rand"
	"encoding/json"
)

// constants for data generation
const (
	TotalRecords = 1_000_000
	FileName	 = "stock_transactions.jsonl"
)

// Transaction structure
type Transaction struct {
	TxnID	 string  `json:"txn_id"`
	Timestamp string `json:"timestamp"`
	Ticket string `json:"ticker"`
	Operation string `json:"operation"` // BUY | SELL
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
	Status string `json:"status"`
}

// Base tickers
var tickers = []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META", "BRK.B", "JPM", "V"}
var operations = []string{"BUY", "SELL"}
var statuses = []string{"COMPLETED", "COMPLETED", "COMPLETED", "PENDING", "REJECTED"}


// Here goes the magic
func main() {
	fmt.Printf("🚀 Starting data generation for [%d] records...  \n", TotalRecords)
	start := time.Now()

	file, err := os.Create(FileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Buffered writer (in memory) for performance, because writing line by line is slow.
	writer := bufio.NewWriter(file)

	// seed to make sure we have random data every time we run the generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < TotalRecords; i++ {
		// Generate random transactions
		txn := Transaction{
			TxnID:     fmt.Sprintf("txn-%d-%d", rng.Int63(), i),
			Timestamp: randomDate(rng),
			Ticket:    tickers[rng.Intn(len(tickers))],
			Operation: operations[rng.Intn(len(operations))],
			Price:     10.0 + rng.Float64()*(1000.0-10.0), // Price between 10.0 and 1000.0
			Quantity:  rng.Intn(100) + 1,
			Status:    statuses[rng.Intn(len(statuses))],
		}

		// Serialize to JSON and write to file
		jsonData, _ := json.Marshal(txn)

		// Write in buffer
		_, error := writer.Write(jsonData)
		if error != nil {
			panic(error)
		}

		writer.WriteByte('\n') // New line for JSONL format

		// Progress indicator
		if (i+1)%100_000 == 0 {
			fmt.Printf(" 📝 Generated %d records...\n", i+1)
		}
	}

	// Flush the buffer to ensure all data is written to file
	defer writer.Flush()
	elapsed := time.Since(start)
	fmt.Printf(" 💾 Data written to file: %s\n", FileName)
	fmt.Printf(" ⏱️  Time taken: %s\n", elapsed)

	fmt.Printf("\n ✨ Data generation completed successfully! \n")
}

// Helper function to generate random date strings
func randomDate(rng *rand.Rand) string {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	delta := end - start

	sec := rng.Int63n(delta) + start
	return time.Unix(sec, 0).Format(time.RFC3339)
}