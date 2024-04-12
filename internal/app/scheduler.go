package app

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"time"
)

const (
	durationToCleanBanners = 5 * time.Hour
	//durationToCleanBannersLocal = 30 * time.Second
)

type DeletedBannersCleaner interface {
	CleanDeletedBanners(ctx context.Context) error
}

func startCleaningTask(scheduler gocron.Scheduler, deletedBannersCleaner DeletedBannersCleaner, log *slog.Logger, ctx context.Context) error {
	if _, err := scheduler.NewJob(
		gocron.DurationJob(durationToCleanBanners),
		gocron.NewTask(
			func(cleaner DeletedBannersCleaner, log *slog.Logger, ctx context.Context) {
				if err := cleaner.CleanDeletedBanners(ctx); err != nil {
					log.Error("error cleaning deleted banners")
				}
				log.Info("deleted banner was cleaned by cron job")
			},
			deletedBannersCleaner,
			log,
			ctx,
		),
	); err != nil {
		return err
	}
	return nil
}
