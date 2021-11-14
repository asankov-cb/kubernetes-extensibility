package handlers

import (
	"encoding/json"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"

	"github.com/asankov-cb/kubernetes-extensibility/admission/validators/conferencetalks"
	"github.com/sirupsen/logrus"
)

// Handler is the HTTP handler that handles the admission request.
type Handler struct {
	validator *conferencetalks.Validator
}

// NewHandler returns new Handler with the given conferencetalks.Validator.
func NewHandler(validator *conferencetalks.Validator) *Handler {
	return &Handler{validator: validator}
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logrus.Infoln("Handling validation request")

	admissionRequest := admissionv1.AdmissionReview{}
	if err := json.NewDecoder(r.Body).Decode(&admissionRequest); err != nil {
		logrus.Error("Error while decoding request body: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	admissionResponse, err := h.validator.Process(admissionRequest.Request)
	if err != nil {
		logrus.Error("Error while processing request: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(rw).Encode(admissionv1.AdmissionReview{
		// K8s expects the incoming response to have TypeMeta with API version
		// the same one as the one in the request.
		TypeMeta: admissionRequest.TypeMeta,
		Response: admissionResponse,
	}); err != nil {
		logrus.Error("Error while encoding response: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Infoln("Successfully handled validation request")
}
