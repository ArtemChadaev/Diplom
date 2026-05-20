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

// adminCreateEmployeeProfile handles POST /api/v1/admin/employees
// adminCreateEmployeeProfile godoc
// @Summary      Create employee profile (Admin)
// @Description  Creates a new employee profile. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body    body      dto.CreateEmployeeProfileRequest  true  "Fields to create"
// @Success      201     {object}  dto.EmployeeProfileResponse
// @Failure      400     {object}  dto.ErrorResponse  "invalid JSON or invalid fields"
// @Failure      403     {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      500     {object}  dto.ErrorResponse  "internal error"
// @Router       /api/v1/admin/employees [post]
func (h *Handler) adminCreateEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	var gdpRaw json.RawMessage
	if req.GDPTrainingHistory != nil {
		if b, err := json.Marshal(req.GDPTrainingHistory); err == nil {
			gdpRaw = b
		}
	}

	input := domain.CreateEmployeeProfileInput{
		UserID:             req.UserID,
		EmployeeCode:       req.EmployeeCode,
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
		GDPTrainingHistory: gdpRaw,
	}

	created, err := h.service.EmployeeProfile.CreateProfile(r.Context(), callerID, callerRole, input)
	if err != nil {
		handleEmployeeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, profileToResponse(created))
}

// adminGetEmployeeProfile handles GET /api/v1/admin/employees/{userID}
// adminGetEmployeeProfile godoc
// @Summary      Get employee profile (Admin)
// @Description  Returns the details of a specific employee profile. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        userID  path      int  true  "User ID of the employee"
// @Success      200     {object}  dto.EmployeeProfileResponse
// @Failure      400     {object}  dto.ErrorResponse  "invalid user ID"
// @Failure      403     {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404     {object}  dto.ErrorResponse  "profile not found"
// @Router       /api/v1/admin/employees/{userID} [get]
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
// adminListEmployeeProfiles godoc
// @Summary      List employee profiles (Admin)
// @Description  Returns a list of employee profiles with optional pagination. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        limit   query     int  false  "Limit (default 10)"
// @Param        offset  query     int  false  "Offset (default 0)"
// @Success      200     {array}   dto.EmployeeProfileResponse
// @Failure      403     {object}  dto.ErrorResponse  "insufficient permissions"
// @Router       /api/v1/admin/employees [get]
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
// adminPatchEmployeeProfile godoc
// @Summary      Update employee profile (Admin)
// @Description  Partially updates fields in an employee's profile. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        userID  path      int  true  "User ID of the employee"
// @Param        body    body      dto.PatchEmployeeProfileRequest  true  "Fields to update"
// @Success      200     {object}  dto.EmployeeProfileResponse
// @Failure      400     {object}  dto.ErrorResponse  "invalid user ID, invalid JSON or invalid fields"
// @Failure      403     {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404     {object}  dto.ErrorResponse  "profile not found"
// @Router       /api/v1/admin/employees/{userID} [patch]
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

	var gdpRaw json.RawMessage
	if req.GDPTrainingHistory != nil {
		if b, err := json.Marshal(req.GDPTrainingHistory); err == nil {
			gdpRaw = b
		}
	}

	input := domain.UpdateEmployeeProfileInput{
		EmployeeCode:       req.EmployeeCode,
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
		GDPTrainingHistory: gdpRaw,
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
	if errors.Is(err, domain.ErrInvalidEmployeeCode) {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	var gdpRaw json.RawMessage
	if len(p.GDPTrainingHistory) > 0 {
		if bytes, err := json.Marshal(p.GDPTrainingHistory); err == nil {
			gdpRaw = bytes
		}
	}

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
		GDPTrainingHistory: gdpRaw,
	}
}