# Reads the Specifications of one Deployment and one Service from a YAML File and Act Accordingly

This example program demonstrates a classical use case by doing the following:

- read a YAML file with one deployment and one related service

- split it in k8s resources

- decode deployment from YAML to JSON

- unmarshall deployment from JSON to k8s resource

- try to create the deployment, and update it, if creation failed

- decode service from YAML to JSON

- unmarshall service from JSON to k8s resource

- try to create the service, and update it, if creation failed

You can adopt the source code from this example to write programs that manage
other types of resources through the Kubernetes API.

## Running this example

Make sure you have a Kubernetes cluster and `kubectl` is configured:

    kubectl get nodes

Compile this example on your workstation:

```
cd create-update-delete-deployment
go build -o ./app
```

Now, run this application on your workstation with your local kubeconfig file:

```
./app

Usage of ./app:
  -f string
    	absolute path to the YAML configuration file (default "configuration.yaml")
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file (default "$HOME/.kube/config")
```

Running this command will execute the following operations on your cluster:

1. **Create Deployment:** This will create a 2 replica Deployment. Verify with
   `kubectl get pods`.
2. **Update Deployment:** This will update the Deployment resource created in
   previous step by setting the replica count to 1 and changing the container
   image to `nginx:1.13`. You are encouraged to inspect the retry loop that
   handles conflicts. Verify the new replica count and container image with
   `kubectl describe deployment demo`.
3. **Rollback Deployment:** This will rollback the Deployment to the last
   revision. In this case, it's the revision that was created in Step 1.
   Use `kubectl describe` to verify the container image is now `nginx:1.12`.
   Also note that the Deployment's replica count is still 1; this is because a
   Deployment revision is created if and only if the Deployment's pod template
   (`.spec.template`) is changed.
4. **List Deployments:** This will retrieve Deployments in the `default`
   namespace and print their names and replica counts.
5. **Delete Deployment:** This will delete the Deployment object and its
   dependent ReplicaSet resource. Verify with `kubectl get deployments`.

Each step is separated by an interactive prompt. You must hit the
<kbd>Return</kbd> key to proceeed to the next step. You can use these prompts as
a break to take time to  run `kubectl` and inspect the result of the operations
executed.

You should see an output like the following:

```
Creating deployment...
Created deployment "demo-deployment".
-> Press Return key to continue.

Updating deployment...
Updated deployment...
-> Press Return key to continue.

Rolling back deployment...
Rolled back deployment...
-> Press Return key to continue.

Listing deployments in namespace "default":
 * demo-deployment (1 replicas)
-> Press Return key to continue.

Deleting deployment...
Deleted deployment.
```

## Cleanup

You can clean up the created deployment and service with:

    kubectl delete -f ./configuration.yaml

Accordingly, using other used YAML configuration and kubeconfig file, respectively.





