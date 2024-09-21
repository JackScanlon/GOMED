package codes

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync/atomic"

	"snomed/src/csv"
	"snomed/src/pg"
	"snomed/src/trud"

	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

type ModelType interface {
	Concept | Description | Relationship | RefsetLang | RefsetMap | CtvMap
}

func getReader[T ModelType](obj *T) (csv.ReaderFn, error) {
	hnd, ok := any(obj).(Reader)
	if !ok {
		return nil, fmt.Errorf("failed to cast type with Reader interface")
	}

	return hnd.Reader(), nil
}

func componentFromSource[T ModelType](db *pg.Driver, fp string, tableName string, obj *T) error {
	fmt.Printf("creating Table<name: %s, from: %s>...\n", tableName, fp)

	columns, err := pg.GetColumnNamesOf(*obj)
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

	reader, err := getReader(obj)
	if err != nil {
		return err
	}
	gocsv.SetCSVReader(reader)

	group.Go(func() error {
		if err := gocsv.UnmarshalToChan(file, rowChan); err != nil {
			return err
		}

		return nil
	})

	var sz atomic.Uint32
	data := make([][]any, 0)

	var push = func() (int64, error) {
		return db.
			GetPool().
			CopyFrom(
				ctx,
				pgx.Identifier{TableSchema, tableName},
				columns,
				pgx.CopyFromRows(data),
			)
	}

	for row := range rowChan {
		hnd, ok := any(row).(Processor)
		if !ok {
			return fmt.Errorf("failed to cast type with Processor interface")
		}

		process, flat, err := hnd.Process(row)
		if err != nil {
			return err
		} else if !process {
			continue
		}

		data = append(data, flat)
		sz.Add(1)

		if sz.Load() >= ChunkSize {
			if _, err := push(); err != nil {
				return err
			}

			sz.Swap(0)
			data = make([][]any, 0)
		}
	}

	if sz.Load() > 0 {
		if _, err := push(); err != nil {
			return err
		}
	}

	return group.Wait()
}

func BuildRelease(db *pg.Driver, release *trud.Release, dir string) error {
	for _, table := range SnomedReleaseGroups {
		if !release.Category.Has(table.Category) {
			continue
		}

		tableName := fmt.Sprintf("%s_%s_%s", TablePrefix, SnomedTag, pg.GetSafeName(table.Name))

		// Tmp
		exists, err := db.Exists(TableSchema, tableName)
		if err != nil {
			return err
		} else if exists {
			continue
		}

		if err := db.DropIfExists(TableSchema, tableName); err != nil {
			return err
		}

		rv := reflect.Indirect(reflect.ValueOf(table.Model))
		rt := rv.Interface()
		if err := db.CreateTableFrom(TableSchema, tableName, rt); err != nil {
			return err
		}

		for _, filename := range table.Filenames {
			pattern := path.Join(dir, release.Metadata.Name, table.Dir, fmt.Sprintf(table.Fmt, filename))
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return err
			}

			for _, file := range matches {
				var err error
				switch obj := rt.(type) {
				case Concept:
					err = componentFromSource(db, file, tableName, &obj)
				case Description:
					err = componentFromSource(db, file, tableName, &obj)
				case Relationship:
					err = componentFromSource(db, file, tableName, &obj)
				case RefsetLang:
					err = componentFromSource(db, file, tableName, &obj)
				case RefsetMap:
					err = componentFromSource(db, file, tableName, &obj)
				case CtvMap:
					err = componentFromSource(db, file, tableName, &obj)
				default:
					err = fmt.Errorf("expected one of Concept|Relationship|Description but got %s", rv.Type().Name())
				}

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
