# Set Up Redis Master-Slave in K8S

Details is https://kubernetes.io/docs/tutorials/stateless-application/guestbook/

##

```sh
kubectl apply -f k8s/redis/redis-master-deployment.yaml --namespace the-automatic-manager
kubectl apply -f k8s/redis/redis-master-service.yaml --namespace the-automatic-manager
kubectl apply -f k8s/redis/redis-slave-deployment.yaml --namespace the-automatic-manager
kubectl apply -f k8s/redis/redis-slave-service.yaml --namespace the-automatic-manager

```

