```mermaid
graph TD
    subgraph Клиент
        A[Пользователь] --> B[Web Browser]
    end

    subgraph API Gateway
        B -->|HTTP Requests| C[HTTP Server :7777]
        C --> D[Router/ServeMux]
    end

    subgraph Аутентификация
        D -->|/v1/auth/register| E[Register Handler]
        D -->|/v1/auth/login| F[Login Handler]
        E --> G[Auth Service]
        F --> G
        G -->|Validate| H[User Repository]
    end

    subgraph Управление аккаунтом
        D -->|/v1/accounts/*| I[Account Handler]
        I --> J[File Uploader]
        J --> K[Local Storage]
        J --> L[File Repository]
    end

    subgraph База данных
        H --> M[(PostgreSQL)]
        L --> M
    end

    subgraph Middleware
        N[CORS Middleware]
        O[Auth Middleware]
        P[Rate Limiter]
    end

    subgraph Файловая система
        K --> Q[Uploads Directory]
    end

    style A fill:#f9f,stroke:#333,stroke-width:2px
    style M fill:#ccf,stroke:#333,stroke-width:2px
    style Q fill:#cfc,stroke:#333,stroke-width:2px
```