package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/nanlei2000/npkill/pkg/file"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "npkill",
		Usage: "delete node modules",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all node modules folders",
				Action: func(cCtx *cli.Context) error {
					p := cCtx.Args().First()
					cwd, err := os.Getwd()
					if err != nil {
						log.Fatal(err)
					}
					if len(p) == 0 {
						p = cwd
					}

					if !path.IsAbs(p) {
						p = path.Join(cwd, p)
					}

					f := file.New(p)
					f.FindNodeModules()

					return nil
				},
			},
			{
				Name:    "del",
				Aliases: []string{"d"},
				Usage:   "delete node modules folder",
				Action: func(cCtx *cli.Context) error {
					p := cCtx.Args().First()
					cwd, err := os.Getwd()
					if err != nil {
						log.Fatal(err)
					}
					if len(p) == 0 {
						p = cwd
					}
					if !path.IsAbs(p) {
						p = path.Join(cwd, p)
					}
					p = path.Clean(p)
					if filepath.Base(p) != "node_modules" {
						return nil
					}

					err = os.RemoveAll(p)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Removed: %s\n", p)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
