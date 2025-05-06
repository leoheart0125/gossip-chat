package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gossip-chat/internal/chat"
	"gossip-chat/internal/gossip"
	"gossip-chat/internal/mdns"

	libp2p "github.com/libp2p/go-libp2p"
)

func main() {
	ctx := context.Background()
	username := flag.String("name", "anon", "Your nickname")
	flag.Parse()

	// Start P2P host
	host, err := libp2p.New()
	if err != nil {
		log.Fatalf("Failed to create host: %v", err)
	}

	// Enable LAN discovery and wait for at least one peer to be found
	discoveryChan := mdns.InitMDNS(host, "gossip-chat")
	go func() {
		for peerInfo := range discoveryChan {
			fmt.Printf("Discovered peer: %s\n", peerInfo.ID)
			if err := host.Connect(ctx, peerInfo); err != nil {
				fmt.Printf("Failed to connect to peer: %s\n", err)
			} else {
				fmt.Printf("Connected to peer: %s\n", peerInfo.ID)
			}
		}
	}()

	// Setup Gossip Chat
	gc, err := gossip.SetupGossipChat(ctx, host)
	if err != nil {
		log.Fatalf("Failed to setup gossip chat: %v", err)
	}
	defer gc.Close()

	// Start CLI interaction
	chat.StartChat(ctx, *username, gc)
}
