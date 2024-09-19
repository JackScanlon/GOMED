package codes

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Date struct {
	pgtype.Date
}

func (date *Date) MarshalCSV() (string, error) {
	return date.Time.Format("20060102"), nil
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("20060102", csv)
	date.Valid = true

	return err
}

type Concept struct {
	Id                 string `csv:"id" dbType:"VARCHAR" dbMod:"18"`
	EffectiveTime      Date   `csv:"effectiveTime" dbType:"DATE"`
	Active             bool   `csv:"active" dbType:"BOOLEAN"`
	ModuleId           string `csv:"moduleId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	DefinitionStatusId string `csv:"definitionStatusId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
}

type Description struct {
	Id                 string `csv:"id" dbType:"VARCHAR" dbMod:"18"`
	EffectiveTime      Date   `csv:"effectiveTime" dbType:"DATE"`
	Active             bool   `csv:"active" dbType:"BOOLEAN"`
	ModuleId           string `csv:"moduleId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	ConceptId          string `csv:"conceptId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	LanguageCode       string `csv:"languageCode" dbType:"VARCHAR" dbMod:"2"`
	TypeId             string `csv:"typeId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	Term               string `csv:"term" dbType:"VARCHAR" dbMod:"256" dbReference:"public.Concept>id"`
	CaseSignificanceId string `csv:"caseSignificanceId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
}

type Relationship struct {
	Id                   string `csv:"id" dbType:"VARCHAR" dbMod:"18"`
	EffectiveTime        Date   `csv:"effectiveTime" dbType:"DATE"`
	Active               bool   `csv:"active" dbType:"BOOLEAN"`
	ModuleId             string `csv:"moduleId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	SourceId             string `csv:"sourceId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	DestinationId        string `csv:"destinationId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	RelationshipGroupId  string `csv:"relationshipGroup" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	TypeId               string `csv:"typeId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	CharacteristicTypeId string `csv:"characteristicTypeId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
	ModifierId           string `csv:"modifierId" dbType:"VARCHAR" dbMod:"18" dbReference:"public.Concept>id"`
}
