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
        D -->|/api/v1/auth/register| E[Register Handler]
        D -->|/api/v1/auth/login| F[Login Handler]
        E --> G[Auth Service]
        F --> G
        G -->|Validate| H[User Repository]
    end

    subgraph Управление пользователями
        D -->|/api/v1/users/{id}| AA[User Handler]
        AA --> AB[User Repository]
        AA -->|DTO| AC[GetUserResponse]
    end

    subgraph Управление аккаунтом
        D -->|/api/v1/account/*| I[Account Handler]
        I --> J[File Uploader]
        J --> K[Local Storage]
        J --> L[File Repository]
    end

    subgraph Управление концертами
        D -->|/admin/v1/concerts/*| R[Concert Handler]
        R --> S[Concert Service]
        S --> T[Concert Repository]
        S --> U[Venue Repository]
        S --> V[Band Repository]
        R -->|DTO| RA[ConcertResponse]
    end

    subgraph Управление местами проведения
        D -->|/admin/v1/venues/*| W[Venue Handler]
        W --> X[Venue Service]
        X --> U
        W -->|DTO| WA[VenueResponse]
    end

    subgraph Управление группами
        D -->|/admin/v1/bands/*| Y[Band Handler]
        Y --> Z[Band Service]
        Z --> V
        Y -->|DTO| YA[BandResponse]
    end

    subgraph Модели/DTO
        AC --> AD[Data Transfer Objects]
        RA --> AD
        WA --> AD
        YA --> AD
    end

    subgraph База данных
        H --> M[(PostgreSQL)]
        AB --> M
        L --> M
        T --> M
        U --> M
        V --> M
    end

    subgraph Middleware
        N[CORS Middleware]
        O[Auth Middleware]
        P[Admin Middleware]
    end

    subgraph Файловая система
        K --> Q[Uploads Directory]
    end

    style A fill:#f9f,stroke:#333,stroke-width:2px
    style M fill:#ccf,stroke:#333,stroke-width:2px
    style Q fill:#cfc,stroke:#333,stroke-width:2px
    style R fill:#fcc,stroke:#333,stroke-width:2px
    style W fill:#fcc,stroke:#333,stroke-width:2px
    style Y fill:#fcc,stroke:#333,stroke-width:2px
    style AD fill:#ffc,stroke:#333,stroke-width:2px
```