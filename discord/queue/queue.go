package queue

import (
	"context"
	"fmt"
	"sync"

	"git.sr.ht/~barveyhirdman/chainkills/backend"
)

type Controller struct {
	backend backend.Engine
	mx      *sync.Mutex
	Guild   map[string]*Guild
	Queue   chan any
}

func NewController() (*Controller, error) {
	backend, err := backend.Backend()
	if err != nil {
		return nil, err
	}

	return &Controller{
		backend: backend,
		mx:      &sync.Mutex{},
		Guild:   make(map[string]*Guild),
		Queue:   make(chan any, 100),
	}, nil
}

func (c *Controller) AddGuild(ctx context.Context, guildID string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.Guild[guildID]; !ok {
		c.Guild[guildID] = &Guild{
			mx:                 &sync.Mutex{},
			ID:                 guildID,
			IgnoredSystemIDs:   make([]string, 0),
			IgnoredSystemNames: make([]string, 0),
			IgnoredRegions:     make([]string, 0),
			Channels:           make([]string, 0),
			Outbox:             make(chan any, 100),
		}
	}
}

func (c *Controller) RefreshGuild(ctx context.Context, guildID string) error {
	ignoredSystemIDs := make([]string, 0)
	ignoredSystemNames := make([]string, 0)
	ignoredRegions := make([]string, 0)
	registeredChannels := make([]string, 0)

	if systemIDs, err := c.backend.GetIgnoredSystemIDs(ctx, guildID); err == nil {
		ignoredSystemIDs = append(ignoredSystemIDs, systemIDs...)
	}

	if systemNames, err := c.backend.GetIgnoredSystemNames(ctx, guildID); err == nil {
		ignoredSystemNames = append(ignoredSystemNames, systemNames...)
	}

	if regionIDs, err := c.backend.GetIgnoredRegionIDs(ctx, guildID); err == nil {
		ignoredRegions = append(ignoredRegions, regionIDs...)
	}

	if channels, err := c.backend.GetRegisteredChannelsByGuild(ctx, guildID); err == nil {
		for _, ch := range channels {
			registeredChannels = append(registeredChannels, ch.ChannelID)
		}
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if guild, ok := c.Guild[guildID]; ok {
		guild.setIgnoredSystemIDs(ignoredSystemIDs)
		guild.setIgnoredRegions(ignoredRegions)
		guild.setChannels(registeredChannels)
	} else {
		return fmt.Errorf("guild %s not found", guildID)
	}
	return nil
}
