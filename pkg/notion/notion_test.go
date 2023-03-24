package notion

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

//setup load the .env, to help integration testing with the AWS credentials etc.
func setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	level, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(level)

}
func TestTaskService_AddTask(t *testing.T) {

	type args struct {
		title              string
		customerproperties map[string]string
	}
	tests := []struct {
		name string
		args args
		//	want    notion.Page
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create new to do task",
			args: args{
				title: "this is a testing task",
				customerproperties: map[string]string{
					"Source":   "Calendar",
					"Location": "http://test.com",
				},
			},
			//		want:
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notionClient, _ := NewTaskService(os.Getenv("NOTION_KEY"), os.Getenv("NOTION_DATABASE_ID"))

			tsk := Task{
				Title:            tt.args.title,
				CustomProperties: tt.args.customerproperties,
			}

			got, err := notionClient.AddTask(tsk)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("AddTask() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	setup()
	log.SetLevel(log.DebugLevel)
	type fields struct {
	}
	type args struct {
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    notion.DatabaseQueryResponse
		wantErr bool
	}{
		{
			name:    "list to do tasks",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notionClient, _ := NewTaskService(os.Getenv("NOTION_KEY"), os.Getenv("NOTION_DATABASE_ID"))

			got, err := notionClient.ListTasks()
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, task := range got {
				fmt.Printf("page: %v\n", task)
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("TaskService.ListTasks() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestTaskService_GetToDoTaskByTitle(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	type fields struct {
	}
	type args struct {
		title string
	}
	tests := []struct {
		name   string
		fields fields
		args   args

		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "test getting to do task by title",
			fields: fields{},
			args: args{
				title: "test",
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewTaskService(os.Getenv("NOTION_KEY"), os.Getenv("NOTION_DATABASE_ID"))
			got, err := p.GetToDoTaskByTitle(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.GetToDoTaskByTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("expect something not found")
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("TaskService.GetToDoTaskByTitle() = %v, want %v", got, tt.want)
			//}
		})
	}
}
