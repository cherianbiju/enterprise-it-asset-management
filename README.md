# Enterprise IT Asset Management System

A production-grade, full-stack IT Asset Management platform engineered with a high-performance **Go (Golang)** REST API backend, modern **Angular** standalone component frontend, and robust **PostgreSQL** relational database.

---

## ⚡ Quick Start for Evaluator / Tutor (Run in 3 Minutes)

Follow these 3 simple steps to launch the entire system from scratch:

```bash
# 1. Clone the repository and enter the directory
git clone https://github.com/YOUR_USERNAME/enterprise-it-asset-management.git
cd enterprise-it-asset-management

# 2. Setup environment config (if backend/.env does not already exist)
cp backend/.env.example backend/.env

# 3. Start PostgreSQL Database (Docker will automatically load schema & seed data)
docker compose up -d

# 4. Start Go Backend (in Terminal 1)
cd backend
go run cmd/main.go
# -> Backend is running at http://localhost:8080 (Health check: http://localhost:8080/health)

# 5. Start Angular Frontend (in Terminal 2)
cd frontend
npm install
npm start
# -> Open http://localhost:4200 in your browser
```

---

## Architecture & Technology Stack

| Layer | Technology | Key Capabilities |
|---|---|---|
| **Database** | PostgreSQL | Foreign keys, unique constraints, transactional ACID integrity, audit logging |
| **Backend** | Go (1.22+) | Native `net/http` routing with Go 1.22 pattern matching, connection pooling, CORS middleware |
| **Frontend** | Angular (18+/19) | Modern standalone components, reactive forms, signals, HttpClient, responsive design |
| **API Format** | RESTful JSON | Standardized status codes, pagination metadata, and error handling |

---

## Directory Structure

```text
it-asset-management/
├── database/
│   └── schema.sql                             # SQL schema with tables, constraints, indexes & seed data
├── backend/
│   ├── cmd/
│   │   └── main.go                            # Server entry point & graceful shutdown
│   ├── internal/
│   │   ├── config/
│   │   │   └── db.go                          # PostgreSQL connection pool & .env loader
│   │   ├── models/
│   │   │   └── asset.go                       # Data structs, DTOs & status types
│   │   ├── repository/
│   │   │   └── asset_repository.go            # SQL queries, ACID transactions & business rules
│   │   ├── handlers/
│   │   │   └── asset_handler.go               # HTTP controllers for all REST endpoints
│   │   └── middleware/
│   │       └── cors.go                        # Cross-Origin Resource Sharing (CORS) handler
│   ├── go.mod                                 # Go module dependencies
│   ├── go.sum                                 # Go dependency checksums
│   ├── .env.example                           # Example environment variables
│   └── .env                                   # Local runtime environment variables
└── frontend/
    ├── src/
    │   ├── app/
    │   │   ├── models/
    │   │   │   └── asset.model.ts             # TypeScript interfaces matching backend models
    │   │   ├── services/
    │   │   │   └── asset.service.ts           # Angular HttpClient service connecting to Go API
    │   │   ├── components/
    │   │   │   ├── dashboard/                 # Metrics overview & shortcuts
    │   │   │   │   ├── dashboard.component.ts
    │   │   │   │   ├── dashboard.component.html
    │   │   │   │   └── dashboard.component.css
    │   │   │   ├── asset-list/                # Inventory table, filters, pagination & action modals
    │   │   │   │   ├── asset-list.component.ts
    │   │   │   │   ├── asset-list.component.html
    │   │   │   │   └── asset-list.component.css
    │   │   │   ├── asset-form/                # Register new asset & edit asset attributes
    │   │   │   │   ├── asset-form.component.ts
    │   │   │   │   ├── asset-form.component.html
    │   │   │   │   └── asset-form.component.css
    │   │   │   └── asset-history/             # Visual audit timeline of device events
    │   │   │       ├── asset-history.component.ts
    │   │   │       ├── asset-history.component.html
    │   │   │       └── asset-history.component.css
    │   │   ├── app.component.ts               # Root layout & enterprise header navigation
    │   │   ├── app.routes.ts                  # Client-side routing configuration
    │   │   └── app.config.ts                  # Standalone application bootstrap configuration
    │   ├── index.html                         # Base HTML template
    │   ├── main.ts                            # Angular bootstrap entry point
    │   └── styles.css                         # Global CSS styles, badges, and variables
    ├── angular.json                           # Angular workspace configuration
    ├── package.json                           # Dependencies and scripts
    └── README.md                              # Frontend documentation
```

---

## Business Rules Enforced

1. **Asset Fields & Constraints**:
   - `id`: Auto-incrementing primary key.
   - `asset_code`: Unique identifier (e.g. `AST-1001`). Immutable after creation.
   - `asset_name`: Description / model name.
   - `asset_type`: Category (`LAPTOP`, `DESKTOP`, `MONITOR`, `SERVER`, `PHONE`, `NETWORK`, `PERIPHERAL`, `OTHER`).
   - `serial_number`: Unique hardware manufacturer identifier. Immutable after creation.
   - `purchase_date`: Valid date format (`YYYY-MM-DD`).
   - `status`: One of `AVAILABLE`, `ASSIGNED`, `UNDER_REPAIR`, `RETIRED`.
   - `assigned_to`: Name of employee or department (populated only when assigned).
   - `location`: Physical site, building, floor, or storage room.

2. **Assignment Rules**:
   - **Validation Check**: An asset **CANNOT** be assigned if its status is already `ASSIGNED`, `UNDER_REPAIR`, or `RETIRED`. Only assets in `AVAILABLE` status can be deployed.
   - **Return Workflow**: Returning an asset transitions its status back to `AVAILABLE`, clears the `assigned_to` field, updates its storage location, and automatically writes an audit record to `asset_history` (`action = 'RETURNED'`).
   - **Repair Workflow**: Marking an asset for repair sets its status to `UNDER_REPAIR` and logs the issue diagnostics to `asset_history` (`action = 'REPAIR'`). Retired assets cannot be sent for repair.
   - **Retirement Workflow**: Decommissioning marks the asset as `RETIRED`, unassigns any current custodian, and logs the decommission reason into `asset_history` (`action = 'RETIRED'`).

---

## Step-by-Step Setup Instructions

### 1. Database Setup (PostgreSQL)

1. Ensure PostgreSQL is installed and running on your system.
2. Connect to PostgreSQL and create the database:
   ```sql
   CREATE DATABASE it_asset_management;
   ```
3. Run the migration and seed script:
   ```bash
   psql -U postgres -d it_asset_management -f database/schema.sql
   ```

---

### 2. Backend Setup (Go)

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Verify or update your environment settings in `.env`:
   ```env
   SERVER_PORT=8080
   CORS_ORIGIN=http://localhost:4200
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=it_asset_management
   DB_SSLMODE=disable
   ```
3. Download dependencies:
   ```bash
   go mod download
   ```
4. Start the backend service:
   ```bash
   go run cmd/main.go
   ```
   The backend will start and listen on `http://localhost:8080`.
   You can verify it is healthy by visiting: `http://localhost:8080/health`.

---

### 3. Frontend Setup (Angular)

1. Open a new terminal and navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install npm dependencies:
   ```bash
   npm install
   ```
3. Launch the Angular development server:
   ```bash
   npm start
   ```
4. Open your browser and navigate to:
   ```text
   http://localhost:4200
   ```

---

## REST API Reference

### Health & Analytics
- `GET /health` — Check server status.
- `GET /api/assets/stats` — Dashboard metrics summary (total, available, assigned, under_repair, retired).

### Assets CRUD
- `GET /api/assets` — Retrieve paginated assets list.
  - Query parameters:
    - `search` (searches code, name, and serial number)
    - `status` (`AVAILABLE`, `ASSIGNED`, `UNDER_REPAIR`, `RETIRED`)
    - `type` (`LAPTOP`, `MONITOR`, etc.)
    - `page` (default: 1)
    - `limit` (default: 10)
- `GET /api/assets/:id` — Retrieve single asset details.
- `POST /api/assets` — Create a new asset.
- `PUT /api/assets/:id` — Update asset attributes.
- `DELETE /api/assets/:id` — Delete asset and cascade remove history.

### Lifecycle Actions & Audit Trail
- `POST /api/assets/:id/assign` — Assign asset to employee.
  - Payload: `{ "assigned_to": "Jane Doe", "location": "Floor 3", "notes": "Onboarding" }`
- `POST /api/assets/:id/return` — Return asset to available inventory.
  - Payload: `{ "location": "IT Storage Room A", "notes": "Project wrapped up" }`
- `POST /api/assets/:id/repair` — Send asset for repair.
  - Payload: `{ "location": "Authorized Vendor", "notes": "Battery replacement" }`
- `POST /api/assets/:id/retire` — Decommission asset.
  - Payload: `{ "location": "E-Waste Bin", "notes": "End of 4-year cycle" }`
- `GET /api/assets/:id/history` — Fetch complete audit trail for an asset.

---

## CORS Verification

The Go backend includes a specialized CORS middleware in [`backend/internal/middleware/cors.go`](file:///c:/Users/CHERIAN%20BIJU/Desktop/it-asset-management/backend/internal/middleware/cors.go) configured to allow requests originating from `http://localhost:4200` with preflight `OPTIONS` handling. You can test this using `curl`:

```bash
curl -i -X OPTIONS http://localhost:8080/api/assets \
  -H "Origin: http://localhost:4200" \
  -H "Access-Control-Request-Method: GET"
```
You will receive an HTTP `204 No Content` response with the appropriate `Access-Control-Allow-Origin: http://localhost:4200` headers.
