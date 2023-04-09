package main

import (
	"fmt"
	"log"
	"os"
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

					absPath, err := filepath.Abs(p)
					if err != nil {
						log.Fatal(err)
					}

					f := file.New(absPath)
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

					absPath, err := filepath.Abs(p)
					if err != nil {
						log.Fatal(err)
					}
					if filepath.Base(absPath) != "node_modules" {
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
