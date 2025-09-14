package models

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Patient struct {
	Fullname    string    `json:"fullname"`
	Gender      string    `json:"gender"`
	Age         int       `json:"age"`
	Contact     string    `json:"contact"`
	Symptoms    string    `json:"symptoms"`
	Treatment   string    `json:"treatment"`
	Assigned_to uuid.UUID `json:"assigned_doctor"`
	Created_by  uuid.UUID `json:"registered_by"`
}

func validateName(fullname string) error {
	if fullname != "" {
		return nil
	}
	return errors.New("fullname must not be empty")

}

// ['male', 'female', 'other']
func validateGender(gender string) error {
	for _, value := range [3]string{"male", "female", "other"} {
		if gender == value {
			return nil
		}
	}
	return errors.New("gender must be one of following - ['male', 'female', 'other']")
}

func validateAge(age int) error {
	if age > 1 || age < 125 {
		return nil
	}
	return errors.New("invalid age value")
}

func validateContact(contact string) error {
	if len(contact) == 10 {
		return nil
	}
	return errors.New("contact no must be of 10 digits")
}

func validateAssignedDoctor(assignedDoctor uuid.UUID) error {
	if assignedDoctor != uuid.Nil {
		return nil
	}
	return errors.New("doctor needs to be assigned")
}

func validateRegisteredBy(registeredBy uuid.UUID) error {
	if registeredBy != uuid.Nil {
		return nil
	}
	return errors.New("receptionist needs to be assigned")
}

func ValidatePatientReq(patientRequest Patient) error {
	var err error

	if err = validateName(patientRequest.Fullname); err != nil {
		return err
	}

	if err = validateGender(patientRequest.Gender); err != nil {
		return err
	}

	if err = validateAge(patientRequest.Age); err != nil {
		return err
	}

	if err = validateContact(patientRequest.Contact); err != nil {
		return err
	}

	if err = validateAssignedDoctor(patientRequest.Assigned_to); err != nil {
		return err
	}

	if err = validateRegisteredBy(patientRequest.Created_by); err != nil {
		return err
	}

	return nil
}

// Function to check which key exists in request body
func verifyPatientRequestKeys(request []byte) [6]bool {

	var doesFullnameExists bool
	var doesGenderExists bool
	var doesAgeExists bool
	var doesContactExists bool
	var doesAssigned_toExists bool
	var doesCreated_byExists bool
	var data map[string]interface{}

	_ = json.Unmarshal([]byte(request), &data)

	if _, fullnameExists := data["fullname"]; !fullnameExists {
		doesFullnameExists = false
	} else {
		doesFullnameExists = true
	}
	if _, genderExists := data["gender"]; !genderExists {
		doesGenderExists = false
	} else {
		doesGenderExists = true
	}
	if _, contactExists := data["contact"]; !contactExists {
		doesContactExists = false
	} else {
		doesContactExists = true
	}
	if _, assigned_doctorExists := data["assigned_doctor"]; !assigned_doctorExists {
		doesAssigned_toExists = false
	} else {
		doesAssigned_toExists = true
	}
	if _, registered_byExists := data["registered_by"]; !registered_byExists {
		doesCreated_byExists = false
	} else {
		doesCreated_byExists = true
	}

	return [6]bool{doesFullnameExists, doesGenderExists, doesAgeExists, doesContactExists, doesAssigned_toExists, doesCreated_byExists}
}

func ValidatePatientPatchReq(request []byte) error {
	var err error
	var patientRequest Patient

	_ = json.Unmarshal(request, &patientRequest)

	// Check which key exists and verify accordingly
	doesKeyExists := verifyPatientRequestKeys(request)

	if doesKeyExists[0] { // if "fullname" exists in request body
		if err = validateName(patientRequest.Fullname); err != nil {
			return err
		}
	}
	if doesKeyExists[1] { // if "gender" exists in request body
		if err = validateGender(patientRequest.Gender); err != nil {
			return err
		}
	}
	if doesKeyExists[2] { // if "age" exists in request body
		if err = validateAge(patientRequest.Age); err != nil {
			return err
		}
	}
	if doesKeyExists[3] { // if "contact" exists in request body
		if err = validateContact(patientRequest.Contact); err != nil {
			return err
		}
	}
	if doesKeyExists[4] { // if "assigned_doctor" exists in request body
		if err = validateAssignedDoctor(patientRequest.Assigned_to); err != nil {
			return err
		}
	}
	if doesKeyExists[5] { // if "registered_by" exists in request body
		if err = validateRegisteredBy(patientRequest.Created_by); err != nil {
			return err
		}
	}

	return nil
}
