package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type ChartInfo struct {
	Name string `json:"name"`
}

type Aggregate struct {
	AvgSpent  int       `json:"average_spent"`
	AvgIn     uint      `json:"average_in"`
	TotalIn   uint      `json:"total_in"`
	TotalOut  int       `json:"total_out"`
	ChartInfo ChartInfo `json:"chart_info"`
}

func processFile(f *File, w http.ResponseWriter, r *http.Request) {
	// content := f.file.Read()
	content, err := io.ReadAll(f.file)
	if err != nil {
		log.Println("failed to read file content: ", err)
		http.ServeFile(w, r, "index.html")
		return
	}

	sr := strings.NewReader(string(content))
	df := dataframe.ReadCSV(
		sr, dataframe.HasHeader(true),
	)

	droppedCols := []string{
		"Transaction ID",
		"Emoji",
		"Currency", // Already have Local Currency Column
		"Amount",   // Already have Local Amount Column
		"Type",     // usually Direct Debit, Card, FasterPayment
		"Notes and #tags",
	}
	alteredDf := df.Drop(droppedCols)

	renameMaps := map[string]string{
		"Local amount":   "Amount",
		"Local currency": "Currency",
		"Money In":       "In",
		"Money Out":      "Out",
	}

	// this is memory inefficient unless the src code
	// is passing pointers, which I doubt
	var df2 dataframe.DataFrame
	df2 = alteredDf
	for from, to := range renameMaps {
		df2 = df2.Rename(to, from)
		if df2.Err != nil {
			log.Printf("could not rename %s | %s because: %s\n", from, to, df2.Err)
		}
	}

	log.Printf("renamed_df: \n%s\n", df2.String())

	if alteredDf.Err != nil {
		log.Println("here- .")
		log.Fatal(df2.Err)
	}

	log.Printf("new nams -> %+v\n", df2.Names())
	fil := df2.Filter(
		dataframe.F{
			Colname:    "In",
			Comparator: series.Greater,
			Comparando: 0,
		},
	)

	totalIn := fil.Col("In").Sum()
	log.Printf("%+v\n", totalIn)

	fil = df2.Filter(
		dataframe.F{
			Colname:    "Out",
			Comparator: series.Less,
			Comparando: 0,
		},
	)
	if fil.Err != nil {
		log.Println("fileter error")
		log.Fatal(fil.Err)
	}

	totalExpenses := fil.Col("Out").Sum()
	log.Printf("%+v\n", totalExpenses)
}
