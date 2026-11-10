# Backend - Sistema de Gestión para una Constructora (Go + PostgreSQL + Docker)

## Estructura de Microservicios

- **`services/project-service`** (Puerto `8081` · Base de datos: `projects_db`)
- **`services/inventory-service`** (Puerto `8082` · Base de datos: `inventory_db` · Patrones *Factory Method* y *Abstract Factory*)
- **`services/budget-service`** (Puerto `8083` · Base de datos: `budgets_db` · Cálculo desacoplado y comunicación HTTP con `project-service`)

## Cómo ejecutar con Docker Compose

Desde la carpeta `backend/`:

```bash
cd backend
docker compose up --build
```

Esto levantará automáticamente:
1. Contenedor de **PostgreSQL 16** creando las 3 bases de datos independientes (`projects_db`, `inventory_db`, `budgets_db`) mediante `docker/init.sql`.
2. **`project-service`** en `http://localhost:8081`
3. **`inventory-service`** en `http://localhost:8082`
4. **`budget-service`** en `http://localhost:8083`

## Cómo ejecutar el Frontend (Next.js)

En otra terminal, desde la carpeta `frontend/`:

```bash
cd frontend
npm run dev
```

Abre `http://localhost:3000` en tu navegador.

