package response

import (
	"encoding/json"
	"net/http"
)

// структура успешного ответа
type successResponse struct {
	Status string `json:"status"` 
	Data   any    `json:"data"`   
}

// структура ошибки
type errorResponse struct {
	Status string `json:"status"` 
	Error  string `json:"error"`  
}

// JSON — универсальный метод для успешного ответа.
// data — это полезная нагрузка (например, пользователь или список товаров).
// statusCode — HTTP код (например, 200 OK, 201 Created).
func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := successResponse{
		Status: "success",
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"status":"error","error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

// Error — универсальный метод для ошибок.
// Возвращает JSON с ключами {"status": "error", "error": "<текст ошибки>"}.
// errorMessage — это обычная ошибка (например errors.New("user not found")).
func ErrorResp(w http.ResponseWriter, statusCode int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := errorResponse{
		Status: "error",
		Error:  err,
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		http.Error(w, `{"status":"error","error":"failed to encode error response"}`, http.StatusInternalServerError)
	}
}

// Shortcut-хелперы для часто используемых статусов.
// Они помогают сократить код в хендлерах.

func BadRequest(w http.ResponseWriter, msg string) {
	ErrorResp(w, http.StatusBadRequest, msg)
}

func NotFound(w http.ResponseWriter, msg string) {
	ErrorResp(w, http.StatusNotFound, msg)
}

func InternalError(w http.ResponseWriter, msg string) {
	ErrorResp(w, http.StatusInternalServerError, msg)
}
