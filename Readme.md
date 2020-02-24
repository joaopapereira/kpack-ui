[![Actions Status](https://github.com/joaopapereira/kpack-ui/workflows/CI/badge.svg)](https://github.com/joaopapereira/kpack-ui/actions)

# kpack UI

This application allow the users to, in a visual way, see the
images that are being managed by kpack as well as a separation
between namespaces.

It should be installed in the same cluster as kpack and reads
the CRD directly.

## Installation (Web Version)

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

## For development (Web Version)

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
LOCAL_START=1 go run cmd/kpack-ui-server/main.go
```

## For development (CLI Version)

Need to install Go and minikube

### If you need to use minikube and local registry follow these instructions:
Start a minikube cluster on Mac
```bash
    minikube start --vm-driver=hyperkit --bootstrapper=kubeadm \
    --insecure-registry "registry.default.svc.cluster.local:5001" 
```
Update `/etc/hosts` by adding the name registry.default.svc.cluster.local on the same line as the entry for localhost. It should look something like this:
```bash
##
127.0.0.1       localhost registry.default.svc.cluster.local
255.255.255.255 broadcasthost
::1             localhost
```

Update the minikube `/etc/hosts` with the host ip for registry.default.svc.cluster.local
 ```bash
    minikube ssh \
    "echo \"192.168.64.1       registry.default.svc.cluster.local\" \
    | sudo tee -a  /etc/hosts"
```

Start the local registry:
```bash
docker-compose up -d
```

Run the app:
```bash
go run cmd/kpack-gui/main.go
```
