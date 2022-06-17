package reminder

type Task struct {
	Title                string
	Message              string
	TitleCommand         string
	MessageCommand       string
	ConditionCommand     string
	Icon                 string
	Interval             int
	NotificationDuration int
}
