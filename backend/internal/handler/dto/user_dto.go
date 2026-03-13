package dto

// IsLoginTakenRequest — входящий запрос для проверки доступности логина.
// Теги json и validate живут только здесь, не в domain.
type IsLoginTakenRequest struct {
	Login string `json:"login" validate:"required,min=3,max=255"`
}

// IsLoginTakenResponse — ответ API.
type IsLoginTakenResponse struct {
	Taken bool `json:"taken"`
}
