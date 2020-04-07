package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var opt option
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "idl filename",
			Value:       "./example/hello.proto",
			Destination: &opt.Proto3FileName,
		},
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc service code",
			Destination: &opt.GenServiceCode,
		},
		cli.StringFlag{
			Name:        "o",
			Usage:       "Specify the output directory",
			Value:       "./output/",
			Destination: &opt.Output,
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		err = genMgr.Run(&opt)
		if err != nil {
			return
		}
		return
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
