package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/sajari/regression"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"profit_on_day/entities"
	"profit_on_day/utils"
)

const (
	aggregateValueCountry  = "country"
	aggregateValueCampaign = "campaign"

	modelSLE = "sle" //simple linear extrapolation
	modelLR  = "lr"  // linear regression
)

func main() {
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
				Value:   "sle",
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
		Action: runCalculation,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type DataGetter interface {
	GetRevenue() float64
	GetCountry() string
	GetCampaignID() string
	GetLtv(i int) float64
}

func runCalculation(cCtx *cli.Context) error {
	source := cCtx.String("source")
	aggregate := cCtx.String("aggregate")
	model := cCtx.String("model")

	switch model {
	case modelSLE:
		if err := CalculateSLE(source, aggregate); err != nil {
			return err
		}
	case modelLR:
		if err := CalculateLR(source, aggregate); err != nil {
			return err
		}
	}

	return nil
}

func CalculateLR(source, aggregate string) error {
	dataArr, err := GetArr(source)
	if err != nil {
		return err
	}

	// Группировка данных по agg
	groupAggData := make(map[string][]DataGetter)
	for _, dp := range dataArr {
		if aggregate == aggregateValueCountry {
			groupAggData[dp.GetCountry()] = append(groupAggData[dp.GetCountry()], dp)
		} else {
			groupAggData[dp.GetCampaignID()] = append(groupAggData[dp.GetCampaignID()], dp)
		}
	}

	// Прогнозирование для каждой agg
	for agg, dpList := range groupAggData {
		// Создание регрессионной модели
		model := new(regression.Regression)
		model.SetObserved("Ltv7")
		model.SetVar(0, "Ltv1")
		model.SetVar(1, "Ltv2")
		model.SetVar(2, "Ltv3")
		model.SetVar(3, "Ltv4")
		model.SetVar(4, "Ltv5")
		model.SetVar(5, "Ltv6")
		model.SetVar(6, "Ltv7")

		// Добавление данных в модель
		for _, dp := range dpList {
			model.Train(regression.DataPoint(dp.GetLtv(7), []float64{dp.GetLtv(1), dp.GetLtv(2), dp.GetLtv(3), dp.GetLtv(4), dp.GetLtv(5), dp.GetLtv(6), dp.GetLtv(7)}))
		}

		err = model.Run()
		if err != nil {
			return err
		}

		inputData := []float64{
			dpList[0].GetLtv(1),
			dpList[0].GetLtv(2),
			dpList[0].GetLtv(3),
			dpList[0].GetLtv(4),
			dpList[0].GetLtv(5),
			dpList[0].GetLtv(6),
			dpList[0].GetLtv(7),
		}
		// Дополняем вектор средними значениями до 60
		for len(inputData) < 60 {
			average := (inputData[0] + inputData[1] + inputData[2] + inputData[3] + inputData[4] + inputData[5] + inputData[6]) / 7
			inputData = append(inputData, average)
		}

		// Прогнозирование значения Ltv60
		prediction, err := model.Predict(inputData)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %f\n", agg, prediction)
	}

	return nil
}

func CalculateSLE(source, aggregate string) error {
	dataArr, err := GetArr(source)
	if err != nil {
		return err
	}

	revenuesByAggMap := dataArr.ToRevenueByAggMap(aggregate)

	byDayMap := make(map[string]float64)
	for agg, revenues := range revenuesByAggMap {
		// среднее значение за один день
		byDayMap[agg] = utils.ArrSum(revenues) / float64(len(revenues))
	}

	for agg, rev := range byDayMap {
		fmt.Printf("%s: %f\n", agg, rev*60)
	}

	return nil
}

type Revenuers []DataGetter

func (r Revenuers) ToRevenueByAggMap(aggregate string) map[string][]float64 {
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
