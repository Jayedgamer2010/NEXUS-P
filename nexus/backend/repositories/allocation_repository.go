package repositories

import (
	"nexus/backend/models"
	"nexus/backend/utils"

	"gorm.io/gorm"
)

type AllocationRepository struct {
	db *gorm.DB
}

func NewAllocationRepository(db *gorm.DB) *AllocationRepository {
	return &AllocationRepository{db: db}
}

func (r *AllocationRepository) ByNodeID(nodeID uint) ([]models.Allocation, error) {
	var allocations []models.Allocation
	err := r.db.Where("node_id = ?", nodeID).Preload("Server").Find(&allocations).Error
	return allocations, err
}

func (r *AllocationRepository) FindUnassignedByNodeID(nodeID uint) ([]models.Allocation, error) {
	var allocations []models.Allocation
	err := r.db.Where("node_id = ? AND server_id IS NULL", nodeID).Find(&allocations).Error
	return allocations, err
}

func (r *AllocationRepository) FindUnassigned(nodeID uint) (*models.Allocation, error) {
	var allocation models.Allocation
	err := r.db.Where("node_id = ? AND server_id IS NULL", nodeID).First(&allocation).Error
	return &allocation, err
}

func (r *AllocationRepository) FindByID(id uint) (*models.Allocation, error) {
	var allocation models.Allocation
	err := r.db.First(&allocation, id).Error
	return &allocation, err
}

func (r *AllocationRepository) Create(allocation *models.Allocation) error {
	return r.db.Create(allocation).Error
}

func (r *AllocationRepository) Update(allocation *models.Allocation) error {
	return r.db.Save(allocation).Error
}

func (r *AllocationRepository) Delete(id uint) error {
	return r.db.Where("server_id IS NULL").Delete(&models.Allocation{}, id).Error
}

func (r *AllocationRepository) List(page, perPage int) ([]models.Allocation, int64, error) {
	page, perPage = utils.SanitizePagination(page, perPage)
	var allocations []models.Allocation
	var total int64

	r.db.Model(&models.Allocation{}).Count(&total)
	err := utils.Paginate(r.db, page, perPage).Preload("Server").Find(&allocations).Error
	return allocations, total, err
}
