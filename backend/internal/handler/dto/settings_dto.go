package dto

type SystemSettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateSettingRequest struct {
	Value string `json:"value" validate:"required"`
}
