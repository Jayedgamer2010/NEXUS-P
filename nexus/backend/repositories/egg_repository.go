package repositories

import (
	"nexus/backend/models"

	"gorm.io/gorm"
)

type EggRepository struct {
	db *gorm.DB
}

func NewEggRepository(db *gorm.DB) *EggRepository {
	return &EggRepository{db: db}
}

func (r *EggRepository) FindAll() ([]models.Egg, error) {
	var eggs []models.Egg
	err := r.db.Find(&eggs).Error
	return eggs, err
}

func (r *EggRepository) FindByID(id uint) (*models.Egg, error) {
	var egg models.Egg
	err := r.db.First(&egg, id).Error
	if err != nil {
		return nil, err
	}
	return &egg, nil
}

func (r *EggRepository) Create(egg *models.Egg) error {
	return r.db.Create(egg).Error
}

func (r *EggRepository) Update(egg *models.Egg) error {
	return r.db.Save(egg).Error
}

func (r *EggRepository) Delete(id uint) error {
	return r.db.Delete(&models.Egg{}, id).Error
}

func (r *EggRepository) CountByEggID(id uint) int64 {
	var count int64
	r.db.Model(&models.Server{}).Where("egg_id = ?", id).Count(&count)
	return count
}
