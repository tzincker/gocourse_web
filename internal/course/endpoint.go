package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tzincker/gocourse_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
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
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
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

		if req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "firstname required"})
			return
		}

		if req.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "lastname required"})
			return
		}

		if req.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "email required"})
			return
		}

		usr, err := s.Create(req.Name, req.StartDate, req.EndDate)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: usr})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		course, err := s.Get(id)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		filters := Filters{
			Name:      v.Get("first_name"),
			StartDate: v.Get("start_date"),
			EndDate:   v.Get("end_date"),
		}

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: courses, Meta: meta})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		var req UpdateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
			return
		}

		err := s.Update(id, &req.Name, &req.StartDate, &req.EndDate)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: map[string]bool{"ok": true}})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: map[string]bool{"ok": true}})
	}
}
