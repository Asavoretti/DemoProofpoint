package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq" // Importa el driver de PostgreSQL
    "github.com/gorilla/mux"
)

// Usuario representa la estructura de un usuario
type Usuario struct {
    ID    int    `json:"id"`
    Nombre  string `json:"nombre"`
    Email string `json:"email"`
}

// Conexión a la base de datos
func conectarBD() (*sql.DB, error) {
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbHost := os.Getenv("DB_HOST")
    dbName := os.Getenv("DB_NAME")

    dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }
    return db, nil
}

// Obtener usuarios
func obtenerUsuarios(w http.ResponseWriter, r *http.Request) {
    db, err := conectarBD()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    filas, err := db.Query("SELECT id, nombre, email FROM usuarios")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer filas.Close()

    var usuarios []Usuario
    for filas.Next() {
        var usuario Usuario
        if err := filas.Scan(&usuario.ID, &usuario.Nombre, &usuario.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        usuarios = append(usuarios, usuario)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(usuarios)
}

// Crear un usuario
func crearUsuario(w http.ResponseWriter, r *http.Request) {
    db, err := conectarBD()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var usuario Usuario
    if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    stmt, err := db.Prepare("INSERT INTO usuarios(nombre, email) VALUES($1, $2)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(usuario.Nombre, usuario.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(usuario)
}

// Servir la página HTML
func servirPagina(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", servirPagina).Methods("GET")
    r.HandleFunc("/usuarios", obtenerUsuarios).Methods("GET")
    r.HandleFunc("/usuarios", crearUsuario).Methods("POST")

    http.Handle("/", r)
    log.Println("Servidor iniciado en el puerto 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
