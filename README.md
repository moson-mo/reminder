# reminder
## A simple background service that periodically sends notifications via D-Bus

reminder allows you to configure tasks that are executed in a specific interval.  
When a task is executed a notification is sent via D-Bus (org.freedesktop.Notifications.Notify).  
The notification message can be enriched with output of a command.

Tasks / Notifications can be configured in ~/.config/reminder/tasks.json

#### Configuration

```
{
	"Tasks": [
		{
			"Title": "Example {result}",
			"Message": "Example message\n{result}\n\nCustomize your notifications by editing tasks.json at ~/.config/reminder/",
			"TitleCommand": "echo \"Output Title\"",
			"MessageCommand": "echo \"Output Message\"",
			"ConditionCommand": "echo \"true\"",
			"Icon": "gtk-preferences",
			"Interval": 600,
			"NotificationDuration": 5
		}
	]
}
```

Option | Description
--- | ---
Title| The title of the notification message. <br>You can use `{result}` as a placeholder for the output of `TitleCommand`|
Message| The message body in the notification message. <br>You can use `{result}` as a placeholder for the output of `MessageCommand`|
TitleCommand<br>(optional)| A command that is being executed when the task runs. <br>The output can be used in `Title`|
MessageCommand<br>(optional)| A command that is being executed when the task runs. <br>The output can be used in `Message`|
ConditionCommand<br>(optional)| A command that is being executed when the task runs. <br>It needs to return the string `true` for the notification to be shown|
Icon<br>(optional)| The name of an icon that should be used in the notification|
Interval| The interval in which the task is executed and notification is shown (in seconds)|
NotificationDuration| Number of seconds a notification is shown|

Here some [examples](examples)