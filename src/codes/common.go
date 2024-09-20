package codes

import (
	"snomed/src/csv"
	"snomed/src/trud"
)

const (
	ChunkSize   uint32 = 1000
	TableSchema string = "public"
	TablePrefix string = "clinicalcode"
	SnomedTag   string = "snomed"
)

type ReleaseItem struct {
	Name      string
	Dir       string
	Fmt       string
	Model     any
	Category  trud.Category
	Filenames []string
}

type Reader interface {
	Reader() csv.ReaderFn
}

type Processor interface {
	Process(any) (bool, []any, error)
}

var (
	SnomedReleaseFiles = [...]string{
		"SnomedCT_UKEditionRF2_PRODUCTION_*",
		"SnomedCT_UKClinicalRF2_PRODUCTION_*",
		"SnomedCT_InternationalRF2_PRODUCTION_*",
	}

	SnomedReleaseGroups = [...]ReleaseItem{
		{
			Name:      "Concept",
			Dir:       "Snapshot/Terminology",
			Fmt:       "sct2_%s_*.txt",
			Model:     new(Concept),
			Category:  trud.SNOMED_RELEASE,
			Filenames: []string{"Concept"},
		},
		{
			Name:      "Description",
			Dir:       "Snapshot/Terminology",
			Fmt:       "sct2_%s_*.txt",
			Model:     new(Description),
			Category:  trud.SNOMED_RELEASE,
			Filenames: []string{"Description"},
		},
		{
			Name:      "Relationship",
			Dir:       "Snapshot/Terminology",
			Fmt:       "sct2_%s_*.txt",
			Model:     new(Relationship),
			Category:  trud.SNOMED_RELEASE,
			Filenames: []string{"Relationship", "StatedRelationship"},
		},
		{
			Name:      "Refset",
			Dir:       "Snapshot/Refset/Map",
			Fmt:       "der2_%s_*.txt",
			Model:     new(RefsetMap),
			Category:  trud.SNOMED_RELEASE,
			Filenames: []string{"iisssciRefset", "sRefset"},
		},
		{
			Name:      "CtvMap",
			Dir:       "Mapping Tables/Updated/*Assured",
			Fmt:       "%s_*.txt",
			Model:     new(CtvMap),
			Category:  trud.SNOMED_READ_MAP,
			Filenames: []string{"ctv3sctmap2_uk", "rcsctmap2_uk", "rcsctmap_enhanced_uk"},
		},
	}

	RefsetIdName = map[string]string{
		// OPSC4
		"1126441000000105": "OPCS4",
		"1382401000000109": "OPCS4",

		// Read
		"999002721000000108": "ReadCodeV2",
		"999002731000000105": "ReadCodeV3",
		"900000000000497000": "ReadCodeV3",

		// ICD
		"446608001": "ICD-10",
		"447562003": "ICD-10",

		"999002271000000101": "ICD-10",
		"999001921000000102": "ICD-10",
		"999001331000000104": "ICD-10",
	}
)
