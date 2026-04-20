package service

import (
	"sync"
	"time"
)

type UserData struct {
	Count     int
	Timestamp time.Time
}

type RateLimiter struct {
	mu    sync.Mutex
	store map[int]*UserData
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		store: make(map[int]*UserData),
	}
}

func (r *RateLimiter) Allow(userID int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	data, exists := r.store[userID]

	if !exists {
		r.store[userID] = &UserData{1, now}
		return true
	}

	if now.Sub(data.Timestamp) > time.Minute {
		data.Count = 1
		data.Timestamp = now
		return true
	}

	if data.Count >= 5 {
		return false
	}

	data.Count++
	return true
}
