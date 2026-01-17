import json
import time
import hashlib

# Configuración
INPUT_FILE = "../data-generator/stock_transactions.jsonl"
OUTPUT_FILE = "aapl_sales.jsonl"
DIFFICULTY_LEVEL = 100  # Ajusta para simular carga CPU

def heavy_computation(record):
    """
    Simula una tarea pesada de CPU (Proof of Work).
    Replica exactamente la lógica de Node.js: SHA256 recursivo.
    """
    payload = f"{record['txn_id']}{record['timestamp']}"
    current_hash = payload.encode()

    for _ in range(DIFFICULTY_LEVEL):
        # hashlib.sha256().digest() retorna bytes
        # Esto mantiene el loop eficiente pero intenso en CPU
        current_hash = hashlib.sha256(current_hash).digest()

    return current_hash.hex()

def run_benchmark():
    print(f"🚀 Iniciando procesamiento con Python (Single Threaded + GIL)...")

    # Métricas
    googl_count = 0
    total_volume = 0.0
    aapl_sales_count = 0
    processed_lines = 0

    start_time = time.time()

    # Abrimos archivos (Streaming I/O)
    # Python usa iteradores perezosos por defecto al leer archivos
    with open(INPUT_FILE, 'r', encoding='utf-8') as infile, \
         open(OUTPUT_FILE, 'w', encoding='utf-8') as outfile:

        for line in infile:
            try:
                record = json.loads(line)
            except json.JSONDecodeError:
                continue

            processedLines = processed_lines + 1

            # --- CPU BOUND TASK ---
            # NOTA DE ARQUITECTURA:
            # Aquí Python sufre por el GIL (Global Interpreter Lock).
            # Al igual que Node bloquea su Event Loop, Python bloquea el hilo principal.
            # Ningún otro hilo (si los hubiera) podría ejecutarse en paralelo real aquí.
            heavy_computation(record)
            # ---------------------

            # 1. Búsqueda
            if record['ticker'] == 'GOOGL':
                googl_count += 1

            # 2. Matemática
            total_volume += record['price'] * record['quantity']

            # 3. Escritura (I/O)
            if record['ticker'] == 'AAPL' and record['operation'] == 'SELL':
                outfile.write(line)
                aapl_sales_count += 1

            processed_lines = processedLines

    end_time = time.time()
    elapsed_time = end_time - start_time

    # Reporte
    print(f"\n🏁 Procesamiento Terminado")
    print(f"⏱️  Tiempo Total: {elapsed_time:.4f} segundos")
    print(f"📄 Líneas procesadas: {processed_lines}")
    print(f"🔎 'GOOGL' encontradas: {googl_count}")
    print(f"💰 Volumen Total: ${total_volume:,.2f}") # Formato moneda nativo de f-string
    print(f"💾 Ventas de 'AAPL' guardadas: {aapl_sales_count}")
    print(f"📁 Archivo generado: {OUTPUT_FILE}")

if __name__ == "__main__":
    run_benchmark()