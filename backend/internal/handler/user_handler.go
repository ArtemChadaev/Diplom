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
)

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
