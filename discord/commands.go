package discord

import (
	"context"
	"fmt"
	"log/slog"

	"git.sr.ht/~barveyhirdman/chainkills/backend/repository"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var packageName = "git.sr.ht/~barveyhirdman/chainkills/discord"

var IgnoreSystemIDCommand = &discordgo.ApplicationCommand{
	Name:        "ignore-system-id",
	Description: "Ignore a system by ID",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "system_id",
			Description: "The ID of the system to ignore",
			Required:    true,
		},
	},
}

var IgnoreRegionIDCommand = &discordgo.ApplicationCommand{
	Name:        "ignore-region-id",
	Description: "Ignore a region by ID",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "region_id",
			Description: "The ID of the region to ignore",
			Required:    true,
		},
	},
}

func HandleIgnoreSystemID(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sctx, span := otel.Tracer(packageName).Start(context.Background(), "HandleIgnoreSystemID")
	defer span.End()

	backend, err := repository.New()
	if err != nil {
		slog.Error("failed to get backend", "error", err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return
	}

	if err := backend.IgnoreSystemID(sctx, i.ApplicationCommandData().Options[0].IntValue()); err != nil {
		slog.Error("failed to add ignored system id", "error", err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return
	}

	systemID := i.ApplicationCommandData().Options[0].IntValue()

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"System ID %d has been ignored",
				systemID,
			),
		},
	}); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slog.Error("failed to respond to interaction", "error", err)
	}

	span.SetStatus(codes.Ok, "ok")
}

func HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, span := otel.Tracer(packageName).Start(context.Background(), "HandleSlashCommand")
	defer span.End()

	// Check if the command is a slash command
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ignore-system-id":
		HandleIgnoreSystemID(ctx, s, i)
	case "ignore-region-id":
		HandleIgnoreRegionID(ctx, s, i)
	}
}

func HandleIgnoreRegionID(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sctx, span := otel.Tracer(packageName).Start(context.Background(), "HandleIgnoreRegionID")
	defer span.End()

	span.SetStatus(codes.Ok, "ok")

	regionID := i.ApplicationCommandData().Options[0].IntValue()

	response, err := ignoreEntityID(sctx, "region_id", regionID)
	if err != nil {
		slog.Error("failed to add region to ignored entity list", "id", regionID, "error", err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		slog.Error("failed to respond to interaction", "error", err)
	}
}

func ignoreEntityID(ctx context.Context, kind string, value int64) (*discordgo.InteractionResponse, error) {
	sctx, span := otel.Tracer(packageName).Start(ctx, "ignoreEntityID")
	defer span.End()

	backend, err := repository.New()
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to ignore %s: %d", kind, value),
			},
		}, err
	}

	message := ""

	switch kind {
	case "system":
		err = backend.IgnoreSystemID(sctx, value)
		message = fmt.Sprintf("System ID %d has been ignored", value)
	case "region":
		err = backend.IgnoreRegionID(sctx, value)
		message = fmt.Sprintf("Region ID %d has been ignored", value)
	default:
		err = fmt.Errorf("unknown kind: %s", kind)
	}

	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to ignore %s: %d", kind, value),
			},
		}, err
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}

	return response, nil
}
