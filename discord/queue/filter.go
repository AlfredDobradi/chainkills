package queue

import (
	"context"
	"fmt"
	"log/slog"

	"git.sr.ht/~barveyhirdman/chainkills/systems"
)

func (g *Guild) filter(ctx context.Context, km systems.Killmail) bool {
	filters := map[string]func(context.Context, systems.Killmail) bool{
		"system_name": g.filterBySystemName,
		"region_id":   g.filterByRegionID,
		"system_id":   g.filterBySystemID,
		// g.filterNPC,
	}

	for reason, filter := range filters {
		if !filter(ctx, km) {
			slog.Debug("filtered out killmail", "killmail", km.KillmailID, "guild", g.ID, "reason", reason)
			return false
		}
	}

	return true
}

func (g *Guild) filterNPC(ctx context.Context, killmail systems.Killmail) bool { //nolint:unused
	return !killmail.Zkill.NPC
}

func (g *Guild) filterBySystemID(ctx context.Context, killmail systems.Killmail) bool {
	for _, id := range g.IgnoredSystemIDs {
		if id == fmt.Sprintf("%d", killmail.SolarSystemID) {
			return false
		}
	}

	return true
}

func (g *Guild) filterBySystemName(ctx context.Context, killmail systems.Killmail) bool {
	sys, ok := systems.GetSystem(killmail.SolarSystemID)
	if !ok {
		return false
	}

	for _, name := range g.IgnoredSystemNames {
		if name == sys.SystemName {
			return false
		}
	}

	return true
}

func (g *Guild) filterByRegionID(ctx context.Context, killmail systems.Killmail) bool {
	sys, ok := systems.GetSystem(killmail.SolarSystemID)
	if !ok {
		return false
	}

	for _, id := range g.IgnoredRegions {
		if id == fmt.Sprintf("%d", sys.RegionID) {
			return false
		}
	}

	return true
}
