package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
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
	fmt.Println("source =", cCtx.String("source"))
	fmt.Println("model =", cCtx.String("model"))
	fmt.Println("aggregate =", cCtx.String("aggregate"))

	byDayMap, err := generateAvgByDayMap(cCtx.String("source"), cCtx.String("aggregate"))
	if err != nil {
		return err
	}

	fmt.Println("Среднее за 1 день, исходя из статистики за 7 дней")
	for agg, rev := range byDayMap {
		fmt.Printf("%s: %v\n", agg, rev)
	}
	fmt.Println("=================================")
	fmt.Println("Предпологаемое за 60 дней")
	for agg, rev := range byDayMap {
		fmt.Printf("%s: %v\n", agg, rev*60)
	}

	return nil
}

func generateAvgByDayMap(source, aggregate string) (map[string]float64, error) {
	avgByDayMap := make(map[string]float64)

	revenuesByAggMap := make(map[string][]float64)
	// TODO: Add better solution for understand source type (.ext)
	if source == "test_data.json" {
		// TODO: If we have very big file we should read it by strings
		content, err := os.ReadFile(fmt.Sprintf("./%s", source))
		if err != nil {
			log.Println("Error when opening file: ", err)
			return nil, err
		}

		var payload []entities.JsonStruct
		err = json.Unmarshal(content, &payload)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
			return nil, err
		}

		for _, sourceStruct := range payload {
			revenue := sourceStruct.Revenue()

			if aggregate == aggregateValueCountry {
				_, ok := revenuesByAggMap[sourceStruct.Country]
				if ok {
					revenuesByAggMap[sourceStruct.Country] = append(revenuesByAggMap[sourceStruct.Country], revenue)
					continue
				}

				revenuesByAggMap[sourceStruct.Country] = []float64{revenue}
			} else {
				_, ok := revenuesByAggMap[sourceStruct.CampaignId]
				if ok {
					revenuesByAggMap[sourceStruct.CampaignId] = append(revenuesByAggMap[sourceStruct.CampaignId], revenue)
					continue
				}

				revenuesByAggMap[sourceStruct.CampaignId] = []float64{revenue}
			}
		}
	} else {
		content, err := os.Open("test_data.csv")
		if err != nil {
			panic(err)
		}
		defer content.Close()

		var payload []entities.CsvStruct
		err = gocsv.Unmarshal(content, &payload)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
			return nil, err
		}

		for _, sourceStruct := range payload {
			revenue := sourceStruct.Revenue()

			if aggregate == aggregateValueCountry {
				_, ok := revenuesByAggMap[sourceStruct.Country]
				if ok {
					revenuesByAggMap[sourceStruct.Country] = append(revenuesByAggMap[sourceStruct.Country], revenue)
					continue
				}

				revenuesByAggMap[sourceStruct.Country] = []float64{revenue}
			} else {
				_, ok := revenuesByAggMap[sourceStruct.CampaignId]
				if ok {
					revenuesByAggMap[sourceStruct.CampaignId] = append(revenuesByAggMap[sourceStruct.CampaignId], revenue)
					continue
				}

				revenuesByAggMap[sourceStruct.CampaignId] = []float64{revenue}
			}
		}
	}

	for agg, revenues := range revenuesByAggMap {
		// среднее значение за одни день
		avgByDayMap[agg] = arrSum(revenues) / float64(len(revenues))
	}

	return avgByDayMap, nil
}

func arrSum(arr []float64) float64 {
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
