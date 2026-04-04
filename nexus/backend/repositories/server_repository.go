package repositories

import (
	"errors"

	"nexus/backend/models"
	"nexus/backend/utils"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) All(page, perPage int) ([]models.Server, int64, error) {
	page, perPage = utils.SanitizePagination(page, perPage)
	var servers []models.Server
	var total int64

	r.db.Model(&models.Server{}).Count(&total)
	err := utils.Paginate(r.db, page, perPage).
		Preload("User").
		Preload("Node").
		Preload("Egg").
		Preload("Allocation").
		Find(&servers).Error
	return servers, total, err
}

func (r *ServerRepository) FindByID(id uint) (*models.Server, error) {
	var server models.Server
	err := r.db.Preload("User").Preload("Node").Preload("Egg").Preload("Allocation").
		First(&server, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &server, err
}

func (r *ServerRepository) FindByUUID(uuid string) (*models.Server, error) {
	var server models.Server
	err := r.db.Preload("User").Preload("Node").Preload("Egg").Preload("Allocation").
		Where("uuid = ?", uuid).First(&server).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &server, err
}

func (r *ServerRepository) FindByUserID(userID uint) ([]models.Server, error) {
	var servers []models.Server
	err := r.db.Where("user_id = ?", userID).Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&models.Server{}).Count(&total).Error
	return total, err
}

func (r *ServerRepository) CountRunning() (int64, error) {
	var total int64
	err := r.db.Model(&models.Server{}).Where("status = ?", models.StatusRunning).Count(&total).Error
	return total, err
}

func (r *ServerRepository) Create(server *models.Server) error {
	return r.db.Create(server).Error
}

func (r *ServerRepository) Update(server *models.Server) error {
	return r.db.Save(server).Error
}

func (r *ServerRepository) UpdateStatus(uuid string, status string) error {
	return r.db.Model(&models.Server{}).Where("uuid = ?", uuid).Update("status", status).Error
}

func (r *ServerRepository) Delete(server *models.Server) error {
	return r.db.Delete(server).Error
}
