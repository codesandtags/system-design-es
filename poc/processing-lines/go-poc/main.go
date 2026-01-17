package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// Configuración
const (
	InputFile       = "../data-generator/stock_transactions.jsonl"
	OutputFile      = "aapl_sales.jsonl"
	DifficultyLevel = 100 // Misma dificultad que Python/Node
)

// Estructura de datos
type Transaction struct {
	TxnID     string  `json:"txn_id"`
	Timestamp string  `json:"timestamp"`
	Ticker    string  `json:"ticker"`
	Operation string  `json:"operation"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Status    string  `json:"status"`
}

// Result agrupa los cálculos de cada worker para enviarlos al main
type Result struct {
	Volume      float64
	GooglFound  int
	AaplSale    bool
	RawLine     string // Para escribir en el archivo si es necesario
}

// heavyComputation simula la carga de CPU (Proof of Work)
// Replica exactamente la lógica de Python/Node: SHA256 recursivo sobre bytes
func heavyComputation(txn Transaction) string {
	payload := fmt.Sprintf("%s%s", txn.TxnID, txn.Timestamp)
	currentHash := []byte(payload)

	for i := 0; i < DifficultyLevel; i++ {
		// Sum256 retorna un array [32]byte
		hash := sha256.Sum256(currentHash)
		// Convertimos el array a slice para la siguiente iteración
		currentHash = hash[:]
	}

	return hex.EncodeToString(currentHash)
}

// worker es la función que ejecutarán múltiples goroutines en paralelo
func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range jobs {
		var txn Transaction
		// Parsear JSON (CPU bound)
		if err := json.Unmarshal([]byte(line), &txn); err != nil {
			continue
		}

		// --- CPU BOUND TASK (Aquí es donde Go brilla) ---
		// Esta función pesada se ejecutará en paralelo en todos los cores
		heavyComputation(txn)
		// -----------------------------------------------

		res := Result{
			Volume: txn.Price * float64(txn.Quantity),
		}

		// Búsqueda
		if txn.Ticker == "GOOGL" {
			res.GooglFound = 1
		}

		// Filtrado para escritura
		if txn.Ticker == "AAPL" && txn.Operation == "SELL" {
			res.AaplSale = true
			res.RawLine = line
		}

		results <- res
	}
}

// writer maneja la escritura a disco en una goroutine separada
// Esto evita lock contention en el archivo
func writer(filename string, writeChan <-chan string, done chan<- bool) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for line := range writeChan {
		writer.WriteString(line)
		writer.WriteByte('\n')
	}
	done <- true
}

func main() {
	numWorkers := runtime.NumCPU()
	fmt.Printf("🚀 Iniciando procesamiento con Go (%d Workers / Goroutines)...\n", numWorkers)

	start := time.Now()

	// Canales
	jobs := make(chan string, 100)
	results := make(chan Result, 100)
	writeChan := make(chan string, 100)
	writeDone := make(chan bool)

	var wg sync.WaitGroup

	// 1. Lanzar Escritor
	go writer(OutputFile, writeChan, writeDone)

	// 2. Lanzar Workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// 3. Lanzar Lector
	go func() {
		file, err := os.Open(InputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
		close(jobs)
	}()

	// 4. Goroutine para cerrar RESULTS (SOLO results)
	go func() {
		wg.Wait()      // Esperar a que los workers terminen
		close(results) // Ya nadie más enviará resultados
		// ⚠️ NO cerramos writeChan aquí todavía, porque el main podría necesitarlo
	}()

	// 5. Loop Principal (Main)
	var (
		processedLines int
		googlCount     int
		totalVolume    float64
		aaplSalesCount int
	)

	// Procesar resultados hasta que se vacíe el canal results
	for res := range results {
		processedLines++
		totalVolume += res.Volume
		googlCount += res.GooglFound

		if res.AaplSale {
			aaplSalesCount++
			// Ahora es seguro: writeChan sigue abierto mientras estemos en este loop
			writeChan <- res.RawLine
		}
	}

	// 6. AHORA sí cerramos el canal de escritura
	// Porque ya salimos del loop y estamos seguros de que no enviaremos nada más
	close(writeChan)

	// Esperar a que el escritor termine de guardar en disco
	<-writeDone

	elapsed := time.Since(start)

	fmt.Println("\n🏁 Procesamiento Terminado")
	fmt.Printf("⏱️  Tiempo Total: %.4f segundos\n", elapsed.Seconds())
	fmt.Printf("📄 Líneas procesadas: %d\n", processedLines)
	fmt.Printf("🔎 'GOOGL' encontradas: %d\n", googlCount)
	fmt.Printf("💰 Volumen Total: $%.2f\n", totalVolume)
	fmt.Printf("💾 Ventas de 'AAPL' guardadas: %d\n", aaplSalesCount)
	fmt.Printf("📁 Archivo generado: %s\n", OutputFile)
}