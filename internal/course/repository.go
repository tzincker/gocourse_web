package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tzincker/gocourse_web/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *domain.Course) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate *time.Time, endDate *time.Time) error
		Count(filters Filters) (int64, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(course *domain.Course) (*domain.Course, error) {
	if err := r.db.Create(course).Error; err != nil {
		r.log.Printf("error: %v", err)
		return nil, err
	}

	r.log.Println("course created with id: ", course.ID)
	return course, nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	var courses []domain.Course
	tx := repo.db.Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&courses)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("courses got")
	return courses, nil
}

func (repo *repo) Get(id string) (*domain.Course, error) {
	course := domain.Course{ID: id}

	result := repo.db.First(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	repo.log.Println("course found with id: ", course.ID)
	return &course, nil
}

func (repo *repo) Delete(id string) error {
	course := domain.Course{ID: id}

	result := repo.db.Delete(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	repo.log.Println("course deleted with id: ", course.ID)
	return nil
}

func (repo *repo) Update(
	id string,
	name *string,
	startDate *time.Time,
	endDate *time.Time,
) error {

	values := make(map[string]any)

	if name != nil {
		values["first_name"] = name
	}

	if startDate != nil {
		values["last_name"] = startDate
	}

	if endDate != nil {
		values["email"] = endDate
	}

	result := repo.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	repo.log.Println("course updated with id: ", id)
	return nil
}

func (repo *repo) Count(filters Filters) (int64, error) {
	var count int64
	tx := repo.db.Model(domain.Course{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	repo.log.Println("courses count")
	return count, nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
