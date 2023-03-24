package notion

//Task represents the ToDoTask in Notion
type Task struct {
	Title string
	// the properties of Task List database is not free to add
	CustomProperties map[string]string
}
