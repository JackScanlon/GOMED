package codes

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"

	"snomed/src/csv"
	"snomed/src/pg"
	"snomed/src/trud"

	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

func createFromSource[T any](db *pg.Driver, fp string, tableName string) error {
	var obj *T = new(T)
	fmt.Printf("creating Table<name: %s, from: %s>...\n", tableName, fp)

	columns, err := pg.GetColumnNamesFrom(*obj)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	rowChan := make(chan T, 1)
	group, ctx := errgroup.WithContext(context.Background())

	gocsv.SetCSVReader(func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.Comma = '\t'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		reader.TrimLeadingSpace = true
		return reader
	})

	group.Go(func() error {
		if err := gocsv.UnmarshalToChan(file, rowChan); err != nil {
			return err
		}

		return nil
	})

	var sz atomic.Uint64
	data := make([][]any, 0)

	var push = func() (int64, error) {
		return db.
			GetPool().
			CopyFrom(
				ctx,
				pgx.Identifier{CodelistSchema, tableName},
				columns,
				pgx.CopyFromRows(data),
			)
	}

	for row := range rowChan {
		flat, err := pg.FlattenRow(row)
		if err != nil {
			return err
		}

		if flat[0] == "" {
			break
		}

		data = append(data, flat)
		sz.Add(1)

		if sz.Load() >= 1000 {
			if _, err := push(); err != nil {
				return err
			}

			sz.Swap(0)
			data = make([][]any, 0)
		}
	}

	if _, err := push(); err != nil {
		return err
	}

	return group.Wait()
}

func createSnomedRelease(db *pg.Driver, release *trud.Release, dir string) error {
	for _, table := range TableMappings {
		tableName := fmt.Sprintf("%s_%s_%s", CodelistParent, CodelistPrefix, strings.ToLower(table.Name))
		if err := db.DropIfExists(CodelistSchema, tableName); err != nil {
			return err
		}

		rv := reflect.Indirect(reflect.ValueOf(table.Model))
		rt := rv.Interface()
		if err := db.CreateTableFrom(CodelistSchema, tableName, rt); err != nil {
			return err
		}

		for _, filename := range table.Filenames {
			pattern := path.Join(dir, release.Metadata.Name, CodelistFileDir, fmt.Sprintf(CodelistFileFmt, filename))
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return err
			}

			for _, file := range matches {
				var err error
				switch rt.(type) {
				case Concept:
					err = createFromSource[Concept](db, file, tableName)
				case Description:
					err = createFromSource[Description](db, file, tableName)
				case Relationship:
					err = createFromSource[Relationship](db, file, tableName)
				default:
					err = fmt.Errorf("expected one of Concept|Relationship|Description but got %s", rv.Type().Name())
				}

				if err != nil {
					return err
				}
			}
		}
	}

	// Create ICD mapping

	return nil
}

func TryCreateTables(db *pg.Driver, release *trud.Release, dir string) error {
	/*
		TODO:
			- [ ] Impl. snomed codelist generator
			- [ ] Handle code mapping table generation
	*/
	if release.IsCategory(trud.SNOMED_RELEASE) {
		return createSnomedRelease(db, release, dir)
	} else if release.IsCategory(trud.SNOMED_READ_MAP) {
		// Create table(s)

		// Create ReadCode mapping(s)

	}

	return nil
}
