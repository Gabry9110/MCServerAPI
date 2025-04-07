package main

import (
	"errors"
	"os/exec"
)

func handleService(operation string) error {
	cmd := exec.Command("systemctl", operation, SERVER_SERVICE_NAME)
	_, err := cmd.Output()
	if err != nil {
		// If there was any error with the command, handle it
		returnedError := "Service stop error: " + err.Error()
		return errors.New(returnedError)
	}
	return nil
}
