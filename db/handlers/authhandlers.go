package handlers

import (
    "database/sql"
    "dbservice/models"
    "encoding/json"
    "net/http"

	"github.com/gorilla/mux"
)

//Новый пользователь
func CreateUser(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		var req models.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return 
		}

		//Валидируем
		if req.Username == "" || req.PasswordHash == "" {
			http.Error(w, "Username and password_hash are required", http.StatusBadRequest)
			return 
		}

		//Проверка существования
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, req.Username).Scan(&exists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return 
		}

		if exists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return 
		}

		//Создаём пользователя
		var user models.User
		 err = db.QueryRow(
            `INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id, username, created_at`,
            req.Username,
            req.PasswordHash,
        ).Scan(&user.ID, &user.Username, &user.CreatedAt)

		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

//Получаем пользователя по имени
func GetUserByUsername(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		username := vars["username"]
		
		//И вновь валидируем
		if username == ""{
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		
		//Получаем пользователя из БД
		var user models.User
		err := db.QueryRow(
			`SELECT id, username, password_hash, created_at 
			FROM users 
			WHERE username = $1`, username,
		).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		userID := vars["id"]

		//И тут тоже валидируем
		if userID == ""{
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		//Получаем пользователя из БД
		var user models.User
		err := db.QueryRow(
			`SELECT id, username, password_hash, created_at 
			FROM users 
			WHERE id = $1`, userID,
		).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}