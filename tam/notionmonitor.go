package automaticmanager

import (
	"github.com/jyouturner/automaticmanager/pkg/notion"
	log "github.com/sirupsen/logrus"
)

//ProcessToDoTasks create the Task Service and run the DoneAndToDo tasks
func ProcessToDoTasks(client *notion.TaskService) {

	err := client.DoneAndCreateTasks()
	if err != nil {
		log.Fatal(err)
	}

}
