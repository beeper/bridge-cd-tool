package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type BeeperEnv string

const (
	EnvDevelopment BeeperEnv = "DEV"
	EnvStaging     BeeperEnv = "STAGING"
	EnvProduction  BeeperEnv = "PROD"
)

type BeeperChannel string

const (
	ChannelStable   BeeperChannel = "STABLE"
	ChannelNightly  BeeperChannel = "NIGHTLY"
	ChannelInternal BeeperChannel = "INTERNAL"
)

type BridgeType string

const (
	BridgeTelegram       BridgeType = "telegram"
	BridgeWhatsApp       BridgeType = "whatsapp"
	BridgeFacebook       BridgeType = "facebook"
	BridgeGoogleChat     BridgeType = "googlechat"
	BridgeGroupMe        BridgeType = "groupme"
	BridgeTwitter        BridgeType = "twitter"
	BridgeSignal         BridgeType = "signal"
	BridgeSignald        BridgeType = "signald"
	BridgeInstagram      BridgeType = "instagram"
	BridgeDiscord        BridgeType = "discordgo"
	BridgeSlack          BridgeType = "slackgo"
	BridgeLinkedIn       BridgeType = "linkedin"
	BridgeImessageCloud  BridgeType = "cloud-mac-stack"
	BridgeHungryserv     BridgeType = "hungryserv"
	BridgeDummy          BridgeType = "dummybridge"
	BridgeDummyWebsocket BridgeType = "dummybridgews"
)

var defaultNotifications = []BridgeUpdateNotification{
	{Environment: EnvDevelopment, Channel: ChannelStable},
	{Environment: EnvStaging, Channel: ChannelStable},
	{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true},
}

var bridgeNotifications = map[BridgeType][]BridgeUpdateNotification{
	BridgeTelegram:   defaultNotifications,
	BridgeWhatsApp:   defaultNotifications,
	BridgeFacebook:   defaultNotifications,
	BridgeGoogleChat: defaultNotifications,
	BridgeGroupMe:    defaultNotifications,
	BridgeTwitter:    defaultNotifications,
	BridgeSignal:     defaultNotifications,
	BridgeSignald:    defaultNotifications,
	BridgeInstagram:  defaultNotifications,
	BridgeDiscord:    defaultNotifications,
	BridgeSlack:      defaultNotifications,
	BridgeLinkedIn:   defaultNotifications,
	BridgeHungryserv: defaultNotifications,
	BridgeDummy: {
		{Environment: EnvDevelopment, Channel: ChannelStable},
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
		{Environment: EnvStaging, Channel: ChannelStable},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
	},
	BridgeDummyWebsocket: {},
	BridgeImessageCloud: {
		{Environment: EnvDevelopment, Channel: ChannelStable, DeployNext: true},
		{Environment: EnvStaging, Channel: ChannelStable, DeployNext: true},
		{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true},
	},
}

const DefaultImageTemplate = "{{.Image}}:{{.Commit}}-amd64"

var imageTemplateOverrides = map[BridgeType]string{
	BridgeDummy:         "{{.Image}}:{{.Commit}}",
	BridgeGroupMe:       "{{.Image}}:{{.Commit}}",
	BridgeHungryserv:    "{{.Image}}:{{.Commit}}",
	BridgeLinkedIn:      "{{.Image}}:{{.Commit}}",
	BridgeSignald:       "{{.Image}}:{{.Commit}}",
	BridgeImessageCloud: "{{.Commit}}",
}

const DefaultTargetRepoTemplate = "%s/bridge/%s"

var targetImageRepoOverrides = map[BridgeType]string{
	BridgeHungryserv: "/hungryserv",
	BridgeSignald:    "/signald",
}

func (bridgeType BridgeType) NotificationTargets() []BridgeUpdateNotification {
	notifications, ok := bridgeNotifications[bridgeType]
	if !ok || len(notifications) == 0 {
		log.Fatalf("No Beeper notifications defined for %q", bridgeType)
	}
	return notifications
}

func (bridgeType BridgeType) FormatImage(image, commit string) string {
	templateString, ok := imageTemplateOverrides[bridgeType]
	if !ok {
		templateString = DefaultImageTemplate
	}

	var result bytes.Buffer
	tmpl := template.Must(template.New("t").Parse(templateString))
	err := tmpl.Execute(&result, map[string]string{
		"Image":  image,
		"Commit": commit,
	})

	if err != nil {
		log.Fatalf("Failed to format image for %q", bridgeType)
	}
	return result.String()
}

func (bridgeType BridgeType) TargetRepo(registry string) string {
	repo, ok := targetImageRepoOverrides[bridgeType]
	if !ok {
		return fmt.Sprintf(DefaultTargetRepoTemplate, registry, string(bridgeType))
	}
	return registry + repo
}
