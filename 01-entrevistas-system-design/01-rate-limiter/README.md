# Rate Limiter Distribuido


```mermaid
flowchart TD
    Client[📱 Clientes (Miles de Requests)] --> LB[⚖️ Load Balancer / Nginx]

    subgraph "Cluster de Aplicación (Tu Código)"
        LB --> S1[🖥️ Servidor API 1]
        LB --> S2[🖥️ Servidor API 2]
        LB --> S3[🖥️ Servidor API 3]
    end

    subgraph "Capa de Datos (Shared State)"
        Redis[("🔴 Redis (Master)
        Clave: 'user_123'
        Tokens: 5")]
    end

    S1 -- "1. Evaluar Lua Script" --> Redis
    S2 -- "1. Evaluar Lua Script" --> Redis
    S3 -- "1. Evaluar Lua Script" --> Redis

    Redis -- "2. Permitir/Denegar" --> S1
    Redis -- "2. Permitir/Denegar" --> S2
    Redis -- "2. Permitir/Denegar" --> S3

    style Redis fill:#ff4444,color:white,stroke:#333
    style LB fill:#f9f,stroke:#333
    style Client fill:#fff,stroke:#333
```
