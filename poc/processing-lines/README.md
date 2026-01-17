# ⚡ 1 Million Rows Challenge: Language Performance PoC

Este repositorio contiene una Prueba de Concepto (PoC) diseñada para comparar el rendimiento, el consumo de recursos y la ergonomía del código de diferentes lenguajes de programación al procesar grandes volúmenes de datos.

El objetivo no es solo ver "cuál es más rápido", sino entender **cómo maneja cada lenguaje la concurrencia, el uso de memoria y la arquitectura I/O** frente a una tarea intensiva.

## 🎯 El Objetivo

Procesar un archivo **JSONL (JSON Lines) de 1 millón de registros** (aprox. 150MB - 200MB) simulando transacciones de bolsa, realizando operaciones mixtas de CPU e I/O.

## 📂 El Dataset

El archivo `stock_transactions.jsonl` es generado sintéticamente y contiene 1,000,000 de líneas. Cada línea es un objeto JSON con esta estructura:

```json
{
  "txn_id": "txn-123456789",
  "timestamp": "2025-05-20T10:00:00Z",
  "ticker": "AAPL",
  "operation": "BUY",
  "price": 150.5,
  "quantity": 10,
  "status": "COMPLETED"
}
```

## 🛠 Las 3 Operaciones (The Benchmark)

Cada implementación debe realizar estas tres tareas en una sola ejecución (single pass si es posible o optimizado según el lenguaje:

### 1. 🔎 Búsqueda (Search)

- **Tarea**: Contar cuántas veces aparece el Ticker GOOGL en el archivo.

Lo que prueba: Velocidad de parsing y comparación de strings.

### 2. 🧮 Matemática (CPU Bound)

- **Tarea**: Calcular el Volumen Total Transaccional en USD de todo el archivo.
- **Fórmula**: Sumatoria de (Price \* Quantity) de cada registro.
- **Lo que prueba**: Manejo de tipos numéricos (Floats), aritmética y acumulación de estado (concurrencia segura).

### 3. 💾 I/O y Filtrado (Write)

- **Tarea**: Filtrar todas las transacciones donde el Ticker sea AAPL (Apple) y la operación sea SELL. Escribir estos registros en un nuevo archivo llamado aapl_sales.jsonl.
- **Lo que prueba**: Escritura en disco eficiente y lógica condicional compuesta.

## Estructura del proyecto

```
├── data-generator/      # Script en Go para generar el dataset
├── python-poc/          # Implementación en Python (Single Threaded / Baseline)
├── javascript-poc/      # Implementación en Node.js (Streams / Event Loop)
├── go-poc/              # Implementación en Go (Goroutines / Worker Pool)
├── elixir-poc/          # Implementación en Elixir (BEAM / Actor Model)
└── stock_transactions.jsonl # (Generado localmente, gitignored)
```

## 🚀 Cómo Ejecuta

```
cd data-generator
go run generator.go
mv stock_transactions.jsonl ../
```

## Resultados

### Python PoC

**Expectativas**

Python es un lenguaje increíble, pero aquí veremos su talón de Aquiles: el Global Interpreter Lock (GIL) y el costo de ser un lenguaje interpretado.

**Estrategia: Streaming (Lazy Evaluation)**

No vamos a cargar el archivo en memoria (read()). Vamos a iterar línea por línea. Esto mantiene la memoria RAM baja (O(1)), aunque el CPU sufrirá parseando JSON un millón de veces secuencialmente.

```sh
🚀 Iniciando procesamiento con Python (Single Threaded)...

🏁 Procesamiento Terminado
⏱️  Tiempo Total: 1.5132 segundos
📄 Líneas procesadas: 1000000
🔎 'GOOGL' encontradas: 99677
💰 Volumen Total: $25,475,929,343.77
💾 Ventas de 'AAPL' guardadas: 50156
📁 Archivo generado: aapl_sales.jsonl
```
