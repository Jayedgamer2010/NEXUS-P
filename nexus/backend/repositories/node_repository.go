package repositories

import (
	"errors"

	"nexus/backend/models"
	"nexus/backend/utils"

	"gorm.io/gorm"
)

type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) All(page, perPage int) ([]models.Node, int64, error) {
	page, perPage = utils.SanitizePagination(page, perPage)
	var nodes []models.Node
	var total int64

	r.db.Model(&models.Node{}).Count(&total)
	err := utils.Paginate(r.db, page, perPage).Find(&nodes).Error
	return nodes, total, err
}

func (r *NodeRepository) FindByID(id uint) (*models.Node, error) {
	var node models.Node
	err := r.db.First(&node, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &node, err
}

func (r *NodeRepository) Create(node *models.Node) error {
	return r.db.Create(node).Error
}

func (r *NodeRepository) Update(node *models.Node) error {
	return r.db.Save(node).Error
}

func (r *NodeRepository) Delete(node *models.Node) error {
	return r.db.Delete(node).Error
}

func (r *NodeRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&models.Node{}).Count(&total).Error
	return total, err
}

func (r *NodeRepository) FindWithAllocations(id uint) (*models.Node, error) {
	var node models.Node
	err := r.db.Preload("Allocations").First(&node, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &node, err
}
