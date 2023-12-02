// internal/repositories/advertisement_repository.go
package repositories

import (
	"database/sql"
	"log"

	"advertisements-service/internal/entities"
)

type AdvertisementRepository struct {
	DB *sql.DB
}

func NewAdvertisementRepository(db *sql.DB) *AdvertisementRepository {
	return &AdvertisementRepository{DB: db}
}

func (r *AdvertisementRepository) GetPaginatedAdvertisements(offset, pageSize int, sortBy, sortOrder string) ([]entities.Advertisement, error) {
	query := "SELECT * FROM advertisements ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"

	rows, err := r.DB.Query(query, pageSize, offset)
	if err != nil {
		log.Printf("Error querying paginated advertisements: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var advertisements []entities.Advertisement

	for rows.Next() {
		var ad entities.Advertisement
		err := rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
		if err != nil {
			log.Printf("Error scanning advertisement row: %v\n", err)
			return nil, err
		}
		advertisements = append(advertisements, ad)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over paginated advertisement rows: %v\n", err)
		return nil, err
	}

	return advertisements, nil
}

func (r *AdvertisementRepository) GetAll() ([]entities.Advertisement, error) {
	rows, err := r.DB.Query("SELECT * FROM advertisements")
	if err != nil {
		log.Printf("Error querying advertisements: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var advertisements []entities.Advertisement

	for rows.Next() {
		var ad entities.Advertisement
		err := rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
		if err != nil {
			log.Printf("Error scanning advertisement row: %v\n", err)
			return nil, err
		}
		advertisements = append(advertisements, ad)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over advertisement rows: %v\n", err)
		return nil, err
	}

	return advertisements, nil
}

func (r *AdvertisementRepository) GetByID(id int) (*entities.Advertisement, error) {
	row := r.DB.QueryRow("SELECT * FROM advertisements WHERE id = ?", id)

	var ad entities.Advertisement
	err := row.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil // No row found
	} else if err != nil {
		log.Printf("Error scanning advertisement row: %v\n", err)
		return nil, err
	}

	return &ad, nil
}

func (r *AdvertisementRepository) Create(advertisement *entities.Advertisement) error {
	log.Printf("%s",advertisement.Description)
	_, err := r.DB.Exec("INSERT INTO advertisements (title, description, price, is_active) VALUES (?, ?, ?, ?)",
		advertisement.Title, advertisement.Description, advertisement.Price, advertisement.IsActive)
	if err != nil {
		log.Printf("Error inserting advertisement: %v\n", err)
		return err
	}

	return nil
}

func (r *AdvertisementRepository) Update(id int, advertisement *entities.Advertisement) error {
	_, err := r.DB.Exec("UPDATE advertisements SET title=?, description=?, price=?, is_active=? WHERE id=?",
		advertisement.Title, advertisement.Description, advertisement.Price, advertisement.IsActive, id)
	if err != nil {
		log.Printf("Error updating advertisement: %v\n", err)
		return err
	}

	return nil
}

func (r *AdvertisementRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM advertisements WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting advertisement: %v\n", err)
		return err
	}

	return nil
}
