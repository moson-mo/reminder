package reminder

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
)

var con *dbus.Conn

// Start starts the notifier
func Start() error {
	var err error
	con, err = dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	defer con.Close()

	conf, err := loadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			conf, err = createConfig()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, task := range conf.Tasks {
		if task.NotificationDuration*1000 > math.MaxInt32 {
			fmt.Println("Duration setting for " + task.Title + " is too large. Skipping")
			continue
		}
		go runTask(task)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	return nil
}

// runs a task (to be started from a go-routine)
func runTask(task Task) {
	for {
		if task.ConditionCommand != "" {
			cRes, err := runCommand(task.ConditionCommand)
			if err != nil {
				fmt.Println("Error running condition command", err)
			}
			if !strings.HasPrefix(cRes, "true") {
				time.Sleep(time.Second * time.Duration(task.Interval))
				continue
			}
		}

		tRes := ""
		mRes := ""
		if task.TitleCommand != "" {
			var err error
			tRes, err = runCommand(task.TitleCommand)
			if err != nil {
				fmt.Println("Error running title command", err)
			}
		}
		if task.MessageCommand != "" {
			var err error
			mRes, err = runCommand(task.MessageCommand)
			if err != nil {
				fmt.Println("Error running message command", err)
			}
		}

		title := strings.ReplaceAll(task.Title, "{result}", tRes)
		message := strings.ReplaceAll(task.Message, "{result}", mRes)

		err := notify(title, message, task.Icon, task.NotificationDuration)
		if err != nil {
			fmt.Println("Error sending notification", err)
		}
		time.Sleep(time.Second * time.Duration(task.Interval))
	}
}

// get users shell
func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh" // fallback
	}
	return shell
}

// execute command and return output
func runCommand(command string) (string, error) {
	shell := getShell()
	args := []string{"-c", command}
	cmd := exec.Command(shell, args...)
	result, err := cmd.Output()
	if err != nil {
		return err.Error(), err
	}
	return string(result), nil
}

// send notification via D-Bus
func notify(title, message, icon string, duration int) error {
	obj := con.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
		icon, title, message, []string{},
		map[string]dbus.Variant{}, int32(1000*duration))
	if call.Err != nil {
		return call.Err
	}
	return nil
}
