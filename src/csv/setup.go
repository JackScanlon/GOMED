package csv

import (
	"io"

	"github.com/gocarina/gocsv"
)

type ProcFn func(any) (bool, []any, error)
type ReaderFn func(io.Reader) gocsv.CSVReader
