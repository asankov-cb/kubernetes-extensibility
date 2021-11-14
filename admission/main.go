package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/asankov-cb/kubernetes-extensibility/admission/handlers"
	"github.com/asankov-cb/kubernetes-extensibility/admission/validators/conferencetalks"
	istaconv1 "github.com/asankov-cb/kubernetes-extensibility/api/v1"
	"github.com/sirupsen/logrus"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`

	port = 8080
)

// paths to the cert and key files
var (
	certPath = filepath.Join(tlsDir, tlsCertFile)
	keyPath  = filepath.Join(tlsDir, tlsKeyFile)
)

var (
	scheme = runtime.NewScheme()
)

// init registers the API resources to the scheme
// used to initialize the Kubernetes client.
//
// This is needed in order for Kubernetes to be able to map
// API responses to their Go types.
func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(istaconv1.AddToScheme(scheme))
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatalf("Error while initializing in-cluster Kubernetes config - service probably not running in Kubernetes: %v", err)
	}

	k8sClient, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		logrus.Fatalf("Error while initializing Kubernetes client - service probably not running in Kubernetes: %v", err)
	}

	validator := conferencetalks.NewValidator(k8sClient)
	handler := handlers.NewHandler(validator)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	logrus.Infoln("Admission controller listening on port [%d]", port)
	if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
		logrus.Fatal(err)
	}
}
