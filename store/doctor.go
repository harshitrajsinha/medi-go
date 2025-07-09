package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type doctorQueryResponse struct {
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// Queries doctor details based on doctor id
func (d *Store) GetDoctorById(id string) (interface{}, error) {
	var queryData doctorQueryResponse
	var err error
	var redisKey string

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second) // if database takes too long, the query should be cancelled automatically after 45 seconds
	defer cancel()

	redisCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // if redis takes too long, the query should be cancelled automatically after 15 seconds
	defer cancel()

	if d.rdb != nil {

		redisKey = fmt.Sprintf("doctor:id:%s", id)
		cached, err := d.rdb.Get(redisCtx, redisKey).Result()
		if err == nil {
			log.Printf("Cache hit for doctor:id:%s", id)
			// Cache hit - deserialize JSON
			err := json.Unmarshal([]byte(cached), &queryData)
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}

			return queryData, nil
		}

		// Cache miss - fetch from DB
		log.Printf("Cache miss for doctor:id:%s", id)
	}

	err = d.db.QueryRowContext(ctx, "SELECT fullname, email, specialization, created_at, updated_at FROM doctor WHERE doctor_id=$1", id).Scan(
		&queryData.Fullname, &queryData.Email, &queryData.Specialization, &queryData.CreatedAt, &queryData.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queryData, nil // return empty model
		}
		return queryData, err // return empty model
	}

	if d.rdb != nil {
		// Store in Redis for 60 minutes (sice doctor data unlikely to update frequently)
		jsonData, _ := json.Marshal(queryData)
		log.Printf("Cache store for doctor:id:%s", id)
		d.rdb.Set(ctx, redisKey, jsonData, 60*time.Minute)
	}

	return queryData, err
}
