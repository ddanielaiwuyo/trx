package main

import (
	"io"
	"log"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type ChartInfo struct {
	Name string `json:"name"`
}

type Response struct {
	AvgOut     float64   `json:"average_spent"`
	TotalOut   float64   `json:"total_out"`
	HighestOut float64   `json:"highest_out"`
	HighestIn  float64   `json:"highest_in"`
	AvgIn      float64   `json:"average_in"`
	TotalIn    float64   `json:"total_in"`
	OutValues  []float64 `json:"out_values"`
	InValues   []float64 `json:"in_values"`

	ChartInfo ChartInfo `json:"chart_info"`
}

func (f *File) handleMonzoBankStatement() (*Response, error) {
	content, err := io.ReadAll(f.file)
	if err != nil {
		return nil, err
	}

	originalDf := dataframe.ReadCSV(
		strings.NewReader(string(content)),
		dataframe.HasHeader(true),
	)

	df := originalDf.Copy()
	df = df.Drop([]string{
		"Transaction ID",
		"Emoji",
		"Currency", // Already have Local Currency Column
		"Amount",   // Already have Local Amount Column
		"Type",     // usually Direct Debit, Card, FasterPayment
		"Notes and #tags",
	})

	renameMap := map[string]string{
		"Amount":   "Local amount",
		"Currency": "Local currency",
		"In":       "Money In",
		"Out":      "Money Out",
	}

	for tgt, old := range renameMap {
		df = df.Rename(tgt, old)
		if df.Err != nil {
			return nil, df.Err
		}
	}

	log.Println("altered orignal data")
	log.Printf("%s\n", df.String())

	var res Response

	fil := df.Filter(
		dataframe.F{
			Colname:    "In",
			Comparator: series.Greater,
			Comparando: 0,
		},
	)

	if fil.Err != nil {
		log.Println("error occured while filtering for In col")
		return nil, fil.Err
	}

	inCol := fil.Col("In")

	res.AvgIn = inCol.Mean()
	res.TotalIn = inCol.Sum()
	res.HighestIn = inCol.Max()
	res.InValues = inCol.Float()

	fil = df.Filter(
		dataframe.F{
			Colname:    "Out",
			Comparator: series.Less,
			Comparando: 0,
		},
	)

	if fil.Err != nil {
		log.Println("error occured while filtering for Out col")
		return nil, fil.Err
	}

	outCol := fil.Col("Out")
	res.AvgOut = outCol.Mean()
	res.TotalOut = outCol.Sum()
	res.HighestOut = outCol.Min()
	res.OutValues = outCol.Float()
	res.ChartInfo.Name = "Testing Chart Name To Generate"

	return &res, nil
}
