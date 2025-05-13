package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type doctorQueryResponse struct {
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func (d *Store) GetDoctorById(ctx context.Context, id string) (interface{}, error) {
	var queryData doctorQueryResponse

	key := fmt.Sprintf("doctor:id:%s", id)
	cached, err := d.rdb.Get(ctx, key).Result()
	if err == nil {
		log.Printf("Cache hit for doctor:id:%s", id)
		// Cache hit - deserialize JSON
		err := json.Unmarshal([]byte(cached), &queryData)
		if err != nil {
			return nil, gqlerror.Errorf("%s", err)
		}

		return queryData, nil
	}

	// Cache miss - fetch from DB
	log.Printf("Cache miss for doctor:id:%s", id)

	err = d.db.QueryRowContext(ctx, "SELECT fullname, email, specialization, created_at, updated_at FROM doctor WHERE doctor_id=$1", id).Scan(
		&queryData.Fullname, &queryData.Email, &queryData.Specialization, &queryData.CreatedAt, &queryData.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queryData, nil // return empty model
		}
		return queryData, err // return empty model
	}

	// Store in Redis for 10 minutes
	jsonData, _ := json.Marshal(queryData)
	log.Printf("Cache store for doctor:id:%s", id)
	d.rdb.Set(ctx, key, jsonData, 10*time.Minute)

	return queryData, err
}
