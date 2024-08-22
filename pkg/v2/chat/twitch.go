package chat

import (
	"context"
	"log/slog"

	"github.com/gempir/go-twitch-irc/v4"
)

func toChatMsg(twitchMsg twitch.PrivateMessage) ChatMsg {
	return ChatMsg{
		Name: twitchMsg.User.DisplayName,
		Bits: twitchMsg.Bits,
		Msg:  twitchMsg.Message,
	}
}

func NewTwitchChat(ctx context.Context, chat string) (chan ChatMsg, error) {
	messages := make(chan ChatMsg)

	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		messages <- toChatMsg(msg)
	})

	client.Join(chat)

	// TODO: on disconnect send done to reconnect everything from the top
	go func() {
		slog.Warn("NewTwitchChat# connecting...")
		err := client.Connect()
		slog.Warn("twitch client disconnected", "err", err)
	}()

	go func() {
		<-ctx.Done()
		slog.Debug("NewTwitchChat#disconnected via context")
		client.Disconnect()
	}()

	return messages, nil
}
