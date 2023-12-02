// internal/usecases/advertisement_service.go
package usecases

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"

	"advertisements-service/internal/entities"
	"advertisements-service/internal/repositories"
)

type AdvertisementRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	IsActive    bool    `json:"is_active"`
}

type AdvertisementService struct {
	Repo     *repositories.AdvertisementRepository
	validate *validator.Validate
}

func NewAdvertisementService(repo *repositories.AdvertisementRepository) *AdvertisementService {
	return &AdvertisementService{
		Repo:     repo,
		validate: validator.New(),
	}
}

func (s *AdvertisementService) GetAdvertisementsPage(page, pageSize int, sortBy, sortOrder string) ([]entities.Advertisement, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("invalid page or pageSize")
	}

	validSortColumns := map[string]bool{"price": true, "created_at": true}
	if !validSortColumns[sortBy] {
		return nil, errors.New("invalid sortBy parameter")
	}

	validSortOrders := map[string]bool{"asc": true, "desc": true}
	if !validSortOrders[strings.ToLower(sortOrder)] {
		return nil, errors.New("invalid sortOrder parameter")
	}

	offset := (page - 1) * pageSize

	advertisements, err := s.Repo.GetPaginatedAdvertisements(offset, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	return advertisements, nil
}

func (s *AdvertisementService) GetAllAdvertisements() ([]entities.Advertisement, error) {
	return s.Repo.GetAll()
}

func (s *AdvertisementService) GetAdvertisementByID(id int) (*entities.Advertisement, error) {
	return s.Repo.GetByID(id)
}


func (s *AdvertisementService) CreateAdvertisement(newAdvertisement *AdvertisementRequest) error {

	if err := s.validate.Struct(newAdvertisement); err != nil {
		return s.validationError(err)
	}

	ad := entities.Advertisement{
		Title:       newAdvertisement.Title,
		Description: newAdvertisement.Description,
		Price:       newAdvertisement.Price,
		IsActive:    newAdvertisement.IsActive,
	}

	return s.Repo.Create(&ad)
}


func (s *AdvertisementService) UpdateAdvertisement(id int, updatedAdvertisement *AdvertisementRequest) error {
	
	if err := s.validate.Struct(updatedAdvertisement); err != nil {
		return s.validationError(err)
	}

	existingAd, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingAd == nil {
		return errors.New("advertisement not found")
	}

	if updatedAdvertisement.Title != "" {
		existingAd.Title = updatedAdvertisement.Title
	}
	if updatedAdvertisement.Description != "" {
		existingAd.Description = updatedAdvertisement.Description
	}
	if updatedAdvertisement.Price != 0 {
		existingAd.Price = updatedAdvertisement.Price
	}

	return s.Repo.Update(id, existingAd)
}

func (s *AdvertisementService) DeleteAdvertisement(id int) error {
	return s.Repo.Delete(id)
}

func (s *AdvertisementService) validationError(err error) error {
	var validationErrors []string

	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, s.buildValidationErrorString(err))
	}

	return errors.New(strings.Join(validationErrors, ", "))
}

func (s *AdvertisementService) buildValidationErrorString(err validator.FieldError) string {
	fieldName := err.Field()

	return fieldName + " is required"
}

