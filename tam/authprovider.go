package automaticmanager

const GOOGLE = "GOOGLE"
const ATLANSSIAN = "ATLANSSIAN"
const GITHUB = "GITHUB"

var ProviderScope = map[string]string{
	//GOOGLE:     "googlecalendar.CalendarReadonlyScope",
	GOOGLE: "https://www.googleapis.com/auth/calendar.events.readonly https://www.googleapis.com/auth/gmail.send",
	//ATLANSSIAN: "offline_access read:jira-user read:jira-work write:jira-work read:confluence-user write:confluence-content read:confluence-content.all read:confluence-space.summary",
	ATLANSSIAN: "offline_access read:page:confluence write:page:confluence read:user:confluence read:dashboard:jira read:issue:jira read:field:jira read:issue-worklog:jira read:issue-status:jira read:user:jira read:label:jira read:project:jira read:status:jira",
	GITHUB:     "repo",
}
