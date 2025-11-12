package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/jackc/pgx/v5"
)

// Employee struct matches your table
type Employee struct {
    ID        int     `json:"id"`
    FirstName string  `json:"firstname"`
    LastName  string  `json:"lastname"`
    Age       int     `json:"age"`
    Email     string  `json:"email"`
}

// Global connection
var conn *pgx.Conn

func main() {
    ctx := context.Background()

    // 1️⃣ Connect to the database
    var err error
    conn, err = pgx.Connect(ctx, "psql 'postgresql://neondb_owner:npg_4LIbQX6PayUK@ep-dawn-pine-adu0ytt4-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require'")
    if err != nil {
        log.Fatal("Unable to connect to database:", err)
    }
    defer conn.Close(ctx)
    fmt.Println("✅ Connected to PostgreSQL!")

    // 3️⃣ Register HTTP handlers
    http.HandleFunc("/employees", employeesHandler)

    fmt.Println("🚀 Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// employeesHandler handles both GET and POST
func employeesHandler(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case http.MethodGet:
        getEmployees(ctx, w)
    case http.MethodPost:
        addEmployee(ctx, w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// GET /employees → list all
func getEmployees(ctx context.Context, w http.ResponseWriter) {
    rows, err := conn.Query(ctx, `SELECT id, firstname, lastname, age, email FROM employee_demo ORDER BY id;`)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var e Employee
        err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Age, &e.Email)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employees = append(employees, e)
    }

    json.NewEncoder(w).Encode(employees)
}

// POST /employees → add new
func addEmployee(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    var e Employee
    err := json.NewDecoder(r.Body).Decode(&e)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    query := `
    INSERT INTO employee_demo (firstname, lastname, age, email)
    VALUES ($1, $2, $3, $4)
    RETURNING id;
    `
    err = conn.QueryRow(ctx, query, e.FirstName, e.LastName, e.Age, e.Email).Scan(&e.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(e)
}