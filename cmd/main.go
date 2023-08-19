package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	fmt.Println("GO!")

	app := &cli.App{
		Name: "profit utility",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file_source",
				Value:   "test_data.json",
				Aliases: []string{"source"},
				Usage:   "File for load data",
			},
			&cli.StringFlag{
				Name:    "prediction_model",
				Value:   "first_test_model",
				Aliases: []string{"model"},
				Usage:   "Model used for predictions",
			},
			&cli.StringFlag{
				Name:     "aggregate_by",
				Aliases:  []string{"aggregate"},
				Usage:    "Data aggregation option (country|campaign)",
				Required: true,
				Action: func(ctx *cli.Context, value string) error {
					if value == "country" || value == "campaign" {
						return nil
					}
					return fmt.Errorf("aggregation value = '%s', but only 'country' or 'campaign' allowed", value)
				},
			},
		},
		Action: actionFunc,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionFunc(cCtx *cli.Context) error {
	fmt.Println(cCtx.String("source"))
	fmt.Println(cCtx.String("model"))
	fmt.Println(cCtx.String("aggregate"))

	return nil
}
