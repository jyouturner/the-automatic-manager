


### Set up ECR

create a ECR repository "the-automatic-manager"

### Build Container Image and Push to ECR

````
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ....dkr.ecr.us-west-2.amazonaws.com
docker build . -t the-automatic-manager       
docker tag the-automatic-manager:latest ....dkr.ecr.us-west-2.amazonaws.com/the-automatic-manager:0.0.2
docker push ....dkr.ecr.us-west-2.amazonaws.com/the-automatic-manager:0.0.2
````

### Minicube

````
minikube start
minikube addons configure registry-creds
minikube addons enable registry-creds
````

now set the context to minikube

````
kubectl config use-context minikube
````

create the namespace

````
kubectl create -f namespace.json 
````

deploy the config

````
kubectl apply -f calendar-monitor-config-dev.yml --namespace the-automatic-manager



