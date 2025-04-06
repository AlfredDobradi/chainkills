package queue

import "sync"

type Guild struct {
	mx                 *sync.Mutex
	ID                 string `json:"id"`
	IgnoredSystemIDs   []string
	IgnoredSystemNames []string
	IgnoredRegions     []string
	Channels           []string
	Outbox             chan any
}

func (g *Guild) setIgnoredSystemIDs(ignoredSystemIDs []string) {
	g.mx.Lock()
	defer g.mx.Unlock()

	g.IgnoredSystemIDs = ignoredSystemIDs
}

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
