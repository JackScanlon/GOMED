package codes

import (
	"io"

	"snomed/src/csv"
	"snomed/src/pg"

	"github.com/gocarina/gocsv"
)

type readerHnd struct{}

func (ref readerHnd) Reader() csv.ReaderFn {
	return func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.Comma = '\t'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		return reader
	}
}

// SNOMED
//
//	iissciRefset: `id      effectiveTime   active  moduleId        refsetId        referencedComponentId   mapGroup        mapPriority     mapRule mapAdvice       mapTarget       correlationId   mapBlock`
//	     sRefset: `id      effectiveTime   active  moduleId        refsetId        referencedComponentId   mapTarget`
type RefsetMap struct {
	readerHnd
	Id                    pg.UUID `csv:"id" dbType:"uuid"`
	EffectiveTime         pg.Date `csv:"effectiveTime" dbType:"DATE"`
	Active                bool    `csv:"active" dbType:"BOOLEAN"`
	ModuleId              string  `csv:"moduleId" dbType:"VARCHAR" dbMod:"18" dbReference:"Concept>id"`
	RefsetId              string  `csv:"refsetId" dbType:"VARCHAR" dbMod:"18" dbReference:"Concept>id"`
	ReferencedComponentId string  `csv:"referencedComponentId" dbType:"VARCHAR" dbMod:"18" dbReference:"Concept>id"`
	MapGroup              int     `csv:"mapGroup,default=1" dbType:"INTEGER"`
	MapPriority           int     `csv:"mapPriority,default=1" dbType:"INTEGER"`
	MapTarget             string  `csv:"mapTarget" dbType:"VARCHAR" dbMod:"7"`
	MapOrigin             string  `csv:"mapOrigin,omitempty" dbType:"VARCHAR" dbMod:"12"`
}

func (ref RefsetMap) Process(row any) (process bool, flat []any, err error) {
	if len(ref.MapTarget) < 1 || ref.MapTarget[0] == '#' {
		return false, flat, err
	}

	var origin string
	if mapOrigin, ok := RefsetIdName[ref.RefsetId]; ok {
		origin = mapOrigin
	} else {
		return false, flat, err
	}

	flat, err = pg.FlattenRow(row)
	if err != nil {
		return false, flat, err
	}
	flat[9] = origin

	return true, flat, nil
}

// CTV2 | DATA MIGRATION/Assured
// 	rcsctmap2: `MapId   ReadCode        TermCode        ConceptId       DescriptionId   IS_ASSURED      EffectiveDate   MapStatus`

// CTV2 | DATA MIGRATION/Unassured
// 	rcsctmap_enhanced: `MapId   ReadCode        TermCode        ConceptId       T30ID   T60ID   T198ID  EffectiveDate   MapStatus`

// CTV3 | DATA MIGRATION/Assured
//
//	ctv3sctmap2: `MAPID   CTV3_CONCEPTID  CTV3_TERMID     CTV3_TERMTYPE   SCT_CONCEPTID   SCT_DESCRIPTIONID       MAPSTATUS       EFFECTIVEDATE   IS_ASSURED`
type CtvMap struct {
	readerHnd
	Id            pg.UUID `csv:"id,MapId,MAPID" dbType:"uuid"`
	EffectiveTime pg.Date `csv:"effectiveTime,effectiveDate,EffectiveDate,EFFECTIVEDATE" dbType:"DATE"`
	Active        bool    `csv:"active,MapStatus,MAPSTATUS" dbType:"BOOLEAN"`
	ConceptId     string  `csv:"conceptId,ConceptId,SCT_CONCEPTID" dbType:"VARCHAR" dbMod:"18" dbReference:"Concept>id"`
	DescriptionId string  `csv:"descriptionId,DescriptionId,SCT_DESCRIPTIONID" dbType:"VARCHAR" dbMod:"18" dbReference:"Description>id"`
	IsAssured     bool    `csv:"IS_ASSURED" dbType:"BOOLEAN"`
	MapTarget     string  `csv:"readCode,ReadCode,CTV3_TERMID" dbType:"VARCHAR" dbMod:"7"`
	MapOrigin     string  `csv:"mapOrigin,CTV3_TERMTYPE,omitempty" dbType:"VARCHAR" dbMod:"12"`
}

func (ref CtvMap) Process(row any) (process bool, flat []any, err error) {
	if len(ref.MapTarget) < 1 || ref.MapTarget[0] == '#' {
		return false, flat, err
	}

	flat, err = pg.FlattenRow(row)
	if err != nil {
		return false, flat, err
	}

	if len(ref.MapOrigin) < 1 {
		flat[7] = "ReadCodeV2"
	} else {
		flat[7] = "ReadCodeV3"
	}

	return true, flat, nil
}
