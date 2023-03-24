

```
go test -timeout 30s -v  github.com/jyouturner/automaticmanager/pkg/notion -run TestApiClient_ListToDoTasks
```

Must
* create below properties in the Notion Task List


Add A Page to Database

```
curl -X POST https://api.notion.com/v1/pages \
  -H "Authorization: Bearer secret_replace_with" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2021-08-16" \
  --data "{
    \"parent\": { \"database_id\": \"task_list_database_id\" },
    \"properties\": {
      \"Status\": {
        \"select\": {
          \"name\": \"To Do\"
        }
      },
      \"title\": {
        \"title\": [
          {
            \"text\": {
              \"content\": \"Yurts in Big Sur, California\"
            }
          }
        ]
      }
    }
  }"
```

response
```
{
  "object": "page",
  "id": "3f5a23be-b380-4295-a61a-ce230bb168b9",
  "created_time": "2022-01-10T20:34:00.000Z",
  "last_edited_time": "2022-01-10T20:34:00.000Z",
  "cover": null,
  "icon": null,
  "parent": {
    "type": "database_id",
    "database_id": "e58691b6-2b1c-4d22-a6b7-aba20ba4cd03"
  },
  "archived": false,
  "properties": {
    "Date Created": {
      "id": "'Y6%3C",
      "type": "created_time",
      "created_time": "2022-01-10T20:34:00.000Z"
    },
    "Status": {
      "id": "%5EOE%40",
      "type": "select",
      "select": null
    },
    "Name": {
      "id": "title",
      "type": "title",
      "title": [
        {
          "type": "text",
          "text": {
            "content": "Yurts in Big Sur, California",
            "link": null
          },
          "annotations": {
            "bold": false,
            "italic": false,
            "strikethrough": false,
            "underline": false,
            "code": false,
            "color": "default"
          },
          "plain_text": "Yurts in Big Sur, California",
          "href": null
        }
      ]
    }
  },
  "url": "https://www.notion.so/Yurts-in-Big-Sur-California-3f5a23beb3804295a61ace230bb168b9"
}
```

Get Database

```
curl -X GET https://api.notion.com/v1/databases/task_list_database_id \
  -H "Authorization: Bearer secret_replace_with" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2021-08-16"
```

Response JSON
```
{
  "object": "database",
  "id": "e58691b6-2b1c-4d22-a6b7-aba20ba4cd03",
  "cover": null,
  "icon": {
    "type": "emoji",
    "emoji": "✔️"
  },
  "created_time": "2021-09-27T15:12:00.000Z",
  "last_edited_time": "2022-01-10T20:35:00.000Z",
  "title": [
    {
      "type": "text",
      "text": {
        "content": "Task List",
        "link": null
      },
      "annotations": {
        "bold": false,
        "italic": false,
        "strikethrough": false,
        "underline": false,
        "code": false,
        "color": "default"
      },
      "plain_text": "Task List",
      "href": null
    }
  ],
  "properties": {
    "Date Created": {
      "id": "'Y6%3C",
      "name": "Date Created",
      "type": "created_time",
      "created_time": {}
    },
    "Status": {
      "id": "%5EOE%40",
      "name": "Status",
      "type": "select",
      "select": {
        "options": [
          {
            "id": "1",
            "name": "To Do",
            "color": "red"
          },
          {
            "id": "2",
            "name": "Doing",
            "color": "yellow"
          },
          {
            "id": "3",
            "name": "Done 🙌",
            "color": "green"
          }
        ]
      }
    },
    "Name": {
      "id": "title",
      "name": "Name",
      "type": "title",
      "title": {}
    }
  },
  "parent": {
    "type": "workspace",
    "workspace": true
  },
  "url": "https://www.notion.so/task_list_database_id"
}
```

List Pages of a Database
```
curl -X POST https://api.notion.com/v1/databases/task_list_database_id/query \
  -H "Authorization: Bearer secret_replace_with" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2021-08-16" \
  --data "{
  }"
```

response
```
{
  "object": "list",
  "results": [
    {
      "object": "page",
      "id": "3f5a23be-b380-4295-a61a-ce230bb168b9",
      "created_time": "2022-01-10T20:34:00.000Z",
      "last_edited_time": "2022-01-10T20:35:00.000Z",
      "cover": null,
      "icon": null,
      "parent": {
        "type": "database_id",
        "database_id": "e58691b6-2b1c-4d22-a6b7-aba20ba4cd03"
      },
      "archived": false,
      "properties": {
        "Date Created": {
          "id": "'Y6%3C",
          "type": "created_time",
          "created_time": "2022-01-10T20:34:00.000Z"
        },
        "Status": {
          "id": "%5EOE%40",
          "type": "select",
          "select": {
            "id": "1",
            "name": "To Do",
            "color": "red"
          }
        },
        "Name": {
          "id": "title",
          "type": "title",
          "title": [
            {
              "type": "text",
              "text": {
                "content": "Yurts in Big Sur, California",
                "link": null
              },
              "annotations": {
                "bold": false,
                "italic": false,
                "strikethrough": false,
                "underline": false,
                "code": false,
                "color": "default"
              },
              "plain_text": "Yurts in Big Sur, California",
              "href": null
            }
          ]
        }
      },
      "url": "https://www.notion.so/Yurts-in-Big-Sur-California-3f5a23beb3804295a61ace230bb168b9"
    }
    ],
  "next_cursor": null,
  "has_more": false
}
```

```
curl -X POST https://api.notion.com/v1/databases/e58691b62b1c4d22a6b7aba20ba4cd03/query \
  -H "Authorization: Bearer secret" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2021-08-16" \
  --data "
{
\"filter\": {
    \"property\": \"Status\",
    \"select\": {
        \"equals\": \"To Do\"
    }
  }
}"
```