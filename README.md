# Todo App on Kubernetes with Redis and ScyllaDB

This is a Gin-based Todo application deployed on a local Kubernetes cluster using Kind. It uses Envoy Proxy for JWT-based authentication, ScyllaDB as the primary database, and Redis for caching. The application is containerized with Docker and managed with Helm in a monorepo setup.

---

## Prerequisites

Ensure the following tools are installed on your system:

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/)

---

## Setup Instructions

### 1. Create a Kind Cluster

```bash
kind create cluster --name todo-cluster
```

---

### 2. Install Dependencies (Redis and ScyllaDB) via Helm

#### Redis

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis
```

#### ScyllaDB

```bash
helm install scylladb bitnami/scylladb
```

Wait until all pods are in `Running` state.

---

### 3. Build and Load the Todo App Docker Image

```bash
docker build -t todo-app:latest .
kind load docker-image todo-app:latest
```

---

### 4. Deploy the Todo App with Helm

```bash
helm install todo-app ./charts/todo-app
```

Make sure `values.yaml` contains your environment variables like Redis and ScyllaDB connection strings.

---

### 5. Deploy Envoy Proxy

Create a ConfigMap for Envoy using your working `envoy.yaml`:

```bash
kubectl create configmap envoy-config --from-file=envoy.yaml --dry-run=client -o yaml | kubectl apply -f -
```

Apply the Envoy deployment and service:

```bash
kubectl apply -f envoy-deployment.yaml
```

---

### 6. Port Forward and Test

```bash
kubectl port-forward svc/envoy 8080:8080
```

Use Postman to test the application at:

- `http://127.0.0.1:8080/health` (JWT protected)
- Swagger docs available at the root or defined endpoint.

---

### 7. JWT Tokens

Generate JWTs using any third-party provider or service that matches the issuer and audience in your Envoy config.

---

## Notes

- This application is meant for local testing and development.
- All environment variables are managed through the `values.yaml` file.
