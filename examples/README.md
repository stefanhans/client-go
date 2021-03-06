# client-go Examples

This directory contains examples that cover various use cases and functionality
for client-go.

### Configuration

- [**Authenticate in cluster**](./in-cluster-client-configuration): Configure a
  client while running inside the Kubernetes cluster.
- [**Authenticate out of cluster**](./out-of-cluster-client-configuration):
  Configure a client to access a Kubernetes cluster from outside.

### Basics

- [**Managing resources with API**](./create-update-delete-deployment): Create,
  get, update, delete a Deployment resource.
- [**Reading specified resources from YAML file**](./read-yaml-configuration-deployment):
  Read the specifications of one Deployment and one Service resource from a YAML file and create and/or
  update accordingly.

### Advanced Concepts

- [**Work queues**](./workqueue): Create a hotloop-free controller with the
  rate-limited workqueue and the [informer framework][informer].
- [**Custom Resource Definition (successor of TPR)**](https://git.k8s.io/apiextensions-apiserver/examples/client-go):
  Register a custom resource type with the API, create/update/query this custom
  type, and write a controller that drives the cluster state based on the changes to
  the custom resources.

[informer]: https://godoc.org/k8s.io/client-go/tools/cache#NewInformer
