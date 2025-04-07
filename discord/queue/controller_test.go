package queue

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~barveyhirdman/chainkills/systems"
	"github.com/stretchr/testify/require"
)

func km(systemID int) systems.Killmail {
	now := time.Now()
	return systems.Killmail{
		KillmailID: 126129718,
		Attackers: []systems.CharacterInfo{
			{CharacterID: 2115180655, CorporationID: 98224639, AllianceID: 99005678},
		},
		Victim:            systems.CharacterInfo{CharacterID: 2123233433, CorporationID: 98575144, AllianceID: 99003581},
		OriginalTimestamp: now,
		SolarSystemID:     systemID,
		Zkill: systems.Zkill{
			URL:  "https://zkillboard.com/kill/126129718/",
			Hash: "bde788de3bbfc6ea15f7653cc8f4f3bcf98e3b86",
			NPC:  false,
		},
	}
}

func TestTargetChannels(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	sink = &Controller{
		mx: &sync.Mutex{},
		guilds: map[string]*Guild{
			"1": {
				ID:               "1",
				IgnoredSystemIDs: []string{"30000142", "30000501"},
				IgnoredRegions:   []string{"10000001", "10000070"},
				Channels:         []string{"1"},
			},
			"2": {
				ID:               "2",
				IgnoredSystemIDs: []string{"30000143", "30000501"},
				IgnoredRegions:   []string{"10000002", "10000070"},
				Channels:         []string{"2"},
			},
		},
	}

	tests := []struct {
		name     string
		killmail systems.Killmail
		expected []string
		skip     bool
	}{
		{
			name:     "both guilds",
			killmail: km(30000402),
			expected: []string{"1", "2"},
		},
		{
			name:     "filtered by region in guild one",
			killmail: km(30000016),
			expected: []string{"2"},
		},
		{
			name:     "filtered by system in guild one",
			killmail: km(30000110),
			expected: []string{"2"},
		},
		{
			name:     "filtered by region in guild two",
			killmail: km(30000199),
			expected: []string{"1"},
		},
		{
			name:     "filtered by system in guild two",
			killmail: km(30000143),
			expected: []string{"1"},
		},
		{
			name:     "filtered by system in both",
			killmail: km(30000501),
			expected: nil,
		},
		{
			name:     "filtered by region in both",
			killmail: km(30010141),
			expected: nil,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			if tt.skip {
				t.Skip("skipping test")
			}
			actual := TargetChannels(context.Background(), tt.killmail)

			if tt.expected == nil {
				require.Empty(t, actual)
			} else {
				require.ElementsMatch(t, tt.expected, actual)
			}
		}
		t.Run(tt.name, tf)
	}
}
