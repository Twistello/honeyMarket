package user

import (

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"grishoney/internal/pkg/response"
	"grishoney/internal/pkg/middleware"
)

type UserHandler struct {
	service *UserService
}

// Конструктор
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes — регистрирует все эндпоинты пользователя в роутере
func (h *UserHandler) RegisterRoutes() chi.Router {

	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register",h.Register)
		r.Post("/", h.Login)
	})
	
	r.Route("/", func(r chi.Router) {
		//r.Use(middleware.JWTAuth)
		//r.Use(middleware.Role)
		r.With(middleware.Paginate).Get("/{id}",h.List)
		r.Delete("/",h.Delete)
		r.Put("/{id}/role", h.UpdateRole)
	})

	return r
	
}

// ===== Handlers =====

// Register — регистрация нового пользователя
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResp(w, 400, "invalid request body")
		return
	}

	u, err := h.service.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		response.ErrorResp(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, u)
}

// Login — аутентификация пользователя
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResp(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.service.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		response.ErrorResp(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, u)
}

// GetByID — получение пользователя по ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResp(w, http.StatusBadRequest, "invalid id format")
		return
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.ErrorResp(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, u)
}

// UpdateRole — обновление роли пользователя

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResp(w, http.StatusBadRequest, "invalid id format")
		return
	}

	var req struct {
		NewRole string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResp(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.service.UpdateRole(r.Context(), id, req.NewRole)
	if err != nil {
		response.ErrorResp(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, u)
}


// Delete — удаление пользователя
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResp(w, http.StatusBadRequest, "invalid id format")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		response.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}

// List — список пользователей с пагинацией
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {

	limit, offset := middleware.GetPagination(r)

	users, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		response.ErrorResp(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, users)
}

/*

func getRole(r *http.Request) (role string) {
	ctx := r.Context()

	role, _ = ctx.Value(RoleType{}).(string)

	return
}
	*/