package pg

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Date struct {
	pgtype.Date
}

func (d *Date) MarshalCSV() (string, error) {
	return d.Time.Format("20060102"), nil
}

func (d *Date) UnmarshalCSV(csv string) (err error) {
	d.Time, err = time.Parse("20060102", csv)
	d.Valid = true

	return err
}

type UUID struct {
	pgtype.UUID
}

func (u *UUID) MarshalCSV() (string, error) {
	var str string
	id, err := uuid.ParseBytes(u.Bytes[:])
	if err != nil {
		return str, nil
	}

	return id.String(), nil
}

func (u *UUID) UnmarshalCSV(csv string) (err error) {
	id, err := uuid.Parse(csv)
	if err != nil {
		return nil
	}

	u.Bytes = id
	u.Valid = true
	return nil
}
