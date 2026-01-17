const fs = require("fs");
const readline = require("readline");
const crypto = require("crypto");
const { performance } = require("perf_hooks");

// Configuración
const INPUT_FILE = "../data-generator/stock_transactions.jsonl";
const OUTPUT_FILE = "aapl_sales.jsonl";
const DIFFICULTY_LEVEL = 100; // Ajusta para simular carga CPU

/**
 * Simula la tarea pesada de CPU (Proof of Work)
 * Replicamos la lógica SHA256 recursivo
 */
function heavyComputation(record) {
  const payload = `${record.txn_id}${record.timestamp}`;
  let currentHash = Buffer.from(payload);

  for (let i = 0; i < DIFFICULTY_LEVEL; i++) {
    // crypto.createHash es la forma nativa en Node
    // .update() carga datos, .digest() ejecuta el hash
    currentHash = crypto.createHash("sha256").update(currentHash).digest();
  }

  return currentHash.toString("hex");
}

async function runBenchmark() {
  console.log(
    "🚀 Iniciando procesamiento con Node.js (Streams + Event Loop)...",
  );

  const startTime = performance.now();

  // Métricas
  let googlCount = 0;
  let totalVolume = 0.0;
  let aaplSalesCount = 0;
  let processedLines = 0;

  // Stream de Escritura
  const outputStream = fs.createWriteStream(OUTPUT_FILE, { encoding: "utf8" });

  // Stream de Lectura con Readline (Interfaz línea por línea)
  const fileStream = fs.createReadStream(INPUT_FILE);
  const rl = readline.createInterface({
    input: fileStream,
    crlfDelay: Infinity,
  });

  // Procesamos línea por línea
  // Nota: En Node, esto pausa el Event Loop en cada línea si la CPU está ocupada
  for await (const line of rl) {
    let record;
    try {
      record = JSON.parse(line);
    } catch (e) {
      continue;
    }

    processedLines++;

    // --- CPU BOUND TASK ---
    // Aquí es donde Node sufrirá igual que Python
    // Mientras esto corre, Node no puede hacer nada más (ni I/O) debido a su Event Loop
    heavyComputation(record);
    // ---------------------

    // 1. Búsqueda
    if (record.ticker === "GOOGL") {
      googlCount++;
    }

    // 2. Matemática
    totalVolume += record.price * record.quantity;

    // 3. Escritura (I/O)
    if (record.ticker === "AAPL" && record.operation === "SELL") {
      // write retorna false si el buffer está lleno (backpressure),
      // pero para este ejemplo simple dejaremos que Node lo gestione.
      outputStream.write(line + "\n");
      aaplSalesCount++;
    }
  }

  outputStream.end();

  const endTime = performance.now();
  const elapsedSeconds = (endTime - startTime) / 1000;

  // Reporte
  console.log("\n🏁 Procesamiento Terminado");
  console.log(`⏱️  Tiempo Total: ${elapsedSeconds.toFixed(4)} segundos`);
  console.log(`📄 Líneas procesadas: ${processedLines}`);
  console.log(`🔎 'GOOGL' encontradas: ${googlCount}`);

  // Formatear moneda en JS es un poco más manual sin Intl en loops rápidos
  const formattedVolume = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(totalVolume);
  console.log(`💰 Volumen Total: ${formattedVolume}`);

  console.log(`💾 Ventas de 'AAPL' guardadas: ${aaplSalesCount}`);
  console.log(`📁 Archivo generado: ${OUTPUT_FILE}`);
}

runBenchmark().catch(console.error);
