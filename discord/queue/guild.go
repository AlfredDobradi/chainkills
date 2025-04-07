package queue

import (
	"context"
	"sync"

	"git.sr.ht/~barveyhirdman/chainkills/backend"
)

type Guild struct {
	mx                 *sync.Mutex
	ID                 string `json:"id"`
	IgnoredSystemIDs   []string
	IgnoredSystemNames []string
	IgnoredRegions     []string
	Channels           []string
}

func newGuild(id string) *Guild {
	return &Guild{
		mx:                 &sync.Mutex{},
		ID:                 id,
		IgnoredSystemIDs:   make([]string, 0),
		IgnoredSystemNames: make([]string, 0),
		IgnoredRegions:     make([]string, 0),
		Channels:           make([]string, 0),
	}
}

func (g *Guild) refresh(ctx context.Context, backend backend.Engine) error {
	ignoredSystemIDs := make([]string, 0)
	ignoredSystemNames := make([]string, 0)
	ignoredRegions := make([]string, 0)
	registeredChannels := make([]string, 0)

	if systemIDs, err := backend.GetIgnoredSystemIDs(ctx, g.ID); err == nil {
		ignoredSystemIDs = append(ignoredSystemIDs, systemIDs...)
	}

	if systemNames, err := backend.GetIgnoredSystemNames(ctx, g.ID); err == nil {
		ignoredSystemNames = append(ignoredSystemNames, systemNames...)
	}

	if regionIDs, err := backend.GetIgnoredRegionIDs(ctx, g.ID); err == nil {
		ignoredRegions = append(ignoredRegions, regionIDs...)
	}

	if channels, err := backend.GetRegisteredChannelsByGuild(ctx, g.ID); err == nil {
		for _, ch := range channels {
			registeredChannels = append(registeredChannels, ch.ChannelID)
		}
	}

	g.setIgnoredSystemIDs(ignoredSystemIDs)
	g.setIgnoredSystemNames(ignoredSystemNames)
	g.setIgnoredRegions(ignoredRegions)
	g.setChannels(registeredChannels)
	return nil
}

func (g *Guild) setIgnoredSystemIDs(ignoredSystemIDs []string) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.IgnoredSystemIDs = ignoredSystemIDs
}

// Deprecating: Use setIgnoredSystemIDs instead
func (g *Guild) setIgnoredSystemNames(ignoredSystemNames []string) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.IgnoredSystemIDs = ignoredSystemNames
}

func (g *Guild) setIgnoredRegions(ignoredRegions []string) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.IgnoredRegions = ignoredRegions
}

func (g *Guild) setChannels(channels []string) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.Channels = channels
}
