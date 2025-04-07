package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"git.sr.ht/~barveyhirdman/chainkills/backend"
	"git.sr.ht/~barveyhirdman/chainkills/systems"
)

var sink *Controller

type Controller struct {
	backend backend.Engine
	mx      *sync.Mutex
	guilds  map[string]*Guild
}

func initController() {
	if sink == nil {
		backend, err := backend.Backend()
		if err != nil {
			slog.Error("failed to initialize backend", "error", err)
			os.Exit(1)
		}

		sink = &Controller{
			backend: backend,
			mx:      &sync.Mutex{},
			guilds:  make(map[string]*Guild),
		}
	}
}

func (c *Controller) AddGuild(ctx context.Context, guildID string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.guilds[guildID]; !ok {
		c.guilds[guildID] = newGuild(guildID)
	}
}

func (c *Controller) RefreshGuild(ctx context.Context, guildID string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if guild, ok := c.guilds[guildID]; ok {
		if err := guild.refresh(ctx, c.backend); err != nil {
			return fmt.Errorf("failed to refresh guild: %w", err)
		}
	} else {
		return fmt.Errorf("guild not found")
	}
	return nil
}

func TargetChannels(ctx context.Context, km systems.Killmail) []string {
	if sink == nil {
		initController()
	}
	slog.Debug("finding target channels for killmail", "killmail", km.KillmailID)

	sink.mx.Lock()
	defer sink.mx.Unlock()
	channels := make([]string, 0)
	for _, guild := range sink.guilds {
		if !guild.filter(ctx, km) {
			slog.Debug("skipping guild", "guild", guild.ID)
			continue
		}
		channels = append(channels, guild.Channels...)
	}

	return channels
}
