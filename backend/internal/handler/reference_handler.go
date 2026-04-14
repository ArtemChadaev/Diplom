package handler

import (
	"net/http"

	"github.com/ima/diplom-backend/internal/handler/dto"
)

// listCountries godoc
// @Summary      List countries
// @Description  Returns a list of all countries in ISO alpha-3 format
// @Tags         References
// @Produce      json
// @Success      200  {array}   dto.CountryResponse
// @Router       /api/v1/ref/countries [get]
func (h *Handler) listCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := h.service.Reference.GetCountries(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch countries")
		return
	}

	resp := make([]dto.CountryResponse, len(countries))
	for i, c := range countries {
		resp[i] = dto.FromCountryDomain(c)
	}

	writeJSON(w, http.StatusOK, resp)
}

// searchATC godoc
// @Summary      Search ATC codes
// @Description  Search for ATC (Anatomical Therapeutic Chemical) codes by partial code or name
// @Tags         References
// @Produce      json
// @Param        q    query     string  false  "Search query"
// @Success      200  {array}   dto.ATCCodeResponse
// @Router       /api/v1/ref/atc [get]
func (h *Handler) searchATC(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	codes, err := h.service.Reference.SearchATC(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to search ATC codes")
		return
	}

	resp := make([]dto.ATCCodeResponse, len(codes))
	for i, c := range codes {
		resp[i] = dto.FromATCCodeDomain(c)
	}

	writeJSON(w, http.StatusOK, resp)
}
