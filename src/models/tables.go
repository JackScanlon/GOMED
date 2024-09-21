package models

import (
	"snomed/src/pg"
)

type SctIdentifier string

const (
	SctFullySpecified SctIdentifier = "F"
	SctPreferred      SctIdentifier = "P"
	SctSynonym        SctIdentifier = "S"
	SctTextual        SctIdentifier = "D"
)

type SctCase string

const (
	SctCL SctIdentifier = "CL"
	SctCS SctIdentifier = "CS"
	SctCI SctIdentifier = "CI"
)

type SnomedCodes struct {
	Id            int      `dbType:"integer" dbIsPrimary:"true"`
	Code          string   `dbType:"VARCHAR" dbMod:"18"`
	Description   string   `dbType:"VARCHAR" dbMod:"256"`
	CaseSig       SctCase  `dbType:"sctcase"`
	Active        bool     `dbType:"BOOLEAN"`
	EffectiveTime pg.Date  `dbType:"DATE"`
	Opcs4Codes    []string `dbType:"text[]"`
	Icd10Codes    []string `dbType:"text[]"`
	Readcv2Codes  []string `dbType:"text[]"`
	Readcv3Codes  []string `dbType:"text[]"`
	// SearchVector tsvector
	// Synonyms     tsvector
}
