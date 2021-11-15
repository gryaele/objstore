package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	objStore *memStore
}

func NewApi() *Api {
	return &Api{
		objStore: newMemStore(),
	}
}

func (a *Api) GetObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	objectID := vars["objectID"]
	logger := log.WithField("bucket", bucket).
		WithField("objectID", objectID).
		WithField("method", "GetObject")
	logger.Debugf("received get object request")
	obj, err := a.objStore.get(bucket, objectID)
	switch {
	case err == errNotFound:
		logger.Warnf("object not found")
		http.Error(w, "object not found", http.StatusNotFound)
		return
	case err != nil:
		logger.Error(err)
		http.Error(w, "failed to get object", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(obj); err != nil {
		logger.Errorf("failed to write response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("object successfully retrieved")
	w.WriteHeader(http.StatusOK)
}

func (a *Api) PutObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	objectID := vars["objectID"]
	logger := log.WithField("bucket", bucket).
		WithField("objectID", objectID).
		WithField("method", "PutObject")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read body: %v", err)
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	err = a.objStore.put(bucket, objectID, body)
	switch {
	case err == errAlreadyExists:
		logger.Warnf("object already exists: %v", err)
		http.Error(w, "object already exists", http.StatusPreconditionFailed)
		return
	case err != nil:
		logger.Error(err)
		http.Error(w, "failed to save object", http.StatusInternalServerError)
		return
	}

	rspBody, err := json.Marshal(&struct {
		Id string `json:"id"`
	}{
		Id: objectID,
	})
	if err != nil {
		logger.Errorf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(rspBody); err != nil {
		logger.Errorf("failed to write response: %v", err)
	}
	logger.Infof("Object stored")
}

func (a *Api) DeleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	objectID := vars["objectID"]
	logger := log.WithField("bucket", bucket).
		WithField("objectID", objectID).
		WithField("method", "DeleteObject")
	err := a.objStore.delete(bucket, objectID)
	switch {
	case err == errNotFound:
		logger.Warnf("object not found")
		http.Error(w, "object not found", http.StatusNotFound)
		return
	case err != nil:
		logger.Error(err)
		http.Error(w, "failed to delete object", http.StatusInternalServerError)
		return
	}
	logger.Infof("object deleted")
}

func (a *Api) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/objects/{bucket}/{objectID}", a.PutObject).Methods(http.MethodPut)
	r.HandleFunc("/objects/{bucket}/{objectID}", a.GetObject).Methods(http.MethodGet)
	r.HandleFunc("/objects/{bucket}/{objectID}", a.DeleteObject).Methods(http.MethodDelete)
	return r
}
