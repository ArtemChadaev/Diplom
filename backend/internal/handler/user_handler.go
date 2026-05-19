package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// getMe godoc
// @Summary      Get current user profile
// @Description  Returns profile details for the authenticated user
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.UserProfileResponse
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse  "user not found"
// @Router       /api/v1/users/me [get]
func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userRepo.FindProfileByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := mapToUserProfileResponse(profile)
	writeJSON(w, http.StatusOK, resp)
}

// patchMe godoc
// @Summary      Update current user profile
// @Description  Partially updates the authenticated user's employee profile
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PatchMeRequest  true  "Fields to update"
// @Success      200   {object}  dto.UserProfileResponse
// @Failure      400   {object}  dto.ErrorResponse  "invalid request body"
// @Failure      401   {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      404   {object}  dto.ErrorResponse  "profile not found"
// @Router       /api/v1/users/me [patch]
func (h *Handler) patchMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.PatchMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var corporateEmail *string
	if req.CorporateMail != nil {
		corporateEmail = req.CorporateMail
	} else {
		corporateEmail = req.CorporateEmail
	}

	var birthDate *time.Time
	if req.BirthdayDate != nil {
		birthDate = req.BirthdayDate
	} else {
		birthDate = req.BirthDate
	}

	var avatarURL *string
	if req.AvatarURL2 != nil {
		avatarURL = req.AvatarURL2
	} else {
		avatarURL = req.AvatarURL
	}

	input := domain.UpdateEmployeeProfileInput{
		FullName:           req.FullName,
		Phone:              req.Phone,
		CorporateEmail:     corporateEmail,
		BirthDate:          birthDate,
		AvatarURL:          avatarURL,
		MedicalBookScanURL: req.MedicalBookScanURL,
		GDPTrainingHistory: req.GDPTrainingHistory,
	}

	_, err := h.service.EmployeeProfile.PatchSelfProfile(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeProfileNotFound) {
			writeError(w, http.StatusNotFound, "employee profile not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	profile, err := h.userRepo.FindProfileByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, mapToUserProfileResponse(profile))
}

// listUsers godoc
// @Summary      List users (Admin)
// @Description  Returns a paginated and filtered list of all users and profiles. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        q      query     string  false  "Search query (login, email, full name)"
// @Param        role   query     string  false  "Filter by role"
// @Param        page   query     int     false  "Page number"
// @Param        limit  query     int     false  "Limit"
// @Success      200    {array}   dto.UserProfileResponse
// @Failure      403    {object}  dto.ErrorResponse  "insufficient permissions"
// @Router       /api/v1/admin/users [get]
func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	role := r.URL.Query().Get("role")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	filter := domain.UserListFilter{
		Query: q,
		Role:  domain.UserRole(role),
		Page:  page,
		Limit: limit,
	}

	profiles, total, err := h.userRepo.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))

	var resp []dto.UserProfileResponse
	for _, p := range profiles {
		resp = append(resp, mapToUserProfileResponse(p))
	}
	if resp == nil {
		resp = make([]dto.UserProfileResponse, 0)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getUserByID godoc
// @Summary      Get user by ID (Admin)
// @Description  Returns a user's details and profile. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserProfileResponse
// @Failure      400  {object}  dto.ErrorResponse  "invalid user ID"
// @Failure      403  {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404  {object}  dto.ErrorResponse  "user not found"
// @Router       /api/v1/admin/users/{id} [get]
func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	profile, err := h.userRepo.FindProfileByUserID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, mapToUserProfileResponse(profile))
}

// patchUser godoc
// @Summary      Update user and permissions (Admin)
// @Description  Updates user role, ns/pv access, special zone access. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "User ID"
// @Param        body  body      dto.PatchUserRequest  true  "Fields to update"
// @Success      200   {object}  dto.UserProfileResponse
// @Failure      400   {object}  dto.ErrorResponse  "invalid user ID or request body"
// @Failure      403   {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404   {object}  dto.ErrorResponse  "user not found"
// @Router       /api/v1/admin/users/{id} [patch]
func (h *Handler) patchUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Make sure the user exists first
	_, err = h.userRepo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if req.Role != nil {
		if err := h.userRepo.UpdateRole(r.Context(), id, domain.UserRole(*req.Role)); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update role")
			return
		}
	}

	if req.NsPvAccess != nil {
		if err := h.userRepo.SetNsPvAccess(r.Context(), id, *req.NsPvAccess); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update ns pv access")
			return
		}
	}

	if req.SpecialZoneAccess != nil {
		input := domain.UpdateEmployeeProfileInput{
			SpecialZoneAccess: req.SpecialZoneAccess,
		}
		if _, err := h.employeeProfileRepo.Update(r.Context(), id, input); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update special zone access")
			return
		}
	}

	// Return updated profile
	profile, err := h.userRepo.FindProfileByUserID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, mapToUserProfileResponse(profile))
}

func mapToUserProfileResponse(p *domain.UserProfile) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:           p.ID,
		Email:        p.Email,
		Role:         string(p.Role),
		NsPvAccess:   p.NsPvAccess,
		UkepBound:    p.UkepBound,
		IsBlocked:    p.IsBlocked,
		EmployeeCode: p.EmployeeCode,
		FullName:     p.FullName,
		Position:     p.Position,
		Department:   p.Department,
	}
}
