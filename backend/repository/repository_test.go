package repository

import (
	"context"
	"testing"

	"git.sr.ht/~barveyhirdman/chainkills/config"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *Backend {
	err := config.Read("./testdata/config.test.yaml")
	require.NoError(t, err)

	backend, err := New()
	require.NoError(t, err)

	return backend
}

func TestAddKillmail(t *testing.T) {
	backend := setup(t)

	tests := map[string]struct {
		add      string
		exists   string
		expected bool
	}{
		"add and exists": {
			add:      "12345",
			exists:   "12345",
			expected: true,
		},
		"does not exist": {
			add:      "",
			exists:   "54321",
			expected: false,
		},
	}

	for name, tt := range tests {
		tf := func(t *testing.T) {
			if tt.add != "" {
				err := backend.AddKillmail(context.Background(), tt.add)
				require.NoError(t, err)
			}

			ok, err := backend.KillmailExists(context.Background(), tt.exists)
			require.NoError(t, err)
			require.Equal(t, tt.expected, ok)
		}

		t.Run(name, tf)
	}
}

func TestIgnoredEntities(t *testing.T) {
	backend := setup(t)

	tests := map[string]struct {
		kind     string
		add      []int64
		expected []string
	}{
		"add and get system id": {
			kind:     "system_id",
			add:      []int64{10000001, 10000002},
			expected: []string{"10000001", "10000002"},
		},
		"add and get region id": {
			kind:     "region_id",
			add:      []int64{30000001, 30000002},
			expected: []string{"30000001", "30000002"},
		},
	}

	for name, tt := range tests {
		tf := func(t *testing.T) {
			var ignore func(context.Context, int64) error

			var get func(context.Context) ([]string, error)

			switch tt.kind {
			default:
				fallthrough
			case "system_id":
				ignore = backend.IgnoreSystemID
				get = backend.GetIgnoredSystemIDs
			case "region_id":
				ignore = backend.IgnoreRegionID
				get = backend.GetIgnoredRegionIDs
			}

			if len(tt.add) > 0 {
				for _, id := range tt.add {
					require.NoError(t, ignore(context.Background(), id))
				}
			}

			ids, err := get(context.Background())
			require.NoError(t, err)
			require.ElementsMatch(t, tt.expected, ids)
		}

		t.Run(name, tf)
	}
}
