package service

import (
	"database/sql"
	"log"
	"time"

	"github.com/ruby-107/source-asia/internal/model"
)

var JobQueue = make(chan model.Request, 100)

func StartWorker(db *sql.DB) {
	for job := range JobQueue {
		processJob(db, job, 3)
	}
}

func processJob(db *sql.DB, job model.Request, retries int) {
	var err error

	for i := 0; i < retries; i++ {

		_, err = db.Exec(
			`INSERT INTO production.users (user_id, payload) 
			 VALUES ($1, $2)`,
			job.UserID,
			job.Payload,
		)

		if err == nil {
			log.Println("Inserted successfully for user:", job.UserID)
			return
		}

		log.Println("Retry:", i+1, "Error:", err)
		time.Sleep(1 * time.Second)
	}

	log.Println("❌ Job failed after retries:", err)
}
