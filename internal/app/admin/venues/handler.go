package venues

import (
	"net/http"
	"strconv"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/req"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type VenueHandlerDeps struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        *VenueService
	UserRepository users.IUserRepository
}

type VenueHandler struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        *VenueService
	UserRepository users.IUserRepository
}

func NewVenueHandler(router *http.ServeMux, deps *VenueHandlerDeps) {
	handler := VenueHandler{
		Config:         deps.Config,
		Logger:         deps.Logger,
		Service:        deps.Service,
		UserRepository: deps.UserRepository,
	}

	router.HandleFunc("POST /admin/venues", handler.Create())
	router.HandleFunc("PUT /admin/venues/{id}", handler.Update())
	router.HandleFunc("DELETE /admin/venues/{id}", handler.Delete())
	router.HandleFunc("GET /admin/venues/{id}", handler.GetByID())
	router.HandleFunc("GET /admin/venues", handler.List())
}

// Create godoc
// @Summary Create a new venue
// @Description Create a new venue with the provided details
// @Tags Admin/Venues
// @Accept json
// @Produce json
// @Param request body CreateVenueRequest true "Venue details"
// @Success 201 {object} VenueResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/v1/venues [post]
func (h *VenueHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[CreateVenueRequest](&w, r)
		if err != nil {
			return
		}

		venue, err := h.Service.Create(payload)
		if err != nil {
			h.Logger.Error("Failed to create venue", "error", err.Error())
			res.Json(w, "Failed to create venue", http.StatusInternalServerError)
			return
		}

		response := VenueResponse{
			ID:          venue.ID,
			Name:        venue.Name,
			Description: venue.Description,
			Address:     venue.Address,
			Phone:       venue.Phone,
			Email:       venue.Email,
		}

		res.Json(w, response, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update a venue
// @Description Update an existing venue with the provided details
// @Tags Admin/Venues
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param request body UpdateVenueRequest true "Venue details to update"
// @Success 200 {object} VenueResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/v1/venues/{id} [put]
func (h *VenueHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid venue ID", http.StatusBadRequest)
			return
		}

		payload, err := req.HandleBody[UpdateVenueRequest](&w, r)
		if err != nil {
			return
		}

		venue, err := h.Service.Update(uint(id), payload)
		if err != nil {
			if err.Error() == "venue not found" {
				res.Json(w, "Venue not found", http.StatusNotFound)
				return
			}
			h.Logger.Error("Failed to update venue", "error", err.Error())
			res.Json(w, "Failed to update venue", http.StatusInternalServerError)
			return
		}

		response := VenueResponse{
			ID:          venue.ID,
			Name:        venue.Name,
			Description: venue.Description,
			Address:     venue.Address,
			Phone:       venue.Phone,
			Email:       venue.Email,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// Delete godoc
// @Summary Delete a venue
// @Description Delete an existing venue
// @Tags Admin/Venues
// @Produce json
// @Param id path int true "Venue ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/v1/venues/{id} [delete]
func (h *VenueHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid venue ID", http.StatusBadRequest)
			return
		}

		err = h.Service.Delete(uint(id))
		if err != nil {
			h.Logger.Error("Failed to delete venue", "error", err.Error())
			res.Json(w, "Failed to delete venue", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetByID godoc
// @Summary Get a venue by ID
// @Description Get details of a specific venue
// @Tags Admin/Venues
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} VenueResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/v1/venues/{id} [get]
func (h *VenueHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid venue ID", http.StatusBadRequest)
			return
		}

		venue, err := h.Service.GetByID(uint(id))
		if err != nil {
			res.Json(w, "Venue not found", http.StatusNotFound)
			return
		}

		response := VenueResponse{
			ID:          venue.ID,
			Name:        venue.Name,
			Description: venue.Description,
			Address:     venue.Address,
			Phone:       venue.Phone,
			Email:       venue.Email,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// List godoc
// @Summary List venues
// @Description Get a paginated list of venues
// @Tags Admin/Venues
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ListVenuesResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/v1/venues [get]
func (h *VenueHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		venues, err := h.Service.List(page, pageSize)
		if err != nil {
			h.Logger.Error("Failed to list venues", "error", err.Error())
			res.Json(w, "Failed to list venues", http.StatusInternalServerError)
			return
		}

		response := ListVenuesResponse{
			Items: make([]VenueResponse, len(venues.Items)),
		}

		for i, venue := range venues.Items {
			response.Items[i] = VenueResponse{
				ID:          venue.ID,
				Name:        venue.Name,
				Description: venue.Description,
				Address:     venue.Address,
				Phone:       venue.Phone,
				Email:       venue.Email,
			}
		}

		res.Json(w, response, http.StatusOK)
	}
}
