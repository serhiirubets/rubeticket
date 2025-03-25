package concerts

import (
	"net/http"
	"strconv"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/req"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type ConcertHandlerDeps struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        IConcertService
	UserRepository users.IUserRepository
}

type ConcertHandler struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        IConcertService
	UserRepository users.IUserRepository
}

func NewConcertHandler(router *http.ServeMux, deps *ConcertHandlerDeps) {
	handler := ConcertHandler{
		Config:         deps.Config,
		Logger:         deps.Logger,
		Service:        deps.Service,
		UserRepository: deps.UserRepository,
	}

	router.HandleFunc("POST /admin/concerts", handler.Create())
	router.HandleFunc("PUT /admin/concerts/{id}", handler.Update())
	router.HandleFunc("DELETE /admin/concerts/{id}", handler.Delete())
	router.HandleFunc("GET /admin/concerts/{id}", handler.GetByID())
	router.HandleFunc("GET /admin/concerts", handler.List())
}

// Create godoc
// @Summary Create a new concert
// @Description Create a new concert with the provided details
// @Tags Admin/Concerts
// @Accept json
// @Produce json
// @Param request body CreateConcertRequest true "Concert details"
// @Success 201 {object} ConcertResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/concerts [post]
func (h *ConcertHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[CreateConcertRequest](&w, r)
		if err != nil {
			return
		}

		concert, err := h.Service.Create(payload)
		if err != nil {
			h.Logger.Error("Failed to create concert", "error", err.Error())
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, concert, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update a concert
// @Description Update an existing concert
// @Tags Admin/Concerts
// @Accept json
// @Produce json
// @Param id path int true "Concert ID"
// @Param request body UpdateConcertRequest true "Concert details"
// @Success 200 {object} ConcertResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/concerts/{id} [put]
func (h *ConcertHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid concert ID", http.StatusBadRequest)
			return
		}

		payload, err := req.HandleBody[UpdateConcertRequest](&w, r)
		if err != nil {
			return
		}

		concert, err := h.Service.Update(uint(id), payload)
		if err != nil {
			if err.Error() == "concert not found" {
				res.Json(w, "Concert not found", http.StatusNotFound)
				return
			}
			h.Logger.Error("Failed to update concert", "error", err.Error())
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, concert, http.StatusOK)
	}
}

// Delete godoc
// @Summary Delete a concert
// @Description Delete an existing concert
// @Tags Admin/Concerts
// @Produce json
// @Param id path int true "Concert ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/concerts/{id} [delete]
func (h *ConcertHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid concert ID", http.StatusBadRequest)
			return
		}

		err = h.Service.Delete(uint(id))
		if err != nil {
			h.Logger.Error("Failed to delete concert", "error", err.Error())
			res.Json(w, "Failed to delete concert", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetByID godoc
// @Summary Get a concert by ID
// @Description Get details of a specific concert
// @Tags Admin/Concerts
// @Produce json
// @Param id path int true "Concert ID"
// @Success 200 {object} ConcertResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/concerts/{id} [get]
func (h *ConcertHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid concert ID", http.StatusBadRequest)
			return
		}

		concert, err := h.Service.GetByID(uint(id))
		if err != nil {
			res.Json(w, "Concert not found", http.StatusNotFound)
			return
		}

		res.Json(w, concert, http.StatusOK)
	}
}

// List godoc
// @Summary List concerts
// @Description Get a paginated list of concerts
// @Tags Admin/Concerts
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ListConcertsResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/concerts [get]
func (h *ConcertHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		concerts, err := h.Service.List(page, pageSize)
		if err != nil {
			h.Logger.Error("Failed to list concerts", "error", err.Error())
			res.Json(w, "Failed to list concerts", http.StatusInternalServerError)
			return
		}

		res.Json(w, concerts, http.StatusOK)
	}
}
