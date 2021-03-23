kind create cluster --name rabbit --image kindest/node:v1.18.4
kubectl create ns rabbits
kubectl get storageclass


---

kubectl apply -n rabbits -f .\kubernetes\rabbit-rbac.yaml
kubectl apply -n rabbits -f .\kubernetes\rabbit-configmap.yaml
kubectl apply -n rabbits -f .\kubernetes\rabbit-secret.yaml
kubectl apply -n rabbits -f .\kubernetes\rabbit-statefulset.yaml

---
    
kubectl -n rabbits port-forward rabbitmq-0 8080:15672

---

cd messaging\rabbitmq\applications\publisher
docker build . -t aimvector/rabbitmq-publisher:v1.0.0

kubectl apply -n rabbits -f deployment.yaml

kubectl -n rabbits port-forward <name> 8081:15672
kubectl -n rabbits port-forward  8081:80

---



kubectl logs -f simple-deployment-4098151155-n8bqr