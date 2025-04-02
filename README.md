# Stock Management API

This is a RESTful API built using Go and PostgreSQL to manage stock information. It provides endpoints to create, retrieve, update, and delete stock records.

## Features
- CRUD operations for stock data
- PostgreSQL database integration
- Gorilla Mux for routing
- Environment variables for database credentials

## Requirements
- Go 1.18+
- PostgreSQL
- Git

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

2. Create a `.env` file in the root directory and add your PostgreSQL credentials:
   ```env
   POSTGRESQL_URL=your_database_url
   ```

3. Install dependencies:
   ```sh
   go mod tidy
   ```

## Running the API

1. Build and run the server:
   ```sh
   go run main.go
   ```

2. The server will start at `http://localhost:8080`.

## API Endpoints

### Get all stocks
```http
GET /api/stock/
```

### Get stock by ID
```http
GET /api/stock/{id}
```

### Create a new stock
```http
POST /api/newstock/
```

### Update stock
```http
PUT /api/stock/{id}
```

### Delete stock
```http
DELETE /api/deletestock/{id}
```

## Project Structure
```
📂 psq-project
├── 📂 middleware
│   ├── handlers.go
├── 📂 models
│   ├── models.go
├── 📂 router
│   ├── router.go
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## License
This project is open-source and available under the MIT License.

---
Made by [Pradyumna A J](https://github.com/pradyumnajavalagi)

