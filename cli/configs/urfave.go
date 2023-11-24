package configs

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/urfave/cli/v2"
)

type urfaveCliClient struct {
	instance *cli.App
}

func NewUrfaveCliClient() *urfaveCliClient {
	return &urfaveCliClient{instance: nil}
}

func (instance *urfaveCliClient) Instance() (*cli.App, error) {
	if instance.instance == nil {
		app := &cli.App{
			Name:                 "crawler",
			Version:              "v1.0.0",
			Usage:                "crawl metadata of VieOn Content!",
			EnableBashCompletion: true,

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "sync",
					Aliases: []string{"s"},
					Value:   "true",
					Usage:   "synchronization crawling (recommend to avoid ip blocking!)",
				},
			},

			Action: func(cCtx *cli.Context) error {
				syncFlag := cCtx.String("sync")
				if syncFlag == "true" {
					syncCrawlEndPoint := "http://localhost:3002/api/v1/contents/crawl/sync"
					if res, err := http.Get(syncCrawlEndPoint); err != nil {
						return cli.Exit("Sync crawling failure: "+err.Error(), 1)
					} else if res.StatusCode != http.StatusOK {
						defer res.Body.Close()
						return cli.Exit(fmt.Sprint("Sync crawling failure with code ", res.StatusCode, ": ", err.Error()), 1)
					} else {
						defer res.Body.Close()
						return cli.Exit("Sync Crawling successfully: ", 0)
					}
				} else if syncFlag == "false" {
					return cli.Exit("Async Crawling not support right now: ", 0)
				} else {
					return cli.Exit("Flag sync only can be true or false!", 84)
				}
			},
		}

		sort.Sort(cli.FlagsByName(app.Flags))
		sort.Sort(cli.CommandsByName(app.Commands))

		cli.HelpFlag = &cli.BoolFlag{
			Name:    "help",
			Aliases: []string{"h"},
			Usage:   "show cli tool man documentation",
		}
		cli.VersionFlag = &cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "show cli tool version",
		}
		cli.AppHelpTemplate = fmt.Sprintf("%s\nAuthor: Xuan Hoa Le a.k.a dsoFresherXuanHoa \n\n", cli.AppHelpTemplate)
		instance.instance = app
	}
	return instance.instance, nil
}
