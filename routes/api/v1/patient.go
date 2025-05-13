package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/medi-go/models"
	"github.com/harshitrajsinha/medi-go/store"
)

type PatientRoutes struct {
	service *store.Store
}

func NewPatientRoutes(service *store.Store) *PatientRoutes {
	return &PatientRoutes{
		service: service,
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (p *PatientRoutes) GetAllPatients(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get data from store
	resp, err := p.service.GetAllPatients(ctx)
	if err != nil {
		// send error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Data: resp})
	log.Println("All patients data populated successfully")
}

func (p *PatientRoutes) CreatePatient(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// send error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}
	defer r.Body.Close()

	// Decode request body
	var patientReq models.Patient
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&patientReq)
	if err != nil {
		// send error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		log.Println(err)
		return
	}

	// validate request body
	if err := models.ValidatePatientReq(patientReq); err != nil {
		// send bad request response
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to store
	patientToken, err := p.service.CreatePatient(ctx, &patientReq)
	if err != nil {
		// error while storing data to db
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// send response
	if patientToken != -1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Code: http.StatusCreated, Message: "patient data inserted into DB successfully!", Data: patientToken})
		log.Println("patient data inserted into DB successfully!")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "No rows inserted - Possibly data already exists"})
		log.Println("No rows inserted - Possibly data already exists")
	}

}

func (p *PatientRoutes) GetPatientByTokenID(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()
	params := mux.Vars(r)

	// Get id
	id := params["token_id"]

	if len(id) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid token ID"})
		log.Println("Invalid token ID")
		return
	}

	// Get data from service layer
	resp, err := p.service.GetPatientByTokenID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	// Send response
	var respData []interface{}
	respData = append(respData, resp) // enclose data in an array
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Data: respData})
	log.Println("Patient data populated successfully based on token ID")
}

func (p *PatientRoutes) UpdatePatient(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get id
	params := mux.Vars(r)
	id := params["token_id"]

	if len(id) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid token ID"})
		log.Println("Invalid token ID")
		return
	}
	defer r.Body.Close()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	var patientReq models.Patient
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&patientReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		log.Println(err)
		return
	}

	// validate request body
	if err := models.ValidatePatientReq(patientReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to store to update engine
	updatedPatient, err := p.service.UpdatePatient(ctx, id, &patientReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedPatient > 0 {
		// data is updated successfully
		log.Println("Patient data updated successfully!")
		// Get the updated result
		p.GetPatientByTokenID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "No data present for provided token ID or data already exists"})
		log.Println("value of updatedPatient is ", updatedPatient)
		return
	}
}

func (p *PatientRoutes) UpdatePatientPartial(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()

	// Get id
	params := mux.Vars(r)
	id := params["token_id"]

	if len(id) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid token ID"})
		log.Println("Invalid token ID")
		return
	}
	defer r.Body.Close()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	var patientReq models.Patient
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&patientReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		log.Println(err)
		return
	}

	// validate request body  for partial update
	if err := models.ValidatePatientPatchReq(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: err.Error()})
		log.Println(err)
		return
	}

	// Pass data to store to update engine
	updatedPatient, err := p.service.UpdatePatient(ctx, id, &patientReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while reading data"})
		panic(err)
	}

	if updatedPatient > 0 {
		// data is updated successfully
		log.Println("Patient data updated successfully!")
		// Get the updated result
		p.GetPatientByTokenID(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "No data present for provided token ID or data already exists"})
		log.Println("value of updatedPatient is ", updatedPatient)
		return
	}
}

func (p *PatientRoutes) DeletePatient(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	ctx := r.Context()
	params := mux.Vars(r)

	// Get id
	id := params["token_id"]

	if len(id) <= 0 {
		log.Println(id)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid token ID"})
		log.Println("Invalid token ID")
		return
	}

	// Pass data to service layer to delete engine
	deletedPatient, err := p.service.DeletePatient(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured while deleting data"})
		panic(err)
	}

	if deletedPatient > 0 {
		// data is deleted successfully
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		log.Println("value of deletedPatient is ", deletedPatient)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "No data present for provided patient ID or data already deleted"})
		log.Println("value of deletedPatient is ", deletedPatient)
		return
	}

}
