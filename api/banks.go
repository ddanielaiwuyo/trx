package api

import (
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// Since we have a Monzo, Halifax, Barclays and others
// And the return values will eventually have to be of a specific
// type, so we can map them out on the FE and for consistency

// Today, we'll get this api to send CanonDataJSON for the Monzo and Halifax
// And then use ChartJS on the FE
type CanonDataJSON struct {
	AvgOut     float64   `json:"average_spent"`
	TotalOut   float64   `json:"total_out"`
	HighestOut float64   `json:"highest_out"`
	HighestIn  float64   `json:"highest_in"`
	AvgIn      float64   `json:"average_in"`
	TotalIn    float64   `json:"total_in"`
	OutValues  []float64 `json:"out_values"`
	InValues   []float64 `json:"in_values"`
}

type ExtraDataResponseJSON struct{}

type Bank interface {
	GetDataframe() (*dataframe.DataFrame, error)
	GetCanonData() (*CanonDataJSON, error)
	GetExtraData() (*ExtraDataResponseJSON, error)
}

type File struct {
	F        multipart.File
	BankName string
}

type MonzoBank struct {
	File
	originalDf dataframe.DataFrame
}

var (
	MonzoColumnsDrop = []string{
		"Transaction ID",
		"Emoji",
		"Currency",
		"Amount",
		"Type",
		"Notes and #tags",
	}

	renameColumns = map[string]string{
		"Amount":   "Local amount",
		"Currency": "Local currency",
		"In":       "Money In",
		"Out":      "Money Out",
	}
)

func (m *MonzoBank) GetDataframe() (*dataframe.DataFrame, error) {
	content, err := io.ReadAll(m.F)
	if err != nil {
		return nil, fmt.Errorf("Could not read file: %w", err)
	}

	df := dataframe.ReadCSV(strings.NewReader(string(content)), dataframe.HasHeader(true))
	if df.Err != nil {
		return nil, fmt.Errorf("Could not parse CSV file: %w", err)
	}

	df = df.Drop(MonzoColumnsDrop)
	if df.Err != nil {
		return nil, fmt.Errorf("Error dropping monzo_cols_drop %w", err)
	}

	m.originalDf = df

	clonedDf := df.Copy()

	if clonedDf.Err != nil {
		return nil, fmt.Errorf("Could not drop unwanted columns: %w", err)
	}

	return &clonedDf, nil
}

func (m *MonzoBank) GetCanonData() (*CanonDataJSON, error) {
	df, err := m.GetDataframe()
	if err != nil {
		return nil, err
	}
	clone := *df

	for newCol, oldCol := range renameColumns {
		clone = clone.Rename(newCol, oldCol)
		if df.Err != nil {
			return nil, df.Err
		}
	}

	var cr CanonDataJSON

	clone = clone.Filter(
		dataframe.F{
			Colname:    "In",
			Comparator: series.Greater,
			Comparando: 0,
		},
	)

	if clone.Err != nil {
		return nil, fmt.Errorf("Could not filter dataframe: %w", df.Err)
	}

	logger.Println("Cloned Filter")
	inCol := clone.Col("In")
	cr.AvgIn = math.Round(inCol.Mean()*100) / 100
	cr.TotalIn = math.Round(inCol.Sum()*100) / 100

	cr.HighestIn = math.Round(inCol.Max()*100) / 100
	cr.InValues = inCol.Float()

	if df.Err != nil {
		return nil, fmt.Errorf("Could not filter dataframe: %w", df.Err)
	}

	// We need to clone the dataframe after every filter
	newClone := m.originalDf.Copy()

	for newCol, oldCol := range renameColumns {
		newClone = newClone.Rename(newCol, oldCol)
		if df.Err != nil {
			return nil, df.Err
		}
	}

	newClone = newClone.Filter(
		dataframe.F{
			Colname:    "Out",
			Comparator: series.Less,
			Comparando: 0,
		},
	)

	outCol := newClone.Col("Out")
	cr.AvgOut = math.Round(outCol.Mean()*100) / 100
	cr.TotalOut = math.Round(outCol.Sum()*100) / 100
	cr.HighestOut = math.Round(outCol.Min()*100) / 100
	cr.OutValues = outCol.Float()

	return &cr, nil
}

func (m *MonzoBank) GetExtraData() (*ExtraDataResponseJSON, error) {
	return &ExtraDataResponseJSON{}, nil
}
