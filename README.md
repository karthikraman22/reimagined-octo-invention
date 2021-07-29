# reimagined-octo-invention

dev setup

postgres

install helm
add helm repo

helm install 

postgres

helm install postgres --set postgresqlUsername=postgres --set postgresqlPassword=postgres bitnami/postgresql

setup port forwarding 

kubectl port-forward --namespace default svc/postgres-postgresql 5432:5432



nats

helm install nats --set auth.enabled=true,auth.user=nats,auth.password=nats bitnami/nats

setup port forwarding

kubectl port-forward --namespace default svc/nats-client 4222:4222 &

To create configmap for ledger

kubectl apply -f ./k8s/dev-configmap.yaml

deploy the app
kubectl apply -f ./k8s/k8s-deployment.yaml

forward the app
kubectl port-forward --namespace default svc/ledger-service 10001


to check
kubectl proxy

