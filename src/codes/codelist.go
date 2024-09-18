package codes

import (
	"fmt"
	"path"
	"path/filepath"

	"snomed/src/trud"
	// "github.com/gocarina/gocsv"
)

var (
	codelistFileDir  = "Full/Terminology"
	codelistFileFmt  = "sct2_%s_*.txt"
	codelistReleases = [...]string{
		"SnomedCT_UKEditionRF2_PRODUCTION_*",
		"SnomedCT_UKClinicalRF2_PRODUCTION_*",
		"SnomedCT_InternationalRF2_PRODUCTION_*",
	}

	codelistPrefix = "clinicalcode"
	codelistTables = map[string]string{
		"Concept":            "snomed_concept",
		"Description":        "snomed_description",
		"StatedRelationship": "snomed_relationship",
		"Relationship":       "snomed_relationship",
	}
)

type Concept struct {
	Id                 string `csv:"id"`
	EffectiveTime      string `csv:"effectiveTime"`
	Active             bool   `csv:"active"`
	ModuleId           string `csv:"moduleId"`
	DefinitionStatusId string `csv:"definitionStatusId"`
}

type Description struct {
	Id               string `csv:"id"`
	EffectiveTime    string `csv:"effectiveTime"`
	Active           bool   `csv:"active"`
	ModuleId         string `csv:"moduleId"`
	ConceptId        string `csv:"conceptId"`
	LanguageCode     string `csv:"languageCode"`
	TypeId           string `csv:"typeId"`
	Term             string `csv:"term"`
	CaseSignificance string `csv:"caseSignificanceId"`
}

type Relationship struct {
	Id                   string `csv:"id"`
	EffectiveTime        string `csv:"effectiveTime"`
	Active               bool   `csv:"active"`
	ModuleId             string `csv:"moduleId"`
	SourceId             string `csv:"sourceId"`
	DestinationId        string `csv:"destinationId"`
	RelationshipGroup    string `csv:"relationshipGroup"`
	TypeId               string `csv:"typeId"`
	CharacteristicTypeId string `csv:"characteristicTypeId"`
	ModifierId           string `csv:"modifierId"`
}

func createFromSource(releaseName string, tableName string, dir string) error {
	pattern := path.Join(dir, codelistFileDir, fmt.Sprintf(codelistFileFmt, releaseName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	// TODO: async chan to batch table insert?
	for _, file := range matches {
		fmt.Printf("[%s: %s] %s\n", releaseName, tableName, file)
	}

	return nil
}

func TryCreateTables(release *trud.Release, dir string) error {
	/*
		TODO:
			- Impl. snomed codelist generator
			- Handle code mapping table generation
	*/
	if !trud.IsCategory(release, trud.SNOMED_RELEASE) {
		return nil
	}

	for _, releaseDir := range codelistReleases {
		releaseDir := path.Join(dir, release.Metadata.Name, releaseDir)
		for releaseName, tableName := range codelistTables {
			tableName := fmt.Sprintf("%s_%s", codelistPrefix, tableName)

			// TODO: det. whether table exists and is in the correct shape

			if err := createFromSource(releaseName, tableName, releaseDir); err != nil {
				return err
			}
		}
	}

	return nil
}
