package enrollment

import (
	"errors"
	"log"

	"github.com/tzincker/gocourse_web/internal/course"
	"github.com/tzincker/gocourse_web/internal/domain"
	"github.com/tzincker/gocourse_web/internal/user"
)

type (
	Service interface {
		Create(userID, courseID string) (*domain.Enrollment, error)
	}

	service struct {
		log       *log.Logger
		userSrv   user.Service
		courseSrv course.Service
		repo      Repository
	}
)

func NewService(log *log.Logger, userSrv user.Service, courseSrv course.Service, repo Repository) Service {
	return &service{
		log:       log,
		userSrv:   userSrv,
		courseSrv: courseSrv,
		repo:      repo,
	}
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	log.Println("Create enrollment service")
	enrollment := domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}

	if _, err := s.userSrv.Get(enrollment.UserID); err != nil {
		return nil, errors.New("user id does not exist")
	}

	if _, err := s.courseSrv.Get(enrollment.CourseID); err != nil {
		return nil, errors.New("course id does not exist")
	}

	e, err := s.repo.Create(&enrollment)

	if err != nil {
		s.log.Println(err)
	}

	return e, err
}
