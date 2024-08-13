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
	BridgeFacebook       BridgeType = "facebook"
	BridgeFacebookGo     BridgeType = "facebookgo"
	BridgeGoogleChat     BridgeType = "googlechat"
	BridgeGroupMe        BridgeType = "groupme"
	BridgeTwitter        BridgeType = "twitter"
	BridgeSignal         BridgeType = "signal"
	BridgeInstagram      BridgeType = "instagram"
	BridgeInstagramGo    BridgeType = "instagramgo"
	BridgeMeta           BridgeType = "meta"
	BridgeMetaV2         BridgeType = "metav2"
	BridgeDiscord        BridgeType = "discordgo"
	BridgeSlack          BridgeType = "slackgo"
	BridgeGoogleMessages BridgeType = "gmessages"
	BridgeLinkedIn       BridgeType = "linkedin"
	BridgeiMessageCloud  BridgeType = "imessagecloud"
	BridgeiMessagego     BridgeType = "imessagego"
	BridgeHungryserv     BridgeType = "hungryserv"
	BridgeDummy          BridgeType = "dummybridge"
	BridgeDummyWebsocket BridgeType = "dummybridgews"

	BridgeGoogleMessagesV2 BridgeType = "gmessagesv2"
)

var defaultNotifications = []BridgeUpdateNotification{
	{Environment: EnvDevelopment, Channel: ChannelStable},
	{Environment: EnvStaging, Channel: ChannelStable},
	{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true},
}

var bridgeNotifications = map[BridgeType][]BridgeUpdateNotification{
	BridgeTelegram:       defaultNotifications,
	BridgeTelegramV2:     defaultNotifications,
	BridgeWhatsApp:       defaultNotifications,
	BridgeFacebook:       defaultNotifications,
	BridgeGoogleChat:     defaultNotifications,
	BridgeGroupMe:        defaultNotifications,
	BridgeTwitter:        defaultNotifications,
	BridgeSignal:         defaultNotifications,
	BridgeInstagram:      defaultNotifications,
	BridgeiMessagego:     defaultNotifications,
	BridgeDiscord:        defaultNotifications,
	BridgeSlack:          defaultNotifications,
	BridgeGoogleMessages: {},
	BridgeGoogleMessagesV2: {
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeGoogleMessages},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeGoogleMessages},
		{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: true, Bridge: BridgeGoogleMessages},
	},
	BridgeLinkedIn:   defaultNotifications,
	BridgeHungryserv: defaultNotifications,
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
	BridgeMeta: {},
	BridgeMetaV2: {
		// These are the default notifications, but duplicated for each mode
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeFacebookGo},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeFacebookGo},
		//{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: false, Bridge: BridgeFacebookGo},
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeInstagramGo},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeInstagramGo},
		//{Environment: EnvProduction, Channel: ChannelInternal, DeployNext: false, Bridge: BridgeInstagramGo},
	},
}

const DefaultImageTemplate = "{{.Image}}:{{.Commit}}-amd64"

var imageTemplateOverrides = map[BridgeType]string{
	BridgeDummy:            "{{.Image}}:{{.Commit}}",
	BridgeGroupMe:          "{{.Image}}:{{.Commit}}",
	BridgeHungryserv:       "{{.Image}}:{{.Commit}}",
	BridgeLinkedIn:         "{{.Image}}:{{.Commit}}",
	BridgeiMessageCloud:    "{{.Commit}}",
	BridgeiMessagego:       "{{.Image}}:{{.Commit}}",
	BridgeGoogleMessagesV2: "{{.Image}}:v2-{{.Commit}}-amd64",
	BridgeTelegramV2:       "{{.Image}}:v2-{{.Commit}}-amd64",
	BridgeMetaV2:           "{{.Image}}:v2-{{.Commit}}-amd64",
}

const DefaultTargetRepoTemplate = "%s/bridge/%s"

var targetImageRepoOverrides = map[BridgeType]string{
	BridgeHungryserv:       "/hungryserv",
	BridgeGoogleMessagesV2: "/bridge/gmessages",
	BridgeTelegramV2:       "/bridge/telegramgo",
	BridgeMetaV2:           "/bridge/meta",
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
