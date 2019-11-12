[![Actions Status](https://github.com/joaopapereira/kpack-ui/workflows/CI/badge.svg)](https://github.com/joaopapereira/kpack-ui/actions)

# kpack UI

This application allow the users to, in a visual way, see the
images that are being managed by kpack as well as a separation
between namespaces.

It should be installed in the same cluster as kpack and reads
the CRD directly.

## Installation

Update the file in `deployment/kubernetes/deployment.yaml` replacing
the `CHANGEME` tags with the values needed

Access to the cluster using `kubectl` or `kapp`

Run the following command to install it:
```bash
kubectl apply -f deployment/kubernetes/deployment.yaml
```

**Or using kapp**

```bash
kapp deploy -a kpack-ui -f deployment/kubernetes/deployment.yaml
```

## For development

Need to install Go and nodejs

Build the frontend:
```bash
pushd ui
    npm install
    npm run build
popd
```

The prior step is mandatory to ensure that the go application
can display the built frontend
Run the app:
```bash
LOCAL_START=1 go run cmd/kpack-ui/main.go
```
