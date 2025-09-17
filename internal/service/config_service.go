package service

import (
	"context"

	"github.com/arwoosa/form/conf"
	"github.com/arwoosa/form/internal/models"
)

// ConfigService handles configuration-related business logic
type ConfigService struct {
	config *conf.AppConfig
}

// NewConfigService creates a new config service
func NewConfigService(config *conf.AppConfig) *ConfigService {
	return &ConfigService{
		config: config,
	}
}

// GetBusinessConfig returns business-related configuration settings
func (s *ConfigService) GetBusinessConfig(ctx context.Context) (*models.BusinessConfig, error) {
	if s.config == nil || s.config.BusinessRulesConfig == nil {
		return nil, ErrInternalError
	}

	businessConfig := &models.BusinessConfig{
		MaxTemplatesPerMerchant: s.config.BusinessRulesConfig.MaxTemplatesPerMerchant,
	}

	return businessConfig, nil
}