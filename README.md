## open-feature-operator

### Project structure

```
├── Dockerfile
├── Makefile
├── PROJECT
├── README.md
├── agent ( Contains the agent that is injected by Mutating admission webhook)
├── api ( Custom Resource definitions )
├── bin
├── config
├── controllers
├── examples ( Example of how to use the operator )
├── go.mod
├── go.sum
├── hack
├── main.go
└── webhooks ( Mutating admission webhooks to insert the sidecar )
```

### Architecture

As per the issue [here](https://github.com/open-feature/research/issues/1)
High level architecture is as follows:

<img src="images/arch-0.png" width="560">

### Workflow

When wishing to leverage featureflagging within the local pod, the following steps are required:

1. Create a new feature flag custom resource e.g.
```
apiVersion: core.openfeature.dev/v1alpha1
kind: FeatureFlagConfiguration
metadata:
  name: featureflagconfiguration-sample
spec:
  featureFlagSpec: |
    {
      "foo" : "bar"
    } 
```

2. Reference the CR within the pod spec annotations
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  annotations:
    openfeature.dev: "enabled"
    openfeature.dev/featureflagconfiguration: "featureflagconfiguration-sample"
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

3. Example usage from host container

```
root@nginx:/# curl localhost:8080
{
  "foo" : "bar"
} 
```

#### Build deploy

```
docker buildx build --platform="linux/amd64,linux/arm64" -t tibbar/of-operator:v1.2 . --push
IMG=tibbar/of-operator:v1.2 make generate
IMG=tibbar/of-operator:v1.2 make deploy
 ```

 ### TODO

    * [ ] Add a test for the operator
    * [ ] Detect and update configmaps on change
    * [ ] Finalizers
    * [ ] Cleanup on deletion