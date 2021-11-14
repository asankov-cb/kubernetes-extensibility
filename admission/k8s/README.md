# Kubernetes specs

This folder contains the Kubernetes spec files to run the admission controller.

- [admissionreview.yaml](admissionreview.yaml) - contains the `ValidatingWebhookConfiguration` that tells Kubernetes to call this service for admission review. It contains all the details about how to call the service - Kubernetes service, port, route (not configured, default to `/`), CA certificates, etc.
- [auth.yaml](auth.yaml) - the `ClusterRole` and `ClusterRoleBinding` needed for the service to make an API call to the Speakers API
- [deployment.yaml](deployment.yaml) - the Kubernetes `Deployment` to run the web service
- [service.yaml](service.yaml) - the Kubernetes `Service` object to expose the web service
- [secret.yaml](secret.yaml) - the Kubernetes `Secret` object that contains the CA data
