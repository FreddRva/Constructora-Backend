# DOCUMENTACIÓN TÉCNICA — SISTEMA DE GESTIÓN CONSTRUCTORA ERP

**Arquitectura de microservicios en Go con los patrones creacionales Factory Method y Abstract Factory, frontend moderno en Next.js y módulos de negocio con montos en soles peruanos (PEN · `S/`).**

> **Documento externo al sistema.** Por decisión de diseño, la información de arquitectura, patrones de diseño y stack tecnológico no se muestra dentro de la aplicación: el usuario final solo ve la interfaz de negocio (Panel de Operaciones + 5 módulos CRUD). Toda la documentación técnica reside exclusivamente en este documento.

---

## Tabla de Contenido

1. [Alcance](#1-alcance)
2. [Arquitectura aplicada en cada parte del sistema](#2-arquitectura-aplicada-en-cada-parte-del-sistema)
3. [Microservicios](#3-microservicios)
4. [Patrones de diseño (Backend — Go)](#4-patrones-de-diseño-backend--go)
   - [4.1 Factory Method](#41-factory-method)
   - [4.2 Abstract Factory](#42-abstract-factory)
   - [4.3 Relación entre ambos patrones](#43-relación-entre-ambos-patrones)
5. [Frontend — Arquitectura por Capas y Componentes](#5-frontend--arquitectura-por-capas-y-componentes)
6. [Stack Tecnológico por capa](#6-stack-tecnológico-por-capa)
7. [Flujo de una petición CRUD (ejemplo: crear ítem de inventario)](#7-flujo-de-una-petición-crud-ejemplo-crear-ítem-de-inventario)
8. [Convención de moneda y Despliegue en AWS (RDS + S3)](#8-convención-de-moneda-y-despliegue-en-aws-rds--s3)

---

## 1. Alcance

El sistema de gestión para la constructora cubre **cinco módulos CRUD**. Todos los registros y cálculos financieros se expresan en **soles peruanos (`PEN · S/`)**.

| Módulo                            | Descripción                                                                                  | Moneda / Unidad    |
| :-------------------------------- | :------------------------------------------------------------------------------------------- | :----------------- |
| **Proyectos (`proyectos`)**       | Obras residenciales, comerciales e industriales con control de plazos, avance (%) y planos   | Avance (%)         |
| **Inventario (`inventario`)**     | Materiales consumibles y herramientas/maquinaria con control de stock y condición            | Stock por unidad   |
| **Presupuestos (`presupuestos`)** | Presupuestos de obra con costos de materiales, mano de obra, equipos y costo total calculado | Soles (`PEN · S/`) |
| **Personal (`personal`)**         | Ingenieros residentes, maestros de obra y operarios asignados a proyectos con jornal diario  | Soles (`PEN · S/`) |
| **Compras (`compras`)**           | Proveedores autorizados, RUC, órdenes de compra, facturas y cálculo automático de montos     | Soles (`PEN · S/`) |

---

## 2. Arquitectura aplicada en cada parte del sistema

Cada capa del sistema aplica una arquitectura específica elegida según su responsabilidad:

| Parte del sistema              | Arquitectura aplicada                                                                         | Detalle operativo                                                                                                                                                                                          |
| :----------------------------- | :-------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Frontend (interfaz)**        | App Router + Componentes Modulares + Estado Reactivo                                          | Páginas organizadas por dominio (`/proyectos`, `/inventario`, `/presupuestos`, `/personal`, `/compras`), estado global con **Zustand** y caché asíncrono con **TanStack Query**.                           |
| **Backend (lógica)**           | Microservicios con Arquitectura en Capas (`Handler → Service → Factory → Repository → Model`) | **5 servicios Go totalmente independientes**: un proceso, un puerto y una base de datos propia por módulo. Cada servicio se despliega y escala de forma aislada.                                           |
| **Persistencia**               | **Database-per-Service (Amazon RDS PostgreSQL)**                                              | Cada microservicio administra su propia base de datos PostgreSQL (`projects_db`, `inventory_db`, `budgets_db`, `staff_db`, `procurement_db`). Ningún servicio consulta directamente las tablas de otro.    |
| **Almacenamiento de Archivos** | **Object Storage (Amazon S3)**                                                                | Planos, fotos de obra, sustentos PDF y facturas se suben a un bucket de **Amazon S3** con _Presigned URLs_, guardando en PostgreSQL únicamente la referencia (`image_url`, `document_url`, `invoice_url`). |
| **Comunicación**               | HTTP/REST síncrona (JSON)                                                                     | Comunicación cliente-servidor y servicio-a-servicio (ej. `presupuestos-service` valida la existencia del proyecto consultando vía HTTP a `proyectos-service`).                                             |

---

## 3. Microservicios

Cada módulo de negocio vive en un proceso Go independiente dentro de `backend/services/` con su propio puerto, base de datos y fábrica de dominio:

| Servicio (`backend/services/`) | Puerto | Base de Datos (RDS) | Patrón Creacional                         | Responsabilidad                                                                                                                                                     |
| :----------------------------- | :----- | :------------------ | :---------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **`proyectos-service`**        | `8081` | `projects_db`       | Factory                                   | Crear y validar proyectos/obras, controlar su estado (`Planificación`, `En Ejecución`, `Finalizado`) y porcentaje de avance (`0–100%`).                             |
| **`inventario-service`**       | `8082` | `inventory_db`      | **Factory Method** + **Abstract Factory** | Crear y validar ítems diferenciando `MATERIAL` (consumible) y `HERRAMIENTA` (activo reutilizable por unidad), además de generar familias de kits por etapa de obra. |
| **`presupuestos-service`**     | `8083` | `budgets_db`        | Factory + Comunicación HTTP               | Validar partidas de materiales, mano de obra y equipos, calcular el `total_cost` fuera del controlador y verificar el proyecto vía HTTP contra `proyectos-service`. |
| **`personal-service`**         | `8084` | `staff_db`          | Factory                                   | Administrar trabajadores por especialidad (`Ingeniero Residente`, `Maestro de Obra`, `Operario`), jornal diario en Soles (`S/`) y asignación a obra.                |
| **`compras-service`**          | `8085` | `procurement_db`    | Factory                                   | Registrar proveedores, RUC y órdenes de compra calculando automáticamente el total (`quantity × unit_price`) en Soles (`S/`).                                       |

---

## 4. Patrones de diseño (Backend — Go)

Ambos patrones principales pertenecen al catálogo creacional (GoF) y resuelven el problema de no instanciar estructuras con condicionales y valores por defecto dispersos en los controladores HTTP.

### 4.1 Factory Method

- **Intención:** Definir una interfaz creadora (`ItemCreator`) con un método `Create` y delegar en creadores concretos (`MaterialCreator` y `ToolCreator`) la inicialización específica de cada tipo de recurso.
- **Dónde se aplica:** `backend/services/inventario-service/internal/factory/item_factory.go`.

```go
// ItemCreator define la interfaz del patrón Factory Method para crear recursos.
type ItemCreator interface {
	Create(input model.Item) (*model.Item, error)
}

// MaterialCreator crea objetos de tipo MATERIAL (consumibles de obra).
type MaterialCreator struct{}

// ToolCreator crea objetos de tipo HERRAMIENTA (activos reutilizables medidos por "Unidad").
type ToolCreator struct{}

// CreateInventoryItem es el Factory Method centralizado que selecciona el creador concreto.
func CreateInventoryItem(input model.Item) (*model.Item, error)
```

### 4.2 Abstract Factory

- **Intención:** Crear **familias de objetos relacionados o compatibles entre sí** (`Material Base` + `Herramienta Base`) según la etapa constructiva sin acoplar el servicio a las estructuras concretas.
- **Dónde se aplica:** `backend/services/inventario-service/internal/factory/abstract_factory.go`.

| Fábrica Concreta (`ConstructoraKitFactory`) | Etapa de Obra                  | Material de la Familia                      | Herramienta de la Familia                |
| :------------------------------------------ | :----------------------------- | :------------------------------------------ | :--------------------------------------- |
| **`ObraGruesaFactory`**                     | Cimentación y estructura       | Cemento Portland Tipo I (100 Bolsas)        | Mezcladora de Concreto 11p3 (2 Unidades) |
| **`AcabadosFactory`**                       | Revestimientos e instalaciones | Ladrillo Pastelero / Cerámico (10 Millares) | Taladro Percutor Industrial (4 Unidades) |

### 4.3 Relación entre ambos patrones

Dentro de `inventario-service`, el **Factory Method** (`CreateInventoryItem`) resuelve la creación individual de un producto según su categoría (`MATERIAL` vs `HERRAMIENTA`), mientras que el **Abstract Factory** (`GetKitFactory`) produce una **familia completa de productos compatibles** (`Material + Herramienta`) para abastecer una etapa de obra en una sola operación.

---

## 5. Frontend — Arquitectura por Capas y Componentes

Toda la aplicación cliente sigue la estructura estándar de **Next.js App Router + shadcn/ui** en `frontend/`:

```text
frontend/
├── app/                         # Enrutador Next.js (Páginas y módulos de negocio)
│   ├── api/upload/route.ts      # Endpoint firmador AWS Signature V4 (S3 Presigned URLs)
│   ├── proyectos/               # page.tsx, ProjectForm.tsx, ProjectList.tsx
│   ├── inventario/              # page.tsx, InventoryForm.tsx, InventoryList.tsx
│   ├── presupuestos/            # page.tsx, BudgetForm.tsx, BudgetList.tsx
│   ├── personal/                # page.tsx, StaffForm.tsx, StaffList.tsx
│   └── compras/                 # page.tsx, ProcurementForm.tsx, ProcurementList.tsx
├── components/                  # Librería de componentes reutilizables (shadcn/ui)
│   ├── ui/                      # Átomos reutilizables (button.tsx, card.tsx, file-upload.tsx)
│   └── layout/                  # Estructura global (Navbar.tsx, QueryProvider.tsx)
├── hooks/                       # TanStack Query (useProjects, useInventory, useBudgets, useStaff, useProcurement)
├── lib/                         # Utilidades shadcn/ui (utils.ts) y esquemas Zod (schemas.ts)
├── services/                    # Cliente HTTP REST hacia los 5 microservicios Go (api.ts)
├── store/                       # Estado global de UI con Zustand (useConstructoraStore.ts)
├── types/                       # Contratos TypeScript (index.ts)
└── components.json              # Configuración oficial de shadcn/ui
```

---

## 6. Stack Tecnológico por capa

| Capa                            | Tecnología                                            | Uso en el sistema                                                      |
| :------------------------------ | :---------------------------------------------------- | :--------------------------------------------------------------------- |
| **Frontend — Framework**        | React 19 + Next.js 16 (App Router) + TypeScript 5     | Renderizado modular, enrutamiento y tipado estático                    |
| **Frontend — Estilos e Iconos** | Tailwind CSS 4 + shadcn/ui + Lucide React             | Diseño limpio empresarial e iconografía de construcción                |
| **Frontend — Estado y Datos**   | Zustand + TanStack Query + React Hook Form + Zod      | Estado global de edición/filtros, caché HTTP e invalidación automática |
| **Backend — Lenguaje**          | Go 1.24 (`net/http`, `database/sql`, `encoding/json`) | 5 microservicios ligeros sin frameworks pesados                        |
| **Persistencia en Nube**        | **Amazon RDS for PostgreSQL** (`constructora-db`)     | 5 bases de datos relacionales aisladas (_Database-per-Service_)        |
| **Almacenamiento de Archivos**  | **Amazon S3** (`constructora-erp-104282840129`)       | Almacenamiento de planos, fichas técnicas, PDFs y facturas             |
| **Contenedores**                | Docker (Multi-stage builds) + Docker Compose          | Despliegue reproducible local y en AWS                                 |
| **Moneda**                      | Soles Peruanos — `PEN (S/)`                           | Formato `es-PE` (`S/ 1,234.56`) en todo el sistema                     |

---

## 7. Flujo de una petición CRUD (ejemplo: crear ítem de inventario)

1. El usuario completa el formulario en `/inventario` y presiona **Guardar en Inventario**.
2. **React Hook Form + Zod** validan los tipos y restricciones en el cliente.
3. **TanStack Query** (`useInventory`) ejecuta `POST http://localhost:8082/inventory` con el cuerpo JSON.
4. `InventoryHandler.Create` recibe la petición y delega en `InventoryService.Create`.
5. `InventoryService` invoca al **Factory Method** `factory.CreateInventoryItem(input)`, el cual selecciona `MaterialCreator` o `ToolCreator` y asigna la unidad y condición correspondientes.
6. `InventoryRepository.Create` ejecuta el `INSERT INTO inventory_items ... RETURNING id` en la base de datos **`inventory_db`** alojada en **Amazon RDS PostgreSQL**.
7. El microservicio responde `201 Created` con el ítem creado y **TanStack Query** invalida la caché `["inventory"]`, actualizando la lista en pantalla al instante.

---

## 8. Convención de moneda y Despliegue en AWS (RDS + S3)

- **Moneda:** Todos los montos financieros del sistema (presupuestos de materiales, mano de obra, equipos, jornales diarios del personal y órdenes de compra a proveedores) se expresan y formatean en **Soles Peruanos (`PEN · S/`)** utilizando `toLocaleString("es-PE", { minimumFractionDigits: 2 })` (ejemplo: `S/ 165,500.00`).
- **Infraestructura AWS configurada (`us-east-1`):**
  - **Amazon RDS PostgreSQL:** Instancia `constructora-db` con auto-inicialización de las 5 bases de datos (`projects_db`, `inventory_db`, `budgets_db`, `staff_db`, `procurement_db`).
  - **Amazon S3:** Bucket `constructora-erp-104282840129` con carga directa mediante firma criptográfica **AWS Signature V4** y lectura segura con _Presigned URLs_.
