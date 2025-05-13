package store

import (
	"context"
	"database/sql"
	"errors"
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

	err := d.db.QueryRowContext(ctx, "SELECT fullname, email, specialization, created_at, updated_at FROM doctor WHERE doctor_id=$1", id).Scan(
		&queryData.Fullname, &queryData.Email, &queryData.Specialization, &queryData.CreatedAt, &queryData.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queryData, nil // return empty model
		}
		return queryData, err // return empty model
	}
	return queryData, err
}
