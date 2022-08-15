package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func runCommand(command string, args ...string) {
	log.Printf("Running `%s %s`", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Error running command:", err)
	}
}

func dockerCredentials() (username, password, registry string) {
	if GitHub && env("GITHUB_REPOSITORY_OWNER") == "beeper" {
		return env("CI_REGISTRY_USERNAME"), env("CI_REGISTRY_PASSWORD"), env("CI_REGISTRY")
	} else {
		return env("BEEPER_REGISTRY_USERNAME"), env("BEEPER_REGISTRY_PASSWORD"), env("BEEPER_REGISTRY")
	}
}

func (bridgeType BridgeType) RetagImage(originalRepo, commitHash string) string {
	username, password, registry := dockerCredentials()
	runCommand("docker", "login",
		"--username", username,
		"--password", password,
		registry,
	)
	sourceImage := bridgeType.FormatImage(originalRepo, commitHash)
	targetImage := bridgeType.FormatImage(bridgeType.TargetRepo(registry), commitHash)
	runCommand("docker", "tag", sourceImage, targetImage)
	runCommand("docker", "push", targetImage)
	runCommand("docker", "rmi", targetImage)
	return targetImage
}
