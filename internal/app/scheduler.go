package app

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"time"
)

type DeletedBannersCleaner interface {
	CleanDeletedBanners(ctx context.Context) error
}

func StartCleaningTask(
	scheduler gocron.Scheduler,
	deletedBannersCleaner DeletedBannersCleaner,
	log *slog.Logger,
	ctx context.Context,
	intervalToCleanBanners time.Duration,
) error {
	_, err := scheduler.NewJob(
		gocron.DurationJob(intervalToCleanBanners),
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
	)
	if err != nil {
		return err
	}
	return nil
}
