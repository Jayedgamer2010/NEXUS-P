package services

import (
	"errors"
	"fmt"

	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EggService struct {
	repo       *repositories.EggRepository
	db         *gorm.DB
	serverRepo *repositories.ServerRepository
}

func NewEggService(
	eggRepo *repositories.EggRepository,
	serverRepo *repositories.ServerRepository,
	db *gorm.DB,
) *EggService {
	return &EggService{repo: eggRepo, serverRepo: serverRepo, db: db}
}

func (s *EggService) All(page, perPage int) ([]models.Egg, int64, error) {
	return s.repo.All(page, perPage)
}

func (s *EggService) FindByID(id uint) (*models.Egg, error) {
	return s.repo.FindByID(id)
}

func (s *EggService) Create(name, description, dockerImage, startup, author string) (*models.Egg, error) {
	egg := &models.Egg{
		UUID:        uuid.New().String(),
		Name:        name,
		Description: description,
		DockerImage: dockerImage,
		Startup:     startup,
		Author:      author,
	}

	if err := s.repo.Create(egg); err != nil {
		return nil, fmt.Errorf("failed to create egg: %w", err)
	}

	return egg, nil
}

func (s *EggService) Update(egg *models.Egg, name, description, dockerImage, startup string) error {
	if name != "" {
		egg.Name = name
	}
	if description != "" {
		egg.Description = description
	}
	if dockerImage != "" {
		egg.DockerImage = dockerImage
	}
	if startup != "" {
		egg.Startup = startup
	}

	return s.repo.Update(egg)
}

func (s *EggService) Delete(egg *models.Egg) error {
	var serverCount int64
	s.db.Model(&models.Server{}).Where("egg_id = ?", egg.ID).Count(&serverCount)
	if serverCount > 0 {
		return errors.New("cannot delete egg with active servers")
	}

	return s.repo.Delete(egg)
}

func (s *EggService) SeedDefaults() error {
	count, _ := s.repo.Count()
	if count > 0 {
		return nil // already seeded
	}

	defaults := []struct {
		Name, Description, DockerImage, Startup string
	}{
		{
			Name:        "Vanilla Minecraft",
			Description: "Vanilla Minecraft Java Edition",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar server.jar",
		},
		{
			Name:        "Paper Minecraft",
			Description: "PaperMC - High performance Minecraft server",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar paper.jar",
		},
		{
			Name:        "Bungeecord",
			Description: "BungeeCord Minecraft proxy server",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar bungeecord.jar",
		},
	}

	for _, d := range defaults {
		if _, err := s.Create(d.Name, d.Description, d.DockerImage, d.Startup, "NEXUS"); err != nil {
			// Log but continue - seed failures aren't fatal
			fmt.Printf("Warning: failed to seed egg %s: %v\n", d.Name, err)
		}
	}

	_ = utils.GenerateUUID() // ensure utils package is imported
	return nil
}
