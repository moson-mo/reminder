package reminder

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type config struct {
	Tasks []Task
}

const fileName = "tasks.json"

// loads a configuration file with tasks
func loadConfig() (*config, error) {
	content, err := os.ReadFile(path.Join(getConfDir(), fileName))
	if err != nil {
		return nil, err
	}

	conf := &config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// creates a sample configuration file
func createConfig() (*config, error) {
	conf := config{
		Tasks: []Task{
			{
				Title:                "Example {result}",
				Message:              "Example message\n{result}\n\nCustomize your notifications by editing tasks.json at ~/.config/reminder/",
				TitleCommand:         "echo \"Output Title\"",
				MessageCommand:       "echo \"Output Message\"",
				ConditionCommand:     "echo \"true\"",
				Icon:                 "gtk-preferences",
				Interval:             600,
				NotificationDuration: 5,
			},
		},
	}

	dir := getConfDir()
	file := path.Join(dir, fileName)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	content, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(file, content, 0644)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

// get's the configuration dir for the current user
func getConfDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return path.Join(dir, "/reminder/")
}
