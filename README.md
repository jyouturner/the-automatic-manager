# The Automatic Manager - Improve Your Workspace Efficiency through Connection and Automation


The "connecting tissue" to help improve engineering manager's productivity by integrating:
* Notion
* Google Calendar
* Jira
* Confluence
* Slack
* Github

## Auth

1. user Auth0 for user authentication
2. Oauth2 for Google, Atlanssian, Github access
3. Personal access token for Jira Software, and Notion access

## Background

I was trying to figure out tooling to help improve productivity through managing the "to do" tasks and note taking.

## Feature Set

1. every day, it will create a journal task like "2022-01-15" in your Notion to-do list. (And move it to done the next morning). This is where you can put in any thoughts that happen (randomly) on daily basis.
2. it will check your Google calendar, find the next 10 events, and add to your to-do list. (it will ignore those "group meetings" like office hours, all-hands because most people don't take notes in those meetings)
3. after the calendar event/meeting is over, it will move those "task" from to-do to "done" state.
4. it scans the to-do tasks. For example, you are having a Zoom meeting with Joe, and you have your notion task "meeting with Joe" open. Say, you and Joe agree that he will create a release plan later. You can press "/" and select "to-do" to add a "to-do" item "follow up with Joe of the release" in the Notion task page. The system will recognize this item and create a Task "follow up with Joe of the release" in your To-Do list. This way that task stays with you even though the meeting is "done".
5. It can check your JIRA issues, if any issue stays in "In Dev" status for more than 5 business days, it can flag them and create To Do task for you to follow up.


## Installation

There are 3 ways that you can use this service. 

### Run Them At Your Local
In this case, all the configuration items are stored in local file (.env)

1. set up the access credentials in Google, Atlanssian and Notion
2. create the .env file
3. run the programs

### Deploy Them to AWS Lambda
In this case, the configuration items are handled differently, 
* access credentials in the AWS Secret Manager
* others in the local .env file
* oauth tokens are saved in AWS S3

Follow below steps to deploy:
1. set up the access credentials in Google, Atlanssian and Notion
2. create .env file
3. go to AWS secret manager, and create the secret with the credentials from step 1
4. run "make deploy"

### Deploy Them to K8S
In this case, the configuration items are handled differently,
* access credentials in the k8s secret
* other configurations are stored in k8s config map
* oauth2 tokens are stored in Redis (also deployed to k8s)

Follow below steps to deploy:
1. set up the access credentials in Google, Atlanssian and Notion
2. update the k8s/config yml files
3. run the kubectl commands

## Start

````
go run cmd/auth0web/main.go
````

sign in http://localhost:3000/user
google auth http://localhost:3000/tam/auth/google
atlanssian auth http://localhost:3000/tam/auth/atlanssian

## Choice of Programming Language and Frameworks

At the same time I am interested in Go. It is like going back to the 2000 with Java (before the J2EE thing). I kind of like it.

## Project Structure

Try to stay with the Go way.

## Architecture

## Infrastructure

## Authentication


## General Configuration (non-secrets)

Most of us have been very used to getting configuration items from environment variables. It is convenient and supported widely in both programming language and infrastructure. As far as configuration is concerned, here are the goals:

1. it has to make local testing easy
2. it has to make deployment easy
3. it has to make security team happy (usually not easy)

### Non-Secrets Configuration
I decided to use DotEnv (.env) which is likely the most popular way to handle enviornment variables(?). Go has GoDotEnv module too. So, in our project, we have dot env files like below

````
.env
.env.dev
.env.qa
.env.prod
````


### Secrets

## Manually deploy to AWS Lambda

````
serverless deploy --stage dev --region us-west-2 --verbose --aws-profile iamadmin-coopers-cose
````

## Issues
[ISSUES.md]

## Setup

### .env

create a file .env at the project root foler.

### AWS Set up

create a S3 bucket, and set in the ATLASSIAN_TOKEN_S3_BUCKET field of the .env file


### Atlanssian Setup

update the .env, and set the 2 parameters

````
ATLASSIAN_TOKEN_S3_BUCKET=[the bucket name]
ATLASSIAN_TOKEN_S3_FILE=storage/token/atlanssian/token.json
````

create the Atlanssian App

First go to https://developer.atlanssian.com/console/myapps/, create a OAuth 2.0 integration.
Choose "Authorization", select "Use rotating refresh tokens", enter the callback URL "https://localhost:8080"
Choose "permissions", 
click "Add" to "Confluence API", then click "Config", and enable "read confluence"
click "Add" to "Jira platform REST API", and check the 3 scopes: read:jira-user,read:jira-work,write:jira-work

### Initiate the first Atlanssian Oauth2 token

run commmand
````
go run cmd/setup/atlanssian/main.go
````
it will print a URL on the terminal, copy it to a browser, the next page will have url like https://localhost:8080/?code=[ACCESS CODE]&state=state-token
copy the access code to the terminal, and the Go program will fetch the access token and refresh token and upload to S3

## Notion Setup

### Create the Integration

### Find the To Do List database id
Open the To Do, at the top right corner, click "Share", the click the "Copy URL", the URL will be something like
https://www.notion.so/e58691b62b1c4d22a6b7aba20ba4cd03?v=8d9dbfa710d442c78f342d64e1a21274
the e58691b62b1c4d22a6b7aba20ba4cd03 is the databse id

edit the .env

````
NOTION_KEY=[secret]
NOTION_DATABASE_ID=[database id]
````
### Share the To Do List with the Integration
Open the To Do, at the top right corner, click "Share", and Invite the integration.


## lambda

````
serverless plugin install -n serverless-pseudo-parameters
````

## Make Sure

Google OAuth Config follow this format
````
{
	"web": {
		"client_id": "....apps.googleusercontent.com",
		"project_id": "...",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_secret": "...",
		"redirect_uri": "https://....execute-api.us-west-2.amazonaws.com/dev/authorized",
		"redirect_uris": ["https://....execute-api.us-west-2.amazonaws.com/dev/authorized"]
	}
}
````

Atlanssian
````
{
	"installed": {
		"client_id": "...",
		"auth_uri": "https://auth.atlassian.com/authorize?audience=api.atlassian.com",
		"token_uri": "https://auth.atlassian.com/oauth/token",
		"client_secret": "...",
		"redirect_uri": "https://....execute-api.us-west-2.amazonaws.com/dev/authorized"
	}
}
````

## References
https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
https://peter.bourgon.org/go-best-practices-2016/#configuration