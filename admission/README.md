# Admission

This service is our implementation of a Kubernetes admission controller.

Its purpose is to verify that before a `ConferenceTalk` is created, a speaker with the ID referenced in the ConferenceTalk object already exists.

The service is configured to use HTTPS with self-signed certificates, because Kubernetes uses HTTPS for admission requests.

## Project structure

The admission controller is a web service that handles HTTPS requests.
It it called by the Kubernetes API server.

Currently, it runs the same handler for any route.

Kubernetes is configured to call it via `ValidatingWebhookConfiguration`. See [admissionreview.yaml](k8s/admissionreview.yaml).

## How to build

The image is already build and can be found at <https://hub.docker.com/repository/docker/asankovcb/kubernetes-extensibility-admission>.

If you want to modify the code and build your own image you follow these instructions:

- Modify the code
- To build the container image, go to the root folder of the project.
- Run:

```shell
docker build -t <repo>/kubernetes-extensibility-admission . -f admission/Dockerfile
docker push <repo>/kubernetes-extensibility-admission
```

- Modify the [Kubernetes deployment spec](k8s/deployment.yaml) to use your image.
