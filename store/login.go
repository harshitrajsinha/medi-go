package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/harshitrajsinha/medi-go/models"
)

func (rec *Store) GetLoginInfo(ctx context.Context, loginReq *models.Credentials) (string, error) {

	var hashedPassword string
	var err error

	if loginReq.Role == "doctor" {
		err = rec.db.QueryRowContext(ctx, "SELECT password_hash FROM doctor WHERE role='doctor' AND email=$1", loginReq.Email).Scan(&hashedPassword)
	} else {
		err = rec.db.QueryRowContext(ctx, "SELECT password_hash FROM staff WHERE role='receptionist' AND email=$1", loginReq.Email).Scan(&hashedPassword)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return hashedPassword, nil // return empty model
		}
		return hashedPassword, err // return empty model
	}

	return hashedPassword, err
}
