package enrollment

import (
	"log"

	"github.com/tzincker/gocourse_web/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enrollment *domain.Enrollment) (*domain.Enrollment, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(enrollment *domain.Enrollment) (*domain.Enrollment, error) {
	result := repo.db.Create(enrollment)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("enrollment created with id: ", enrollment.ID)
	return enrollment, nil
}
