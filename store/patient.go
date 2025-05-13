package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/harshitrajsinha/medi-go/models"
)

type patientQueryResponse struct {
	Fullname   string `json:"fullname"`
	Gender     string `json:"gender"`
	Age        int    `json:"age"`
	Contact    string `json:"contact"`
	Symptoms   string `json:"symptoms"`
	AssignedTo string `json:"assigned_to"`
	TokenID    string `json:"token_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// apply pagination
func (rec *Store) GetAllPatients(ctx context.Context, limit int32, offset int32) (interface{}, error) {

	var total_records int32

	if limit <= 0 {
		limit = 5
	}

	rows, err := rec.db.QueryContext(ctx, "SELECT fullname, gender, age, contact, symptoms, assigned_to, token_id, updated_at, created_at, count(*) over() as total_records FROM patient LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return patientQueryResponse{}, nil // return empty model
		}
		return patientQueryResponse{}, err // return empty model
	}
	defer rows.Close()

	// slice to store all rows
	allPatientData := make([]interface{}, 0)

	// Get each row data into a slice
	for rows.Next() {
		var queryData patientQueryResponse
		var assignedDoctor string
		// var registeredBy string
		// Return single row
		err = rows.Scan(
			&queryData.Fullname, &queryData.Gender, &queryData.Age, &queryData.Contact, &queryData.Symptoms, &assignedDoctor, &queryData.TokenID, &queryData.UpdatedAt, &queryData.CreatedAt, &total_records)
		if err != nil {
			return patientQueryResponse{}, err
		}

		// Get Assigned doctor name
		doctorData, err := rec.GetDoctorById(ctx, assignedDoctor)
		if err != nil {
			queryData.AssignedTo = ""
		} else {
			doctor_name := doctorData.(doctorQueryResponse).Fullname
			queryData.AssignedTo = doctor_name
		}

		// store each row
		allPatientData = append(allPatientData, queryData)
	}

	return allPatientData, nil
}

func (rec *Store) GetPatientByTokenID(ctx context.Context, token_id string) (interface{}, error) {
	var queryData patientQueryResponse
	var assignedDoctor string
	// var registeredBy string

	err := rec.db.QueryRowContext(ctx, "SELECT fullname, gender, age, contact, symptoms, assigned_to, token_id, updated_at, created_at FROM patient WHERE token_id=$1", token_id).Scan(
		&queryData.Fullname, &queryData.Gender, &queryData.Age, &queryData.Contact, &queryData.Symptoms, &assignedDoctor, &queryData.TokenID, &queryData.UpdatedAt, &queryData.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queryData, nil // return empty model
		}
		return queryData, err // return empty model
	}

	// Get Assigned doctor name
	doctorData, err := rec.GetDoctorById(ctx, assignedDoctor)
	if err != nil {
		queryData.AssignedTo = ""
	} else {
		doctor_name := doctorData.(doctorQueryResponse).Fullname
		queryData.AssignedTo = doctor_name
	}

	return queryData, err
}

func (rec *Store) CreatePatient(ctx context.Context, patientMod *models.Patient) (int64, error) {

	var tokenID int64

	// Begin DB transaction
	tx, err := rec.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error while inserting data ", err)
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()

	var query string = "INSERT INTO patient (fullname, gender, age, contact, symptoms, assigned_to, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING token_id"
	err = tx.QueryRowContext(ctx, query, patientMod.Fullname, patientMod.Gender, patientMod.Age, patientMod.Contact, patientMod.Symptoms, patientMod.Assigned_to, patientMod.Created_by).Scan(&tokenID)

	if err != nil {
		log.Println("Error while inserting data ", err)
		return -1, err
	}

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	log.Println("Error while inserting data ", err)
	// 	return -1, err
	// }

	return tokenID, nil
}

func (rec *Store) UpdatePatient(ctx context.Context, tokenID string, patientReq *models.Patient) (int64, error) {

	// DB transaction
	tx, err := rec.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error while updating data ", err)
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	var query strings.Builder
	var args []interface{}
	argCount := 1

	query.WriteString("UPDATE patient SET ")

	// If fullname in request
	if patientReq.Fullname != "" {
		query.WriteString(fmt.Sprintf("fullname=$%d ", argCount))
		args = append(args, patientReq.Fullname)
		argCount++
	}
	if patientReq.Gender != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("gender=$%d ", argCount))
		args = append(args, patientReq.Gender)
		argCount++
	}
	if patientReq.Age > 1 {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("age=$%d ", argCount))
		args = append(args, patientReq.Age)
		argCount++
	}

	if patientReq.Contact != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("contact=$%d ", argCount))
		args = append(args, patientReq.Contact)
		argCount++
	}

	if patientReq.Symptoms != "" {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("symptoms=$%d ", argCount))
		args = append(args, patientReq.Symptoms)
		argCount++
	}

	if patientReq.Assigned_to != uuid.Nil {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("assigned_to=$%d ", argCount))
		args = append(args, patientReq.Assigned_to)
		argCount++
	}

	if patientReq.Created_by != uuid.Nil {
		if argCount > 1 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("created_by=$%d ", argCount))
		args = append(args, patientReq.Created_by)
		argCount++
	}

	query.WriteString(fmt.Sprintf("WHERE token_id=$%d ", argCount))
	args = append(args, tokenID)

	result, err := tx.ExecContext(ctx, query.String(), args...)
	if err != nil {
		log.Println("Error while updating data ", err)
		return -1, err
	}

	rowAffected, err := result.RowsAffected()
	return rowAffected, nil
}

func (rec *Store) DeletePatient(ctx context.Context, tokenID string) (int64, error) {

	// DB transaction
	tx, err := rec.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("Transaction rollback error: ", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Println("Commit rollback error: ", cmErr)
			}
		}
	}()

	var query string = "DELETE FROM patient WHERE token_id=$1"
	result, err := tx.ExecContext(ctx, query, tokenID)
	if err != nil {
		return -1, err
	}
	rowAffected, err := result.RowsAffected()

	return rowAffected, nil

}
