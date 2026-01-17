# 📑 Resumen Ejecutivo: Benchmark de Arquitectura y Concurrencia

**Proyecto:** "The 1 Million Rows Challenge"
**Arquitecto:** Codesandtags

## 1. Objetivo del Estudio

El propósito de esta Prueba de Concepto (PoC) no fue simplemente medir la velocidad sintáctica de diferentes lenguajes, sino **evaluar el comportamiento de distintos Modelos de Concurrencia** frente a una carga de trabajo mixta (IO-Bound y CPU-Bound).

Se buscó responder: _¿Cómo aprovecha cada lenguaje el hardware moderno (Multi-Core) al procesar grandes volúmenes de datos?_

## 2. Definición de la Carga de Trabajo

Procesamiento de un dataset financiero (`stock_transactions.jsonl`) de **1 millón de registros** (aprox. 150MB).

**Las Operaciones:**

1.  **Parsing JSON:** Deserialización de texto a estructuras de memoria (Costo de memoria/parsing).
2.  **Búsqueda (Search):** Filtrado de cadenas (`"GOOGL"`).
3.  **Matemática (Math):** Aritmética de punto flotante y acumulación de estado (`Price * Quantity`).
4.  **Hashing Criptográfico (CPU Stress Test):**
    - _Operación:_ SHA-256 recursivo (100 iteraciones por registro).
    - _Por qué:_ Esta operación fue diseñada para saturar el procesador, obligando a los lenguajes a demostrar sus capacidades de paralelismo real. Sin esto, la prueba hubiera sido dependiente solo de la velocidad del disco SSD.
5.  **I/O Escritura:** Generación de un reporte filtrado para ventas de `"AAPL"`.

## 3. Infraestructura de Prueba

- **Hardware:** Apple MacBook Pro (Chip M3).
- **Arquitectura:** ARM64 con arquitectura híbrida (Performance Cores + Efficiency Cores).
- **Sistema Operativo:** macOS Sonoma.

## 4. Matriz de Lenguajes (Los Contendientes)

Seleccionamos 4 lenguajes que representan los distintos cuadrantes del diseño de sistemas de tipos y modelos de ejecución:

| Lenguaje    | Cuadrante (Tipo/Verificación) | Modelo de Concurrencia         | Filosofía                                                                         |
| :---------- | :---------------------------- | :----------------------------- | :-------------------------------------------------------------------------------- |
| **Python**  | **Dynamic & Strong**          | **Single Thread (GIL)**        | Simplicidad y legibilidad. Ejecución secuencial bloqueante.                       |
| **Node.js** | **Dynamic & Weak**            | **Event Loop (Single Thread)** | I/O no bloqueante. Excelente para redes, pero limitado por un solo hilo para CPU. |
| **Go**      | **Static & Strong**           | **CSP (Goroutines)**           | Paralelismo nativo, memoria compartida via canales, compilado a código máquina.   |
| **Elixir**  | **Dynamic & Strong**          | **Actor Model (BEAM VM)**      | Procesos aislados, inmutabilidad, tolerancia a fallos masiva.                     |

## 5. Resultados del Benchmark (Actualizado)

| Posición | Lenguaje    | Tiempo Total | Factor vs. Ganador | Observación Clave                                                                   |
| :------- | :---------- | :----------- | :----------------- | :---------------------------------------------------------------------------------- |
| 🥇 **1** | **Go**      | **2.36 s**   | 1x                 | Saturación total de CPU (100% uso). Velocidad "Metal-frío".                         |
| 🥈 **2** | **Elixir**  | **17.48 s**  | ~7.4x              | Saturación total de CPU. Overhead natural de la Máquina Virtual (BEAM).             |
| 🥉 **3** | **Python**  | **27.21 s**  | ~11.5x             | Limitado a 1 Core. Rápido en hashing solo porque `hashlib` es C nativo.             |
| 🐢 **4** | **Node.js** | **52.43 s**  | ~22.2x             | Limitado a 1 Core. Sufrió overhead por gestión de objetos/buffers en el Event Loop. |

## 6. Conclusiones Técnicas

1.  **La falacia del "Asíncrono":** Node.js demostró que ser "Non-blocking I/O" no sirve de nada cuando la tarea es intensiva en CPU. El Event Loop se congeló procesando hashes, resultando en el peor tiempo.
2.  **El poder del Hardware:** Solo **Go** y **Elixir** lograron justificar la inversión en un chip M3. Ambos lenguajes activaron automáticamente todos los núcleos (Performance y Efficiency), demostrando escalabilidad horizontal real dentro de una misma máquina.
3.  **Compilado vs. VM:** Aunque Elixir paralelizó perfecto (usando todos los cores), **Go ganó por fuerza bruta** al ser código nativo compilado, eliminando la capa de interpretación de una máquina virtual.
4.  **Veredicto:**
    - Para **Data Pipelines & Number Crunching** de alto rendimiento: **Go** es el rey indiscutible.
    - Para **Sistemas Distribuidos Tolerantes a Fallos** donde la estabilidad importa más que la velocidad cruda: **Elixir** es la opción superior.
    - Para **Scripting rápido y Prototipado**: **Python** sigue siendo imbatible en velocidad de desarrollo (no de ejecución).
