package handlers

import (
    "apiservice/auth"
    "apiservice/models"
    "bytes"
    "encoding/json"
    "net/http"
)

//Регистрируемся
func Register(w http.ResponseWriter, r *http.Request){
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error: Invalid JSON", http.StatusBadRequest)
		return
	}

	//Валидируем
	if req.Username == "" || req.Password == "" {
		http.Error(w, "error: Username and Password are required", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 {
		http.Error(w, "error: Username must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "error: Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	//Хэшируем
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	//Отправляем в db
	createUserReq := map[string]string{
		"username": req.Username,
		"password_hash": hashedPassword,
	}

	jsonData, err := json.Marshal(createUserReq)
	resp, err := http.Post(
		"http://dbservice:8081/users",
		"application/json",
		 bytes.NewBuffer(jsonData),
	)

	if err != nil {
		http.Error(w, "error: Failed to create user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		http.Error(w, "Username already exists", http.StatusConflict)
	}

	//Получаем созданного юзера
	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user data", http.StatusInternalServerError)
		return
	}

	//Наш JWT
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	//Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	})
}

//Логинимся
func Login(w http.ResponseWriter, r *http.Request){
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `Invalid request body`, http.StatusBadRequest)
		return
	}

	//Валидируем
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	//Получаем юзера из db
	resp, err := http.Get("http://dbservice:8080/users/%s" + req.Username)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, `Failed to get user`, http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	//Сверяем пароль
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//Наш JWT
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	//Ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	})
}