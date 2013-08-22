package main

import (
	"os/exec"
)

func gitCommand() exec.Cmd {
	return exec.Command("git")
}

func gitReceivePackCommand() exec.Cmd {
	return exec.Command("git", "receive-pack")
}
