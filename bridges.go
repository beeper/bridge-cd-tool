package main

import (
	"fmt"
	"log"
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
	BridgeInstagram      BridgeType = "instagram"
	BridgeLegacyDiscord  BridgeType = "discord"
	BridgeLegacySlack    BridgeType = "slack"
	BridgeDiscord        BridgeType = "discordgo"
	BridgeSlack          BridgeType = "slackgo"
	BridgeLinkedIn       BridgeType = "linkedin"
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
	BridgeTelegram:      defaultNotifications,
	BridgeWhatsApp:      defaultNotifications,
	BridgeFacebook:      defaultNotifications,
	BridgeGoogleChat:    defaultNotifications,
	BridgeGroupMe:       defaultNotifications,
	BridgeTwitter:       defaultNotifications,
	BridgeSignal:        defaultNotifications,
	BridgeInstagram:     defaultNotifications,
	BridgeLegacyDiscord: defaultNotifications,
	BridgeLegacySlack:   defaultNotifications,
	BridgeDiscord:       defaultNotifications,
	BridgeSlack:         defaultNotifications,
	BridgeLinkedIn:      defaultNotifications,
	BridgeHungryserv:    defaultNotifications,
	BridgeDummy: {
		{Environment: EnvDevelopment, Channel: ChannelStable},
		{Environment: EnvDevelopment, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
		{Environment: EnvStaging, Channel: ChannelStable},
		{Environment: EnvStaging, Channel: ChannelStable, Bridge: BridgeDummyWebsocket},
	},
	BridgeDummyWebsocket: {},
}

const DefaultImageTemplate = "%s:%s-amd64"

var imageTemplateOverrides = map[BridgeType]string{
	BridgeDummy:         "%s:%s",
	BridgeGroupMe:       "%s:%s",
	BridgeHungryserv:    "%s:%s",
	BridgeLegacyDiscord: "%s/discord:%s",
	BridgeLegacySlack:   "%s/slack:%s",
	BridgeLinkedIn:      "%s:%s",
}

const DefaultTargetRepoTemplate = "%s/bridge/%s"

var targetImageRepoOverrides = map[BridgeType]string{
	BridgeHungryserv: "/hungryserv",
}

func (bridgeType BridgeType) NotificationTargets() []BridgeUpdateNotification {
	notifications, ok := bridgeNotifications[bridgeType]
	if !ok || len(notifications) == 0 {
		log.Fatalf("No Beeper notifications defined for %q", bridgeType)
	}
	return notifications
}

func (bridgeType BridgeType) FormatImage(image, commit string) string {
	template, ok := imageTemplateOverrides[bridgeType]
	if !ok {
		template = DefaultImageTemplate
	}
	return fmt.Sprintf(template, image, commit)
}

func (bridgeType BridgeType) TargetRepo(registry string) string {
	repo, ok := targetImageRepoOverrides[bridgeType]
	if !ok {
		return fmt.Sprintf(DefaultTargetRepoTemplate, registry, string(bridgeType))
	}
	return registry + repo
}
