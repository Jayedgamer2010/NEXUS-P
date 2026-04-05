package repositories

import (
	"nexus/backend/models"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) DB() *gorm.DB {
	return r.db
}

func (r *ServerRepository) FindAll(page, perPage int) ([]models.Server, int64, error) {
	var servers []models.Server
	var total int64

	if err := r.db.Model(&models.Server{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := r.db.Preload("User").Preload("Node").Preload("Egg").Preload("Allocation").
		Offset(offset).Limit(perPage).Order("id DESC").Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	return servers, total, nil
}

func (r *ServerRepository) FindByUserID(userID uint) ([]models.Server, error) {
	var servers []models.Server
	err := r.db.Where("user_id = ?", userID).
		Preload("Node").Preload("Egg").Preload("Allocation").
		Order("id DESC").Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) FindByID(id uint) (*models.Server, error) {
	var server models.Server
	err := r.db.Preload("User").Preload("Node").Preload("Egg").Preload("Allocation").
		First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) FindByUUID(uuid string) (*models.Server, error) {
	var server models.Server
	err := r.db.Preload("User").Preload("Node").Preload("Egg").Preload("Allocation").
		Where("uuid = ?", uuid).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) Create(server *models.Server) error {
	return r.db.Create(server).Error
}

func (r *ServerRepository) Update(server *models.Server) error {
	return r.db.Save(server).Error
}

func (r *ServerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Server{}, id).Error
}

func (r *ServerRepository) CountAll() int64 {
	var count int64
	r.db.Model(&models.Server{}).Count(&count)
	return count
}

func (r *ServerRepository) CountRunning() int64 {
	var count int64
	r.db.Model(&models.Server{}).Where("status = ?", models.StatusRunning).Count(&count)
	return count
}

func (r *ServerRepository) FindRecent(limit int) ([]models.Server, error) {
	var servers []models.Server
	err := r.db.Preload("User").Preload("Node").
		Order("created_at DESC").Limit(limit).Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) CountByUserID(userID uint) int64 {
	var count int64
	r.db.Model(&models.Server{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (r *ServerRepository) CountByNodeID(nodeID uint) int64 {
	var count int64
	r.db.Model(&models.Server{}).Where("node_id = ?", nodeID).Count(&count)
	return count
}
