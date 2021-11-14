package conferencetalks

import (
	"context"
	"fmt"

	istaconv1 "github.com/asankov-cb/kubernetes-extensibility/api/v1"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/client-go/kubernetes/scheme"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Validator validates an incoming AdmissionRequest for istacon.org/v1/ConferenceTalk.
type Validator struct {
	k8sClient client.Client
	decoder   runtime.Decoder
}

// NewValidator returns new Validator with the given Kubernetes Client.
func NewValidator(client client.Client) *Validator {
	return &Validator{
		k8sClient: client,
		decoder:   scheme.Codecs.UniversalDeserializer(),
	}
}

// Process validates the incoming AdmissionRequest.
//
// It expectes the request to be for istacon.org/v1/ConferenceTalk.
// If it is not, it will skip the check and return AdmissionResponse with Allowed=true.
//
// For ConferenceTalk object, it will make an API call to the istacon.org/v1/Speaker API
// to get the speaker with the speakerID in the ConferenceTalk object.
// If that speaker does not exist, it will return that the conference talk should not be allowed.
func (v *Validator) Process(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	response := &admissionv1.AdmissionResponse{UID: request.UID, Allowed: true}

	if !shoudValidateRequest(request) {
		return response, nil
	}

	confTalk := istaconv1.ConferenceTalk{}
	if _, _, err := v.decoder.Decode(request.Object.Raw, nil, &confTalk); err != nil {
		return nil, err
	}

	s := istaconv1.Speaker{}
	err := v.k8sClient.Get(context.Background(), client.ObjectKey{
		Namespace: confTalk.Namespace,
		Name:      confTalk.Spec.SpeakerID,
	}, &s)
	if err != nil {
		if apierrors.IsNotFound(err) {
			response.Allowed = false
			response.Result = &metav1.Status{
				Message: fmt.Sprintf("Speaker with ID [%s] does not exist.", confTalk.Spec.SpeakerID),
			}
			return response, nil
		}
		return nil, err
	}

	return response, nil
}

// shoudValidateRequest returns true if the request if of API group "istacon.org" and kind "ConferenceTalk".
//
// Any request that is not of that kind should not even get here, because of the ValidatingWebhookConfiguration,
// however this is an assumptions and this safeguard is here just in case.
func shoudValidateRequest(request *admissionv1.AdmissionRequest) bool {
	return request.Kind.Group == "istacon.org" && request.Kind.Kind == "ConferenceTalk"
}
