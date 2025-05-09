package chat

import (
	"context"
	"fmt"
	"strings"

	"gossip-chat/internal/gossip"

	"github.com/c-bata/go-prompt"
)

// MessageListener handles incoming chat messages
type MessageListener struct {
	pubsub  *gossip.PubSub
	handler func([]byte)
}

// NewMessageListener creates a new message listener
func NewMessageListener(pubsub *gossip.PubSub) *MessageListener {
	return &MessageListener{
		pubsub: pubsub,
		handler: func(data []byte) {
			fmt.Printf("\r%s\n", string(data))
		},
	}
}

// Start begins listening for messages in a goroutine
func (ml *MessageListener) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := ml.pubsub.Messages(ctx)
				if err != nil {
					fmt.Println("Error receiving message:", err)
					if ctx.Err() != nil {
						return
					}
					continue
				}
				ml.handler(msg.Data)
			}
		}
	}()
}

// StartChat 啟動聊天界面
func StartChat(ctx context.Context, username string, g *gossip.PubSub) {
	// 設置消息監聽器

	// Create and start message listener
	listener := NewMessageListener(g)
	listener.Start(ctx)

	// 命令自動完成函數
	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: "/help", Description: "Show available commands"},
			{Text: "/exit", Description: "Exit the chat"},
			{Text: "/clear", Description: "Clear the screen"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	fmt.Println("Chat started! Type your message or command (/help for available commands)")

	// 自定義選項，最小化多餘輸出
	options := []prompt.Option{
		prompt.OptionPrefix("> "),
		prompt.OptionPrefixTextColor(prompt.Green),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	}

	// 處理用戶輸入循環
	for {
		// 使用 go-prompt 讀取輸入，但發送後清除輸入行
		input := prompt.Input("", completer, options...)
		input = strings.TrimSpace(input)

		// 處理特殊命令
		if input == "/exit" {
			fmt.Println("Exiting chat...")
			return
		} else if input == "/help" {
			fmt.Println("Available commands:")
			fmt.Println("  /help  - Show this help message")
			fmt.Println("  /exit  - Exit the chat")
			fmt.Println("  /clear - Clear the screen")
			continue
		} else if input == "/clear" {
			fmt.Print("\033[H\033[2J") // ANSI clear screen
			continue
		}

		// 空消息跳過
		if input == "" {
			continue
		}

		// 發送常規消息，但不顯示在自己的終端上
		msg := fmt.Sprintf("%s: %s", username, input)
		if err := g.Publish(ctx, []byte(msg)); err != nil {
			fmt.Println("Error sending message:", err)
		}
		fmt.Print("\033[1A") // 移到上一行
		fmt.Print("\033[2K") // 清除該行
	}
}
