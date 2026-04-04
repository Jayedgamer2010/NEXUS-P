package repositories

import (
	"errors"

	"nexus/backend/models"
	"nexus/backend/utils"

	"gorm.io/gorm"
)

type EggRepository struct {
	db *gorm.DB
}

func NewEggRepository(db *gorm.DB) *EggRepository {
	return &EggRepository{db: db}
}

func (r *EggRepository) All(page, perPage int) ([]models.Egg, int64, error) {
	page, perPage = utils.SanitizePagination(page, perPage)
	var eggs []models.Egg
	var total int64

	r.db.Model(&models.Egg{}).Count(&total)
	err := utils.Paginate(r.db, page, perPage).Find(&eggs).Error
	return eggs, total, err
}

func (r *EggRepository) FindByID(id uint) (*models.Egg, error) {
	var egg models.Egg
	err := r.db.First(&egg, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &egg, err
}

func (r *EggRepository) Create(egg *models.Egg) error {
	return r.db.Create(egg).Error
}

func (r *EggRepository) Update(egg *models.Egg) error {
	return r.db.Save(egg).Error
}

func (r *EggRepository) Delete(egg *models.Egg) error {
	return r.db.Delete(egg).Error
}

func (r *EggRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&models.Egg{}).Count(&total).Error
	return total, err
}
