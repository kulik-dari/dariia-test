# dariia-test

This repository contains the codebase for the dariia-test project. Below is an overview of the project's purpose and instructions on how to get started.

## Overview

The dariia-test project aims to create a Kubernetes controller that manages `ClusterScan` resources, encapsulating arbitrary jobs and recording job execution results. The controller supports both one-off and recurring executions.

## Getting Started

### Prerequisites

- Go version v1.21.0+
- Docker version 17.03+
- Kubectl version v1.11.3+
- Access to a Kubernetes v1.11.3+ cluster

### Steps to Create and Deploy the Controller

#### Step 1: Review Kubebuilder Documentation

First, I reviewed the [Kubebuilder Book](https://book.kubebuilder.io/). This helped me understand the best practices and methodologies for developing Kubernetes controllers using Go and Kubebuilder. Even though there are other frameworks for different languages, I decided to go with Go and Kubebuilder for this project.

#### Step 2: Set Up the Repository

I created a new GitHub repository for the project. This provided a place to manage my code, track changes, and collaborate if needed.

#### Step 3: Implement the Controller

1. **Create a new Kubebuilder project**:
   - I initialized a new Kubebuilder project. This command set up the basic directory structure and configuration files for my controller.
   
2. **Create the API for the `ClusterScan` resource**:
   - Using Kubebuilder, I generated the API definition for the `ClusterScan` resource. This included creating the CRD (Custom Resource Definition) and its corresponding Go structs.

3. **Define the `ClusterScan` Spec and Status**:
   - I edited the generated `clusterscan_types.go` file to define the `spec` (desired state) and `status` (observed state) fields for the `ClusterScan` resource. This structured the data that my controller would manage.

4. **Implement the Controller Logic**:
   - In the `clusterscan_controller.go` file, I wrote the reconciliation logic. This included the steps for managing the lifecycle of `ClusterScan` resources and creating Kubernetes `Jobs` or `CronJobs` as needed.

5. **Configure ServiceAccount, Role, and RoleBinding**:
   - I defined a ServiceAccount to give my controller an identity within the cluster.
   - I created a Role with the necessary permissions to manage `ClusterScan` resources and create `Jobs` or `CronJobs`.
   - I linked the Role to the ServiceAccount using a RoleBinding, ensuring my controller had the appropriate permissions.

6. **Update the Deployment Configuration**:
   - I modified the deployment manifest to use the newly created ServiceAccount, ensuring the controller runs with the correct permissions.

7. **Build and Push the Docker Image**:
   - I built a Docker image for my controller, which included the compiled Go binary and its dependencies.
   - I pushed this image to a container registry, making it available for deployment in the cluster.

8. **Deploy the Controller**:
   - I applied the deployment manifest to create the controller's Deployment in the cluster. This started the controller, enabling it to begin reconciling `ClusterScan` resources.

#### Step 4: Verify the Setup

1. **Check the Pod Status**:
   - Using `kubectl`, I listed the pods in the namespace where the controller was deployed, ensuring the controller pod was running and not encountering errors.

2. **Check Logs for Errors**:
   - I viewed the logs of the controller pod to identify any errors or issues during startup and operation. This helped diagnose configuration problems or bugs in the code.

3. **Describe the Pod for More Information**:
   - I used `kubectl describe` to get detailed information about the controller pod. This included events, status conditions, and resource usage, which were helpful for troubleshooting.

#### Step 5: Debugging

If the container was crashing, I ensured the application inside the Docker container was configured correctly to connect to the Kubernetes API. This included setting the correct environment variables and ensuring the ServiceAccount had the necessary permissions. I reviewed the logs and deployment configuration to identify any misconfigurations or missing settings.

### Deployment Instructions
1. **Build and push your image to the specified location**:
   ```sh
   make docker-build docker-push IMG=<some-registry>/dariia-test:tag
2. ## Install the CRDs into the cluster

  ```sh
make install
```
3. **BDeploy the Manager to the cluster**:
```sh
make deploy IMG=<some-registry>/dariia-test:tag
```
Note: If you encounter RBAC errors, you may need to grant yourself cluster-admin privileges or be logged in as an admin.

4. **BCreate instances of your solution**:
```sh
kubectl apply -k config/samples/
```
Note: Ensure that the samples have default values to test it out.

### Uninstallation
1. **Delete the instances (CRs) from the cluster**:
```sh
kubectl delete -k config/samples/
```

2. **Delete the APIs (CRDs) from the cluster**:
```sh
make uninstall
```

3. **Undeploy the controller from the cluster**:
```sh
make undeploy
```

### Project Distribution

1. **Build the installer for the image built and published in the registry**:
```sh
make build-installer IMG=<some-registry>/dariia-test:tag
```
Note: The make build-installer target generates an install.yaml file in the dist directory. This file contains all the resources built with Kustomize, which are necessary to install this project without its dependencies.
2. **Using the installer**:
```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/dariia-test/<tag or branch>/dist/install.yaml
```
###  Additional Information
For more details, refer to the [Kubebuilder Book](https://book.kubebuilder.io/).
