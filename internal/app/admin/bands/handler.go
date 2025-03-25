package bands

import (
	"net/http"
	"strconv"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/req"
	"github.com/serhiirubets/rubeticket/internal/pkg/res"
)

type BandHandlerDeps struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        *BandService
	UserRepository users.IUserRepository
}

type BandHandler struct {
	Config         *config.Config
	Logger         log.ILogger
	Service        *BandService
	UserRepository users.IUserRepository
}

func NewBandHandler(router *http.ServeMux, deps *BandHandlerDeps) {
	handler := BandHandler{
		Config:         deps.Config,
		Logger:         deps.Logger,
		Service:        deps.Service,
		UserRepository: deps.UserRepository,
	}

	router.HandleFunc("POST /bands", handler.Create())
	router.HandleFunc("PUT /bands/{id}", handler.Update())
	router.HandleFunc("DELETE /bands/{id}", handler.Delete())
	router.HandleFunc("GET /bands/{id}", handler.GetByID())
	router.HandleFunc("GET /bands", handler.List())
}

// Create godoc
// @Summary Create a new band
// @Description Create a new band with the provided details
// @Tags Admin/Bands
// @Accept json
// @Produce json
// @Param request body CreateBandRequest true "Band details"
// @Success 201 {object} BandResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/bands [post]
func (h *BandHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[CreateBandRequest](&w, r)
		if err != nil {
			return
		}

		band, err := h.Service.Create(payload)
		if err != nil {
			h.Logger.Error("Failed to create band", "error", err.Error())
			res.Json(w, "Failed to create band", http.StatusInternalServerError)
			return
		}

		response := BandResponse{
			ID:          band.ID,
			Name:        band.Name,
			Description: band.Description,
		}

		res.Json(w, response, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update a band
// @Description Update an existing band with the provided details
// @Tags Admin/Bands
// @Accept json
// @Produce json
// @Param id path int true "Band ID"
// @Param request body UpdateBandRequest true "Band details to update"
// @Success 200 {object} BandResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/bands/{id} [put]
func (h *BandHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid band ID", http.StatusBadRequest)
			return
		}

		payload, err := req.HandleBody[UpdateBandRequest](&w, r)
		if err != nil {
			return
		}

		band, err := h.Service.Update(uint(id), payload)
		if err != nil {
			if err.Error() == "band not found" {
				res.Json(w, "Band not found", http.StatusNotFound)
				return
			}
			h.Logger.Error("Failed to update band", "error", err.Error())
			res.Json(w, "Failed to update band", http.StatusInternalServerError)
			return
		}

		response := BandResponse{
			ID:          band.ID,
			Name:        band.Name,
			Description: band.Description,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// Delete godoc
// @Summary Delete a band
// @Description Delete an existing band
// @Tags Admin/Bands
// @Produce json
// @Param id path int true "Band ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/bands/{id} [delete]
func (h *BandHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid band ID", http.StatusBadRequest)
			return
		}

		err = h.Service.Delete(uint(id))
		if err != nil {
			h.Logger.Error("Failed to delete band", "error", err.Error())
			res.Json(w, "Failed to delete band", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetByID godoc
// @Summary Get a band by ID
// @Description Get details of a specific band
// @Tags Admin/Bands
// @Produce json
// @Param id path int true "Band ID"
// @Success 200 {object} BandResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /admin/bands/{id} [get]
func (h *BandHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "Invalid band ID", http.StatusBadRequest)
			return
		}

		band, err := h.Service.GetByID(uint(id))
		if err != nil {
			res.Json(w, "Band not found", http.StatusNotFound)
			return
		}

		response := BandResponse{
			ID:          band.ID,
			Name:        band.Name,
			Description: band.Description,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// List godoc
// @Summary List bands
// @Description Get a paginated list of bands
// @Tags Admin/Bands
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ListBandsResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /admin/bands [get]
func (h *BandHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		bands, err := h.Service.List(page, pageSize)
		if err != nil {
			h.Logger.Error("Failed to list bands", "error", err.Error())
			res.Json(w, "Failed to list bands", http.StatusInternalServerError)
			return
		}

		response := ListBandsResponse{
			Items: make([]BandResponse, len(bands)),
		}

		for i, band := range bands {
			response.Items[i] = BandResponse{
				ID:          band.ID,
				Name:        band.Name,
				Description: band.Description,
			}
		}

		res.Json(w, response, http.StatusOK)
	}
}
