package dl

import (
	"context"
	"fmt"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/downloader"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/spf13/viper"
)

func Run(ctx context.Context, urls []string) error {
	kvd, err := kv.New(kv.Options{
		Path: consts.KVPath,
		NS:   viper.GetString(consts.FlagNamespace),
	})
	if err != nil {
		return err
	}

	c, err := tgc.New(viper.GetString(consts.FlagProxy), kvd, false, floodwait.NewSimpleWaiter())
	if err != nil {
		return err
	}

	return c.Run(ctx, func(ctx context.Context) error {
		status, err := c.Auth().Status(ctx)
		if err != nil {
			return err
		}
		if !status.Authorized {
			return fmt.Errorf("not authorized. please login first")
		}

		umsgs, err := parseURLs(ctx, c.API(), urls)
		if err != nil {
			return err
		}
		// TODO(iyear): files msgs

		return downloader.New(c.API(), viper.GetInt(consts.FlagPartSize), viper.GetInt(consts.FlagThreads), newIter(c.API(), umsgs)).
			Download(ctx, viper.GetInt(consts.FlagLimit))
	})
}
