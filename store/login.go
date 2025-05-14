package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/harshitrajsinha/medi-go/models"
)

func (rec *Store) GetLoginInfo(ctx context.Context, loginReq *models.Credentials) (LoginResponse, error) {

	loginResponse := LoginResponse{}
	var err error

	if loginReq.Role == "doctor" {
		err = rec.db.QueryRowContext(ctx, "SELECT doctor_id, password_hash FROM doctor WHERE role='doctor' AND email=$1", loginReq.Email).Scan(&loginResponse.UserID, &loginResponse.HashedPassword)
	} else {
		err = rec.db.QueryRowContext(ctx, "SELECT staff_id, password_hash FROM staff WHERE role='receptionist' AND email=$1", loginReq.Email).Scan(&loginResponse.UserID, &loginResponse.HashedPassword)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return loginResponse, nil // return empty model
		}
		return loginResponse, err // return empty model
	}

	return loginResponse, err
}
