package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"snomed/src/pg"
)

type Template struct {
	Name    string
	Data    map[string]any
	Content *template.Template
	PgOpts  pg.PgOptions
	HasOpts bool
	Echo    bool
}

type TmplOption func(*Template) error

func WithData(data map[string]any) TmplOption {
	return func(t *Template) error {
		for k, v := range data {
			t.Data[k] = v
		}

		return nil
	}
}

func WithPgOpts(opts ...pg.PgOption) TmplOption {
	return func(t *Template) error {
		db, err := pg.TryGetDB()
		if err != nil {
			return err
		}

		t.PgOpts = db.GetOptions(opts...)
		t.HasOpts = true
		return nil
	}
}

func WithEcho() TmplOption {
	return func(t *Template) error {
		t.Echo = true
		return nil
	}
}

func (t *Template) Exec(args ...interface{}) error {
	db, err := pg.TryGetDB()
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")
	err = t.Content.Execute(buf, t.Data)
	if err != nil {
		return err
	}

	if t.Echo {
		fmt.Printf("executing template<Exec>: '%s'\n", t.Name)
	}

	stmt := db.StmtWithOpts(t.PgOpts)
	_, err = stmt.Exec(buf.String(), args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Template) Query(target interface{}, args ...interface{}) error {
	db, err := pg.TryGetDB()
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")
	err = t.Content.Execute(buf, t.Data)
	if err != nil {
		return err
	}

	if t.Echo {
		fmt.Printf("executing template<Query>: '%s'\n", t.Name)
	}

	stmt := db.StmtWithOpts(t.PgOpts)
	err = stmt.Query(target, buf.String(), args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Template) Get(target interface{}, args ...interface{}) error {
	db, err := pg.TryGetDB()
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")
	err = t.Content.Execute(buf, t.Data)
	if err != nil {
		return err
	}

	if t.Echo {
		fmt.Printf("executing template<Get>: '%s'\n", t.Name)
	}

	stmt := db.StmtWithOpts(t.PgOpts)
	err = stmt.Get(target, buf.String(), args...)
	if err != nil {
		return err
	}

	return nil
}
