package course

import (
	"log"
	"time"

	"github.com/tzincker/gocourse_web/internal/domain"
)

type (
	Filters struct {
		Name      string
		StartDate string
		EndDate   string
	}

	Service interface {
		Create(name, startDate, endDate string) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate *string, endDate *string) error
		Count(filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name, startDate, endDate string) (*domain.Course, error) {
	log.Println("Create course service")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Panicln(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Panicln(err)
		return nil, err
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	c, err := s.repo.Create(&course)

	if err != nil {
		s.log.Println(err)
	}

	return c, err
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	log.Println("Get all courses service")

	courses, err := s.repo.GetAll(filters, offset, limit)

	if err != nil {
		s.log.Println(err)
	}

	return courses, err
}

func (s service) Get(id string) (*domain.Course, error) {
	log.Println("Get course service")

	c, err := s.repo.Get(id)

	if err != nil {
		s.log.Println(err)
	}

	return c, err
}

func (s service) Delete(id string) error {
	log.Println("Delete course service")

	err := s.repo.Delete(id)

	if err != nil {
		s.log.Println(err)
	}

	return nil
}

func (s service) Update(
	id string,
	name *string,
	startDate *string,
	endDate *string,
) error {
	log.Println("Update course service")

	var err error
	var startDateParsed *time.Time
	var endDateParsed *time.Time

	if startDate != nil {

		parsedDate, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Panicln(err)
			return err
		}
		startDateParsed = &parsedDate
	}

	if endDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Panicln(err)
			return err
		}
		endDateParsed = &parsedDate
	}

	err = s.repo.Update(id, name, startDateParsed, endDateParsed)
	if err != nil {
		s.log.Println(err)
	}

	return nil
}

func (s service) Count(filters Filters) (int64, error) {
	log.Println("Get all courses count service")
	count, err := s.repo.Count(filters)
	if err != nil {
		s.log.Println(err)
	}

	return count, err
}
