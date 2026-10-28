package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"it-asset-management/backend/internal/models"
	"it-asset-management/backend/internal/repository"
)

type AssetHandler struct {
	repo repository.AssetRepository
}

func NewAssetHandler(repo repository.AssetRepository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

// GET /api/assets
func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	status := query.Get("status")
	assetType := query.Get("type")

	page := 1
	if p := query.Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10
	if l := query.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	assets, total, err := h.repo.GetAssets(r.Context(), search, status, assetType, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	response := models.PaginatedResponse{
		Data:       assets,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	writeJSON(w, http.StatusOK, response)
}

// GET /api/assets/stats
func (h *AssetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// GET /api/assets/:id
func (h *AssetHandler) GetAssetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	asset, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// POST /api/assets
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.AssetCode) == "" || strings.TrimSpace(req.AssetName) == "" ||
		strings.TrimSpace(req.AssetType) == "" || strings.TrimSpace(req.SerialNumber) == "" ||
		strings.TrimSpace(req.PurchaseDate) == "" || strings.TrimSpace(req.Location) == "" {
		writeError(w, http.StatusBadRequest, "Missing required fields: asset_code, asset_name, asset_type, serial_number, purchase_date, location")
		return
	}

	if req.Status != "" && !req.Status.IsValid() {
		writeError(w, http.StatusBadRequest, "Invalid asset status. Allowed values: AVAILABLE, ASSIGNED, UNDER_REPAIR, RETIRED")
		return
	}

	asset, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateCode) {
			writeError(w, http.StatusConflict, "Asset code or serial number already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, asset)
}

// PUT /api/assets/:id
func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req models.UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.AssetName) == "" || strings.TrimSpace(req.AssetType) == "" ||
		strings.TrimSpace(req.PurchaseDate) == "" || strings.TrimSpace(req.Location) == "" {
		writeError(w, http.StatusBadRequest, "Missing required fields for update")
		return
	}

	asset, err := h.repo.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// DELETE /api/assets/:id
func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /api/assets/:id/assign
func (h *AssetHandler) AssignAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req models.AssignAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.AssignedTo) == "" {
		writeError(w, http.StatusBadRequest, "Field 'assigned_to' is required")
		return
	}

	asset, err := h.repo.Assign(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		if errors.Is(err, repository.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// POST /api/assets/:id/return
func (h *AssetHandler) ReturnAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req models.ReturnAssetRequest
	_ = json.NewDecoder(r.Body).Decode(&req) // optional body

	asset, err := h.repo.ReturnAsset(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		if errors.Is(err, repository.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// POST /api/assets/:id/repair
func (h *AssetHandler) RepairAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req models.StatusChangeRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	asset, err := h.repo.Repair(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		if errors.Is(err, repository.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// POST /api/assets/:id/retire
func (h *AssetHandler) RetireAsset(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req models.StatusChangeRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	asset, err := h.repo.Retire(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		if errors.Is(err, repository.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// GET /api/assets/:id/history
func (h *AssetHandler) GetAssetHistory(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	history, err := h.repo.GetHistory(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, history)
}

func getIDFromRequest(r *http.Request) (int, error) {
	// First check Go 1.22+ PathValue
	if val := r.PathValue("id"); val != "" {
		return strconv.Atoi(val)
	}

	// Fallback to URL path segments
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, part := range parts {
		if part == "assets" && i+1 < len(parts) {
			return strconv.Atoi(parts[i+1])
		}
	}

	if len(parts) > 0 {
		return strconv.Atoi(parts[len(parts)-1])
	}

	return 0, errors.New("cannot extract id")
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
