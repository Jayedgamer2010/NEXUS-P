package repositories

import (
	"errors"
	"nexus/backend/models"

	"gorm.io/gorm"
)

type AllocationRepository struct {
	db *gorm.DB
}

var ErrAllocationAssigned = errors.New("allocation is already assigned to a server")

func NewAllocationRepository(db *gorm.DB) *AllocationRepository {
	return &AllocationRepository{db: db}
}

func (r *AllocationRepository) FindByNodeID(nodeID uint) ([]models.Allocation, error) {
	var allocations []models.Allocation
	err := r.db.Where("node_id = ?", nodeID).Find(&allocations).Error
	return allocations, err
}

func (r *AllocationRepository) FindAvailable(nodeID uint) (*models.Allocation, error) {
	var allocation models.Allocation
	err := r.db.Where("node_id = ? AND server_id IS NULL", nodeID).First(&allocation).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

func (r *AllocationRepository) Assign(allocationID uint, serverID uint) error {
	return r.db.Model(&models.Allocation{}).Where("id = ?", allocationID).
		Update("server_id", serverID).Error
}

func (r *AllocationRepository) Unassign(allocationID uint) error {
	return r.db.Model(&models.Allocation{}).Where("id = ?", allocationID).
		Update("server_id", nil).Error
}

func (r *AllocationRepository) Create(allocation *models.Allocation) error {
	return r.db.Create(allocation).Error
}

func (r *AllocationRepository) Delete(id uint) error {
	var alloc models.Allocation
	if err := r.db.First(&alloc, id).Error; err != nil {
		return err
	}
	if alloc.ServerID != nil {
		return ErrAllocationAssigned
	}
	return r.db.Delete(&models.Allocation{}, id).Error
}
