defmodule ElixirPoc do
  @input_file "../data-generator/stock_transactions.jsonl"
  @output_file "aapl_sales.jsonl"
  @difficulty_level 100 # Misma dificultad que Go/Python/Node

  def run_benchmark do
    IO.puts("🚀 Iniciando procesamiento con Elixir (Flow / Map-Reduce)...")

    {:ok, output_pid} = File.open(@output_file, [:write, :utf8])
    start_time = System.monotonic_time(:millisecond)

    results =
      @input_file
      |> File.stream!(read_ahead: 100_000)
      |> Flow.from_enumerable(max_demand: 100)
      |> Flow.map(&parse_json/1)
      |> Flow.reject(&is_nil/1)
      |> Flow.map(&heavy_computation/1)
      |> Flow.map(fn record ->
        check_and_write(record, output_pid)
      end)
      # --- FASE 1: Reducción Paralela (Local en cada núcleo) ---
      |> Flow.reduce(fn -> %{count: 0, googl: 0, volume: 0.0, aapl_sales: 0} end, fn record, acc ->
        %{
          count: acc.count + 1,
          googl: acc.googl + (if record["ticker"] == "GOOGL", do: 1, else: 0),
          volume: acc.volume + (record["price"] * record["quantity"]),
          aapl_sales: acc.aapl_sales + (if record["aapl_found"], do: 1, else: 0)
        }
      end)
      # --- FASE 2: Shuffle (Enviar todo a un solo proceso) ---
      |> Flow.partition(stages: 1)
      # --- FASE 3: Reducción Global (Sumar los parciales) ---
      |> Flow.reduce(fn -> %{count: 0, googl: 0, volume: 0.0, aapl_sales: 0} end, fn
        {key, value}, acc ->
          Map.update!(acc, key, &(&1 + value))
        _, acc ->
          acc
      end)
      # Ahora estamos seguros de que Flow emitirá exactamente UN mapa con el total
      |> Enum.to_list()
      |> Map.new()

    print_results(results, start_time, output_pid)
  end
  # --- FUNCIONES AUXILIARES ---

  defp parse_json(line) do
    case Jason.decode(line) do
      {:ok, data} -> data
      _ -> nil
    end
  end

  # Simulación de Proof of Work (Recursivo)
  defp heavy_computation(record) do
    payload = "#{record["txn_id"]}#{record["timestamp"]}"

    # En Elixir, :crypto llama directamente a OpenSSL (C Code), es muy rápido.
    # Usamos Enum.reduce para simular el loop for
    _final_hash = Enum.reduce(1..@difficulty_level, payload, fn _, acc ->
      :crypto.hash(:sha256, acc)
    end)

    record # Retornamos el record original para que siga en la tubería
  end

  # Lógica de filtrado y escritura
  defp check_and_write(record, output_pid) do
    aapl_found = record["ticker"] == "AAPL" and record["operation"] == "SELL"

    if aapl_found do
      # Volvemos a codificar a JSON para escribir (o podríamos pasar la línea original)
      # Para ser justos con Go, escribimos.
      {:ok, line} = Jason.encode(record)
      IO.write(output_pid, line <> "\n")
    end

    # Agregamos una marca al record para saber si fue venta de AAPL en el reduce
    Map.put(record, "aapl_found", aapl_found)
  end

  defp print_results(stats, start_time, output_pid) do
    File.close(output_pid)
    end_time = System.monotonic_time(:millisecond)
    elapsed_sec = (end_time - start_time) / 1000.0

    IO.puts("\n🏁 Procesamiento Terminado")
    IO.puts("⏱️  Tiempo Total: #{Float.round(elapsed_sec, 4)} segundos")
    IO.puts("📄 Líneas procesadas: #{stats.count}")
    IO.puts("🔎 'GOOGL' encontradas: #{stats.googl}")

    # Formateo de moneda
    volume_str = :erlang.float_to_binary(stats.volume, [decimals: 2])
    IO.puts("💰 Volumen Total: $#{volume_str}")
    IO.puts("💾 Ventas de 'AAPL' guardadas: #{stats.aapl_sales}")
    IO.puts("📁 Archivo generado: #{@output_file}")
  end
end
