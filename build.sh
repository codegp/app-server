
set -e

cd $GOPATH/src/github.com/codegp/app-server/server
CGO_ENABLED=0 go build
cd $GOPATH/src/github.com/codegp/app-server
docker build -t local/codegp/app-server .
kubectl delete deployment app-server-deployment  && kubectl create -f app-server.yaml
