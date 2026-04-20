package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ruby-107/source-asia/internal/model"
	"github.com/ruby-107/source-asia/internal/service"
)

type Handler struct {
	limiter *service.RateLimiter
	db      *sql.DB
}

func NewHandler(l *service.RateLimiter, db *sql.DB) *Handler {
	return &Handler{
		limiter: l,
		db:      db,
	}
}

func writeResponse(w http.ResponseWriter, status int, success bool, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]interface{}{
		"success": success,
	}

	if success {
		resp["data"] = data
	} else {
		resp["error"] = errMsg
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
	var req model.Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, nil, "invalid request")
		return
	}

	if req.UserID == 0 {
		writeResponse(w, http.StatusBadRequest, false, nil, "user_id is required")
		return
	}

	if !h.limiter.Allow(req.UserID) {
		writeResponse(w, http.StatusTooManyRequests, false, nil, "rate limit exceeded")
		return
	}

	service.JobQueue <- req

	data := map[string]interface{}{
		"user_id":    req.UserID,
		"payload":    req.Payload,
		"status":     "accepted",
		"created_at": time.Now(),
	}

	writeResponse(w, http.StatusAccepted, true, data, "")
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {

	rows, err := h.db.Query(`
		SELECT user_id, COUNT(*) 
		FROM production.users
		GROUP BY user_id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserStat struct {
		UserID       int `json:"user_id"`
		RequestCount int `json:"request_count"`
	}

	var stats []UserStat

	for rows.Next() {
		var stat UserStat

		err := rows.Scan(&stat.UserID, &stat.RequestCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats = append(stats, stat)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stats)
}
