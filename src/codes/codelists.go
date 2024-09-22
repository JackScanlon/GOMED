package codes

import (
	"io"
	"snomed/src/csv"
	"snomed/src/pg"

	"github.com/gocarina/gocsv"
)

type generatorHnd struct{}

type Concept struct {
	generatorHnd
	Id                 string  `csv:"id" dbType:"varchar" dbMod:"18"`
	EffectiveTime      pg.Date `csv:"effectiveTime" dbType:"date"`
	Active             bool    `csv:"active" dbType:"boolean"`
	ModuleId           string  `csv:"moduleId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	DefinitionStatusId string  `csv:"definitionStatusId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
}

type Description struct {
	generatorHnd
	Id                 string  `csv:"id" dbType:"varchar" dbMod:"18"`
	EffectiveTime      pg.Date `csv:"effectiveTime" dbType:"date"`
	Active             bool    `csv:"active" dbType:"boolean"`
	ModuleId           string  `csv:"moduleId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	ConceptId          string  `csv:"conceptId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	LanguageCode       string  `csv:"languageCode" dbType:"varchar" dbMod:"2"`
	TypeId             string  `csv:"typeId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	Term               string  `csv:"term" dbType:"varchar" dbMod:"256" dbReference:"Concept>id"`
	CaseSignificanceId string  `csv:"caseSignificanceId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
}

type Relationship struct {
	generatorHnd
	Id                   string  `csv:"id" dbType:"varchar" dbMod:"18"`
	EffectiveTime        pg.Date `csv:"effectiveTime" dbType:"date"`
	Active               bool    `csv:"active" dbType:"boolean"`
	ModuleId             string  `csv:"moduleId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	SourceId             string  `csv:"sourceId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	DestinationId        string  `csv:"destinationId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	RelationshipGroupId  string  `csv:"relationshipGroup" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	TypeId               string  `csv:"typeId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	CharacteristicTypeId string  `csv:"characteristicTypeId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	ModifierId           string  `csv:"modifierId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
}

type RefsetLang struct {
	generatorHnd
	Id                    pg.UUID `csv:"id" dbType:"uuid"`
	EffectiveTime         pg.Date `csv:"effectiveTime" dbType:"date"`
	Active                bool    `csv:"active" dbType:"boolean"`
	ModuleId              string  `csv:"moduleId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	RefsetId              string  `csv:"refsetId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
	ReferencedComponentId string  `csv:"referencedComponentId" dbType:"varchar" dbMod:"18" dbReference:"Description>id"`
	AcceptabilityId       string  `csv:"acceptabilityId" dbType:"varchar" dbMod:"18" dbReference:"Concept>id"`
}

func (p generatorHnd) Process(row any) (process bool, flat []any, err error) {
	flat, err = pg.FlattenRow(row)
	if err != nil {
		return false, flat, err
	}

	if flat[0] == "" {
		return false, flat, nil
	}

	return true, flat, nil
}

func (p generatorHnd) Reader() csv.ReaderFn {
	return func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.Comma = '\t'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		reader.TrimLeadingSpace = true
		return reader
	}
}
