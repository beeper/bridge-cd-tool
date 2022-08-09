package main

import (
	"log"
	"os"
	"strings"
)

func env(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Panicf("Environment variable %q is not set", name)
	}
	return val
}

func main() {
	if branch := env("CI_COMMIT_BRANCH"); branch != "main" && branch != "master" {
		log.Println("Not notifying Beeper about update: not on main branch")
		return
	} else if env("CI_JOB_STAGE") != "deploy" && env("CI_JOB_STATUS") != "success" {
		log.Println("Not notifying Beeper about update: build failed")
		return
	}
	commitMsg := env("CI_COMMIT_MESSAGE")
	if strings.Contains(commitMsg, "[cd skip]") || strings.Contains(commitMsg, "[skip cd]") {
		log.Println("Not notifying Beeper about update: commit message says CD should be skipped")
		return
	}
	bridgeType := BridgeType(env("BEEPER_BRIDGE_TYPE"))
	image := bridgeType.FormatImage(env("CI_REGISTRY_IMAGE"), env("CI_COMMIT_SHA"))
	targets := bridgeType.NotificationTargets()
	log.Printf("Notifying %d channels about %s %s", len(targets), bridgeType, image)
	for _, notif := range targets {
		err := notif.Fill(bridgeType, image).Send()
		if err != nil {
			log.Fatalf("Failed to notify Beeper %s/%s about update to %s: %v", notif.Environment, notif.Channel, notif.Bridge, err)
		}
	}
}
