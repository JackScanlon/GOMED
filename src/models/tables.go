package models

import (
	"snomed/src/pg"
)

type SnomedCodes struct {
	Id            int      `dbType:"integer" dbIsPrimary:"true"`
	Code          string   `dbType:"varchar" dbMod:"18"`
	Description   string   `dbType:"varchar" dbMod:"256"`
	CaseSig       SctCase  `dbType:"integer"`
	Active        bool     `dbType:"boolean"`
	EffectiveTime pg.Date  `dbType:"date"`
	Opcs4Codes    []string `dbType:"text[]"`
	Icd10Codes    []string `dbType:"text[]"`
	Readcv2Codes  []string `dbType:"text[]"`
	Readcv3Codes  []string `dbType:"text[]"`
	// SearchVector tsvector
	// Synonyms     tsvector
}

type OntologyProps struct {
	Code           string `json:"code"`
	CodingSystemId string `json:"coding_system_id"`
}

type OntologyTag struct {
	Id          int64         `dbType:"bigint" dbIsPrimary:"true"`
	Name        string        `dbType:"varchar" dbMod:"256"`
	TypeId      OntologyType  `dbType:"integer"`
	ReferenceId int64         `dbType:"bigint"`
	Properties  OntologyProps `dbType:"jsonb"`
	// SearchVector tsvector
}

type OntologyTagEdge struct {
	Id       int64 `dbType:"bigint" dbIsPrimary:"true"`
	ChildId  int64 `dbType:"bigint" dbReference:"OntologyTag>id"`
	ParentId int64 `dbType:"bigint" dbReference:"OntologyTag>id"`
}
