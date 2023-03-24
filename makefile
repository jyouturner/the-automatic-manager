.PHONY: build clean deploy run_monitor

build:
	export GO111MODULE=on
	# the lambda executable to be deployed by serverless framework to AWS
	env GOOS=linux go build -o bin/lambda/authorization lambda/authorization/main.go
	env GOOS=linux go build -o bin/lambda/authorized lambda/authorized/main.go
	env GOOS=linux go build -o bin/lambda/calendar_monitor lambda/calendarmonitor/main.go
	env GOOS=linux go build -o bin/lambda/notion_monitor lambda/notionmonitor/main.go
	env GOOS=linux go build  -o bin/lambda/jira_monitor lambda/jiramonitor/main.go
	#echo "the command line executable to run the program on your intel mac"
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/monitor_calendar_darwin_amd64 cmd/calendarmonitor/main.go
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/monitor_notion_darwin_amd64 cmd/notionmonitor/main.go
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/monitor_jira_darwin_amd64 cmd/jiramonitor/main.go
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/auth0_web_darwin_amd64 cmd/auth0web/main.go
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/auto_darwin_amd64 cmd/automode/main.go
	#echo "the command line executable to run the program on linux or k8s container"
	#env GOOS=linux go build -v -o bin/monitor_calendar_linux local/monitorcalendar/main.go
	#env GOOS=linux go build -v -o bin/monitor_notion_linux local/monitornotion/main.go
	#env GOOS=linux go build -v -o bin/monitor_jira_linux local/monitorjira/main.go
	#echo "the simple web server to test stuff on k8s"
	#env GOOS=linux go build -o bin/web_server_linux web/main.go
	

clean:
	rm -rf ./bin ./deploy

deploy: clean build
	serverless deploy --config serverless.yml --stage dev --region us-west-2 --verbose

