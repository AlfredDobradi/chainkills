package memory

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	store := &Store{
		keyValue: make(map[any]any),
	}

	tests := []struct {
		name        string
		set         bool
		key         string
		value       string
		expectError error
	}{
		{"hit", true, "key1", "value1", nil},
		{"miss", false, "key2", "value2", redis.Nil},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			if tt.set {
				_, err := store.Set(context.Background(), tt.key, tt.value, 0).Result()
				require.NoError(t, err)
			}

			value, err := store.Get(context.Background(), tt.key).Result()
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.Equal(t, tt.value, value)
			}
		}

		t.Run(tt.name, tf)
	}

	_, err := store.Set(context.Background(), "test_key", "test_value", 0).Result()
	require.NoError(t, err)

	value, err := store.Get(context.Background(), "test_key").Result()
	require.NoError(t, err)

	require.Equal(t, "test_value", value)
}

func TestSets(t *testing.T) {
	makeStore := func() *Store {
		return &Store{
			keyValue: map[any]any{
				"not_set":        "test",
				"already_exists": map[any]struct{}{},
				"has_items": map[any]struct{}{
					"item1": {},
					"item2": {},
				},
			},
		}
	}

	tests := []struct {
		name            string
		key             string
		members         []any
		expectedMembers []any
		expectError     error
		skip            bool
	}{
		{
			name:            "empty",
			key:             "already_exists",
			members:         []any{},
			expectedMembers: []any{},
			expectError:     nil,
		},
		{
			name:            "add items",
			key:             "has_items",
			members:         []any{"item3", "item4"},
			expectedMembers: []any{"item1", "item2", "item3", "item4"},
			expectError:     nil,
		},
		{
			name:            "unique",
			key:             "has_items",
			members:         []any{"item1", "item2"},
			expectedMembers: []any{"item1", "item2"},
			expectError:     nil,
		},
		{
			name:        "not set",
			key:         "not_set",
			members:     []any{"item1", "item2"},
			expectError: redis.Nil,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			store := makeStore()

			if len(tt.members) > 0 {
				_, err := store.SAdd(context.Background(), tt.key, tt.members...).Result()
				if tt.expectError != nil {
					require.ErrorIs(t, err, tt.expectError)
				} else {
					require.NoError(t, err)
				}
			}

			if tt.expectError == nil {
				members, err := store.SMembers(context.Background(), tt.key).Result()
				require.NoError(t, err)
				require.ElementsMatch(t, tt.expectedMembers, members)
			}
		}

		t.Run(tt.name, tf)
	}
}
