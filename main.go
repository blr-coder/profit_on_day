package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"profit_on_day/entities"
)

const (
	aggregateValueCountry  = "country"
	aggregateValueCampaign = "campaign"
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
					if value == aggregateValueCountry || value == aggregateValueCampaign {
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

	resMap := make(map[string][]float64)

	content, err := os.ReadFile(fmt.Sprintf("./%s", cCtx.String("source")))
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []entities.JsonStruct
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for _, sourceStruct := range payload {
		revenue := sourceStruct.Revenue()

		_, ok := resMap[sourceStruct.Country]
		if ok {
			resMap[sourceStruct.Country] = append(resMap[sourceStruct.Country], revenue)
			continue
		}

		resMap[sourceStruct.Country] = []float64{revenue}
	}

	for agg, listByAgg := range resMap {
		avgRevenueBy7days := sum(listByAgg) / float64(len(listByAgg))
		avgRevenueByOneDay := avgRevenueBy7days / 7
		avgRevenueBy60Days := avgRevenueByOneDay * 60
		fmt.Printf("%s: %v\n", agg, avgRevenueBy60Days)
	}

	return nil
}

func generateResMap(source, model, aggregate string) map[string][]float64 {
	resMap := make(map[string][]float64)

	// open file, unmarshal data to []struct, loop by []struct and add agg data to resMap

	return resMap
}

func sum(arr []float64) float64 {
	var sum float64
	idx := 0
	for {
		if idx > len(arr)-1 {
			break
		}
		sum += arr[idx]
		idx++
	}
	return sum
}
