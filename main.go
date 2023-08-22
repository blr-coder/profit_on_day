package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
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

type DataGetter interface {
	GetRevenue() float64
	GetCountry() string
	GetCampaignID() string
}

func actionFunc(cCtx *cli.Context) error {
	fmt.Println("source =", cCtx.String("source"))
	fmt.Println("model =", cCtx.String("model"))
	fmt.Println("aggregate =", cCtx.String("aggregate"))

	source := cCtx.String("source")
	aggregate := cCtx.String("aggregate")

	dataFile, err := os.Open(source)
	if err != nil {
		log.Fatal("Error during os.Open: ", err)
		return err
	}
	defer dataFile.Close()

	dataArr, err := GetArr(source)
	if err != nil {
		return err
	}

	revenuesByAggMap := dataArr.ToAggMap(aggregate)

	byDayMap := make(map[string]float64)
	for agg, revenues := range revenuesByAggMap {
		// среднее значение за один день
		byDayMap[agg] = arrSum(revenues) / float64(len(revenues))
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

type Revenuers []DataGetter

func (r Revenuers) ToAggMap(aggregate string) map[string][]float64 {
	revenuesByAggMap := make(map[string][]float64)
	for _, dataStruct := range r {
		revenue := dataStruct.GetRevenue()

		if aggregate == aggregateValueCountry {
			_, ok := revenuesByAggMap[dataStruct.GetCountry()]
			if ok {
				revenuesByAggMap[dataStruct.GetCountry()] = append(revenuesByAggMap[dataStruct.GetCountry()], revenue)
				continue
			}

			revenuesByAggMap[dataStruct.GetCountry()] = []float64{revenue}
		} else {
			_, ok := revenuesByAggMap[dataStruct.GetCampaignID()]
			if ok {
				revenuesByAggMap[dataStruct.GetCampaignID()] = append(revenuesByAggMap[dataStruct.GetCampaignID()], revenue)
				continue
			}

			revenuesByAggMap[dataStruct.GetCampaignID()] = []float64{revenue}
		}
	}

	return revenuesByAggMap
}

func GetArr(source string) (Revenuers, error) {
	dataFile, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error during os.Open(): %w", err)
	}
	defer dataFile.Close()

	var res []DataGetter

	if path.Ext(source) == ".json" {
		var dataArr []*entities.JsonStruct
		err = json.NewDecoder(dataFile).Decode(&dataArr)
		if err != nil {
			return nil, fmt.Errorf("error during Decode(): %w", err)
		}

		for _, jsonStruct := range dataArr {
			res = append(res, jsonStruct)
		}

		return res, nil
	}

	var dataArr []*entities.CsvStruct
	err = gocsv.Unmarshal(dataFile, &dataArr)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal(): %w", err)
	}

	for _, csvStruct := range dataArr {
		res = append(res, csvStruct)
	}

	return res, nil
}
