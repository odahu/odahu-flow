package postgres

import (
	"github.com/odahu/odahu-flow/packages/policy-server/pkg/handler"
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

type OPAPolicyModel struct {
	gorm.Model
	ID string
	RegoPolicy []byte
	Labels map[string]string
	Tags []string
}

func (s *Store) Create(policy handler.OPAPolicy) error {
	model := OPAPolicyModel{
		ID:         policy.ID,
		RegoPolicy: policy.RegoPolicy,
		Labels:     policy.Labels,
	}
	result := s.DB.Create(&model)

	return result.Error
}

func (s *Store) Delete(ID string) error {
	panic("implement me")
}

func (s *Store) Update(ID string) error {
	panic("implement me")
}

func (s *Store) AutoMigrate() error {
	return s.DB.AutoMigrate(&OPAPolicyModel{})
}