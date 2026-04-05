package repositories

import (
	"nexus/backend/models"

	"gorm.io/gorm"
)

type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) FindAll() ([]models.Node, error) {
	var nodes []models.Node
	err := r.db.Find(&nodes).Error
	return nodes, err
}

func (r *NodeRepository) FindByID(id uint) (*models.Node, error) {
	var node models.Node
	err := r.db.Preload("Allocations").First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepository) Create(node *models.Node) error {
	return r.db.Create(node).Error
}

func (r *NodeRepository) Update(node *models.Node) error {
	return r.db.Save(node).Error
}

func (r *NodeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Node{}, id).Error
}

func (r *NodeRepository) CountAll() int64 {
	var count int64
	r.db.Model(&models.Node{}).Count(&count)
	return count
}
