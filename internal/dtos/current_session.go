package dtos

import "github.com/google/uuid"

type CurrentSessionResponse struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
}

func ToCurrentSessionResponse(id, userId uuid.UUID, userAgent string) CurrentSessionResponse {
	return CurrentSessionResponse{
		ID:        id.String(),
		UserId:    userId.String(),
		UserAgent: userAgent,
	}
}
