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
	BridgeTelegramV2     BridgeType = "telegramv2"
	BridgeWhatsApp       BridgeType = "whatsapp"
	BridgeWhatsAppV2     BridgeType = "whatsappv2"
	BridgeFacebook       BridgeType = "facebook"
	BridgeFacebookGo     BridgeType = "facebookgo"
	BridgeGoogleChat     BridgeType = "googlechat"
	BridgeGroupMe        BridgeType = "groupme"
	BridgeTwitter        BridgeType = "twitter"
	BridgeSignal         BridgeType = "signal"
	BridgeInstagram      BridgeType = "instagram"
	BridgeInstagramGo    BridgeType = "instagramgo"
	BridgeMeta           BridgeType = "meta"
	BridgeDiscord        BridgeType = "discordgo"
	BridgeSlack          BridgeType = "slackgo"
	BridgeGoogleMessages BridgeType = "gmessages"
	BridgeGoogleVoice    BridgeType = "gvoice"
	BridgeLinkedIn       BridgeType = "linkedin"
	BridgeiMessageCloud  BridgeType = "imessagecloud"
	BridgeiMessagego     BridgeType = "imessagego"
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
	BridgeTelegram: {},
	BridgeTelegramV2: {
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeTelegram},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeTelegram},
		{Environment: EnvProduction, Channel: ChannelInternal, Bridge: BridgeTelegram, DeployNext: true},
	},
	BridgeWhatsApp: {},
	BridgeWhatsAppV2: {
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeWhatsApp},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeWhatsApp},
		{Environment: EnvProduction, Channel: ChannelInternal, Bridge: BridgeWhatsApp, DeployNext: false},
	},
	BridgeFacebook:       defaultNotifications,
	BridgeGoogleChat:     defaultNotifications,
	BridgeGroupMe:        defaultNotifications,
	BridgeTwitter:        defaultNotifications,
	BridgeSignal:         defaultNotifications,
	BridgeInstagram:      defaultNotifications,
	BridgeiMessagego:     defaultNotifications,
	BridgeDiscord:        defaultNotifications,
	BridgeSlack:          defaultNotifications,
	BridgeGoogleMessages: defaultNotifications,
	BridgeGoogleVoice:    defaultNotifications,
	BridgeLinkedIn:       defaultNotifications,
	BridgeHungryserv:     defaultNotifications,
	BridgeDummy: {
		{Environment: EnvDevelopment, Channel: ChannelStable},
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
		{Environment: EnvStaging, Channel: ChannelStable},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
	},
	BridgeDummyWebsocket: {},
	BridgeiMessageCloud: {
		{Environment: EnvDevelopment, Channel: ChannelStable, DeployNext: true},
		{Environment: EnvStaging, Channel: ChannelStable, DeployNext: true},
		{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true},
	},
	BridgeMeta: {
		// These are the default notifications, but duplicated for each mode
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeFacebookGo},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeFacebookGo},
		{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true, Bridge: BridgeFacebookGo},
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeInstagramGo},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeInstagramGo},
		{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true, Bridge: BridgeInstagramGo},
	},
}

const DefaultImageTemplate = "{{.Image}}:{{.Commit}}-amd64"

var imageTemplateOverrides = map[BridgeType]string{
	BridgeDummy:         "{{.Image}}:{{.Commit}}",
	BridgeGroupMe:       "{{.Image}}:{{.Commit}}",
	BridgeHungryserv:    "{{.Image}}:{{.Commit}}",
	BridgeLinkedIn:      "{{.Image}}:{{.Commit}}",
	BridgeiMessageCloud: "{{.Commit}}",
	BridgeiMessagego:    "{{.Image}}:{{.Commit}}",
	BridgeTelegramV2:    "{{.Image}}:v2-{{.Commit}}-amd64",
	BridgeWhatsAppV2:    "{{.Image}}:v2-{{.Commit}}-amd64",
}

const DefaultTargetRepoTemplate = "%s/bridge/%s"

var targetImageRepoOverrides = map[BridgeType]string{
	BridgeHungryserv: "/hungryserv",
	BridgeTelegramV2: "/bridge/telegram",
	BridgeWhatsAppV2: "/bridge/whatsapp",
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
