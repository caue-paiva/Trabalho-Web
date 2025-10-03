package mapper

import (
	"time"

	"backend/internal/entities"
)

// Text DTOs

type CreateTextRequest struct {
	Slug     string `json:"slug"`
	Content  string `json:"content"`
	PageID   string `json:"page_id,omitempty"`
	PageSlug string `json:"page_slug,omitempty"`
}

type UpdateTextRequest struct {
	Content  string `json:"content"`
	PageID   string `json:"page_id,omitempty"`
	PageSlug string `json:"page_slug,omitempty"`
}

type TextResponse struct {
	ID            string    `json:"id"`
	Slug          string    `json:"slug"`
	Content       string    `json:"content"`
	PageID        string    `json:"page_id,omitempty"`
	PageSlug      string    `json:"page_slug,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastUpdatedBy string    `json:"last_updated_by,omitempty"`
}

// Mapping functions

func ToTextEntity(req CreateTextRequest) entities.Text {
	return entities.Text{
		Slug:     req.Slug,
		Content:  req.Content,
		PageID:   req.PageID,
		PageSlug: req.PageSlug,
	}
}

func ToTextUpdateEntity(req UpdateTextRequest) entities.Text {
	return entities.Text{
		Content:  req.Content,
		PageID:   req.PageID,
		PageSlug: req.PageSlug,
	}
}

func TextToResponse(text entities.Text) TextResponse {
	return TextResponse{
		ID:            text.ID,
		Slug:          text.Slug,
		Content:       text.Content,
		PageID:        text.PageID,
		PageSlug:      text.PageSlug,
		CreatedAt:     text.CreatedAt,
		UpdatedAt:     text.UpdatedAt,
		LastUpdatedBy: text.LastUpdatedBy,
	}
}

func TextsToResponse(texts []entities.Text) []TextResponse {
	result := make([]TextResponse, len(texts))
	for i, text := range texts {
		result[i] = TextToResponse(text)
	}
	return result
}
