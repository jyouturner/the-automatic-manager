.PHONY: build clean deploy run_monitor

build:
	export GO111MODULE=on
	echo "the lambda executable to be deployed by serverless framework to AWS"
	env GOOS=linux go build -o bin/lambda/authorization example/awsserverless/lambda/authorization/main.go
	env GOOS=linux go build -o bin/lambda/authorized example/awsserverless/lambda/authorized/main.go


clean:
	rm -rf ./bin 

deploy: clean build
	serverless deploy --config serverless.yml --stage dev --region us-west-2 --verbose

run_monitor:
	go run cmd/monitorcalendar/main.go
