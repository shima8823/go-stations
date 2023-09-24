package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	osname, ok := r.Context().Value(middleware.OSNameKey).(string)
	if ok {
		println(osname)
	}

	if r.Method == http.MethodPost {
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if (req.ID == 0 || req.Subject == "") && (req.Subject == "" || req.Description == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Update(r.Context(), &req)
		if err != nil {
			if err.(*model.ErrNotFound) != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		PrevID := r.URL.Query().Get("prev_id")
		Size := r.URL.Query().Get("size")

		if PrevID == "" {
			PrevID = "0"
		}

		if Size == "" {
			Size = "0"
		}

		PrevIDInt, err := strconv.ParseInt(PrevID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		SizeInt, err := strconv.ParseInt(Size, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req := model.ReadTODORequest{
			PrevID: PrevIDInt,
			Size:   SizeInt,
		}

		res, err := h.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := map[string]interface{}{
			"todos": res,
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodDelete {
		req := model.DeleteTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
