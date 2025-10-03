package mapper

import (
	"encoding/base64"
	"fmt"
	"time"

	"backend/internal/entities"
)

// Image DTOs

type CreateImageRequest struct {
	Slug     string `json:"slug,omitempty"`
	Name     string `json:"name"`
	Text     string `json:"text,omitempty"`
	Date     string `json:"date,omitempty"` // ISO format
	Location string `json:"location,omitempty"`
	Data     string `json:"data"` // base64 encoded
}

type UpdateImageRequest struct {
	Slug     string `json:"slug,omitempty"`
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Date     string `json:"date,omitempty"` // ISO format
	Location string `json:"location,omitempty"`
	Data     string `json:"data,omitempty"` // base64 encoded (optional)
}

type ImageResponse struct {
	ID            string    `json:"id"`
	Slug          string    `json:"slug,omitempty"`
	ObjectURL     string    `json:"object_url"`
	Name          string    `json:"name"`
	Text          string    `json:"text,omitempty"`
	Date          time.Time `json:"date,omitempty"`
	Location      string    `json:"location,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastUpdatedBy string    `json:"last_updated_by,omitempty"`
}

// Mapping functions

func ToImageEntity(req CreateImageRequest) (entities.Image, []byte, error) {
	// Decode base64 image data
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return entities.Image{}, nil, fmt.Errorf("invalid base64 data: %w", err)
	}

	// Parse date if provided
	var date time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return entities.Image{}, nil, fmt.Errorf("invalid date format: %w", err)
		}
		date = parsedDate
	}

	img := entities.Image{
		Slug:     req.Slug,
		Name:     req.Name,
		Text:     req.Text,
		Date:     date,
		Location: req.Location,
	}

	return img, data, nil
}

func ToImageUpdateEntity(req UpdateImageRequest) (entities.Image, []byte, error) {
	// Decode base64 image data if provided
	var data []byte
	var err error
	if req.Data != "" {
		data, err = base64.StdEncoding.DecodeString(req.Data)
		if err != nil {
			return entities.Image{}, nil, fmt.Errorf("invalid base64 data: %w", err)
		}
	}

	// Parse date if provided
	var date time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return entities.Image{}, nil, fmt.Errorf("invalid date format: %w", err)
		}
		date = parsedDate
	}

	img := entities.Image{
		Slug:     req.Slug,
		Name:     req.Name,
		Text:     req.Text,
		Date:     date,
		Location: req.Location,
	}

	return img, data, nil
}

func ImageToResponse(img entities.Image) ImageResponse {
	return ImageResponse{
		ID:            img.ID,
		Slug:          img.Slug,
		ObjectURL:     img.ObjectURL,
		Name:          img.Name,
		Text:          img.Text,
		Date:          img.Date,
		Location:      img.Location,
		CreatedAt:     img.CreatedAt,
		UpdatedAt:     img.UpdatedAt,
		LastUpdatedBy: img.LastUpdatedBy,
	}
}

func ImagesToResponse(images []entities.Image) []ImageResponse {
	result := make([]ImageResponse, len(images))
	for i, img := range images {
		result[i] = ImageToResponse(img)
	}
	return result
}
