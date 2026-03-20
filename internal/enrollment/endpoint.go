package enrollment

import (
	"encoding/json"
	"net/http"

	"github.com/tzincker/gocourse_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int32      `json:"status"`
		Data   any        `json:"data,omitempty"`
		Err    string     `json:"error,omitempty"`
		Meta   *meta.Meta `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
			return
		}

		if req.UserID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "user_id required"})
			return
		}

		if req.CourseID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "course_id required"})
			return
		}

		enroll, err := s.Create(req.UserID, req.CourseID)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: enroll})
	}
}
