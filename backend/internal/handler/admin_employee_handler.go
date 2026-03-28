package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

// adminGetEmployeeProfile handles GET /api/v1/admin/employees/{userID}
func (h *Handler) adminGetEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	targetUserID := mustParseIntParam(w, r, "userID")
	if targetUserID == -1 {
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	profile, err := h.service.EmployeeProfile.GetProfile(r.Context(), callerID, callerRole, targetUserID)
	if err != nil {
		handleEmployeeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, profileToResponse(profile))
}

// adminListEmployeeProfiles handles GET /api/v1/admin/employees
func (h *Handler) adminListEmployeeProfiles(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	profiles, err := h.service.EmployeeProfile.ListProfiles(r.Context(), callerID, callerRole, limit, offset)
	if err != nil {
		handleEmployeeError(w, r, err)
		return
	}

	responses := make([]dto.EmployeeProfileResponse, len(profiles))
	for i, p := range profiles {
		responses[i] = *profileToResponse(&p)
	}

	writeJSON(w, http.StatusOK, responses)
}

// adminPatchEmployeeProfile handles PATCH /api/v1/admin/employees/{userID}
func (h *Handler) adminPatchEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	targetUserID := mustParseIntParam(w, r, "userID")
	if targetUserID == -1 {
		return
	}

	var req dto.PatchEmployeeProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	input := domain.UpdateEmployeeProfileInput{
		FullName:           req.FullName,
		CorporateEmail:     req.CorporateEmail,
		Phone:              req.Phone,
		Position:           req.Position,
		Department:         req.Department,
		BirthDate:          req.BirthDate,
		AvatarURL:          req.AvatarURL,
		HireDate:           req.HireDate,
		DismissalDate:      req.DismissalDate,
		MedicalBookScanURL: req.MedicalBookScanURL,
		SpecialZoneAccess:  req.SpecialZoneAccess,
		GDPTrainingHistory: req.GDPTrainingHistory,
	}

	updated, err := h.service.EmployeeProfile.UpdateProfile(r.Context(), callerID, callerRole, targetUserID, input)
	if err != nil {
		handleEmployeeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, profileToResponse(updated))
}

func handleEmployeeError(w http.ResponseWriter, r *http.Request, err error) {
	log := logger.FromContext(r.Context())
	if errors.Is(err, domain.ErrEmployeeProfileNotFound) {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, domain.ErrInsufficientPerms) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	log.Error("employee profile operation failed", "error", err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func mustParseIntParam(w http.ResponseWriter, r *http.Request, param string) int {
	val := chi.URLParam(r, param)
	id, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return -1
	}
	return id
}

func profileToResponse(p *domain.EmployeeProfile) *dto.EmployeeProfileResponse {
	return &dto.EmployeeProfileResponse{
		ID:                 p.ID,
		UserID:             p.UserID,
		EmployeeCode:       p.EmployeeCode,
		FullName:           p.FullName,
		CorporateEmail:     p.CorporateEmail,
		Phone:              p.Phone,
		Position:           p.Position,
		Department:         p.Department,
		BirthDate:          p.BirthDate,
		AvatarURL:          p.AvatarURL,
		HireDate:           p.HireDate,
		DismissalDate:      p.DismissalDate,
		MedicalBookScanURL: p.MedicalBookScanURL,
		SpecialZoneAccess:  p.SpecialZoneAccess,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
