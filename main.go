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

var GitHub = os.Getenv("GITHUB_ACTIONS") == "true"
var GitLab = os.Getenv("GITLAB_CI") == "true"

func main() {
	if GitHub {
		githubMain()
	} else if GitLab {
		gitlabMain()
	} else {
		log.Println("Unknown CI platform")
		return
	}
}

func githubMain() {
	if branch := env("GITHUB_REF_NAME"); branch != "main" && branch != "master" && branch != "xchat" && branch != "megabridge" {
		log.Println("Not notifying Beeper about update: not on main branch")
		return
	}
	bridgeType := BridgeType(env("BEEPER_BRIDGE_TYPE"))
	image := bridgeType.FormatImage(bridgeType.TargetRepo(env("CI_REGISTRY")), env("GITHUB_SHA"))
	doNotify(bridgeType, image)
}

func gitlabMain() {
	if env("CI_JOB_STAGE") != "deploy" && env("CI_JOB_STATUS") != "success" {
		log.Println("Not notifying Beeper about update: build failed")
		return
	}
	branch := env("CI_COMMIT_BRANCH")
	isLatest := branch == "main" || branch == "master" || branch == "xchat" || branch == "megabridge"
	bridgeType := BridgeType(env("BEEPER_BRIDGE_TYPE"))
	image := bridgeType.RetagImage(env("CI_REGISTRY_IMAGE"), env("CI_COMMIT_SHA"), isLatest)
	if !isLatest {
		log.Println("Not notifying Beeper about update: not on main branch")
	} else if commitMsg := env("CI_COMMIT_MESSAGE"); strings.Contains(commitMsg, "[cd skip]") || strings.Contains(commitMsg, "[skip cd]") {
		log.Println("Not notifying Beeper about update: commit message says CD should be skipped")
	} else {
		doNotify(bridgeType, image)
	}
}

func doNotify(bridgeType BridgeType, image string) {
	targets := bridgeType.NotificationTargets()
	log.Printf("Notifying %d channels about %s %s", len(targets), bridgeType, image)
	failed := false
	for _, notif := range targets {
		err := notif.Fill(bridgeType, image).Send()
		if err != nil {
			log.Printf("Failed to notify Beeper %s/%s about update to %s: %v", notif.Environment, notif.Channel, notif.Bridge, err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
