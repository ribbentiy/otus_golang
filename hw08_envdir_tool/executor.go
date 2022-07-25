package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	okCode   int = 0
	failCode int = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return failCode
	}

	cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint
	fmt.Println(cmdExec)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	var appEnv []string
	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			appEnv = append(appEnv, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}
	cmdExec.Env = append(os.Environ(), appEnv...)
	if err := cmdExec.Run(); err != nil {
		return failCode
	}

	return okCode
}
