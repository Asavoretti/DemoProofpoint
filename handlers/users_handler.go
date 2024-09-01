package handlers

import (
	db "DemoProofpoint/dataBase"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

// GetUsers retrieves and returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get users received")
	db, err := db.ConectarBD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre, email FROM usuarios")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []Usuario
	for rows.Next() {
		var user Usuario
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	log.Println("Users successfully sent in the response")
}

// CreateUser adds a new user to the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to create user received")
	db, err := db.ConectarBD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user Usuario
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO usuarios(nombre, email) VALUES(?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Nombre, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	log.Println("User successfully created:", user)
}

// UpdateUser updates the information of a user in the database
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to update user received")
	db, err := db.ConectarBD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	userID := vars["id"]

	var user Usuario
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE usuarios SET nombre=?, email=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Nombre, user.Email, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	log.Println("User successfully updated:", user)
}

// DeleteUser deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to delete user received")
	db, err := db.ConectarBD()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	userID := vars["id"]

	stmt, err := db.Prepare("DELETE FROM usuarios WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("User successfully deleted, ID:", userID)
}

// ServePage serves the HTML page
func ServePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving the HTML page")
	http.ServeFile(w, r, "index.html")
}
