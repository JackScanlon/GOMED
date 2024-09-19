package codes

type TableMap struct {
	Name      string
	Model     any
	Filenames []string
}

var (
	CodelistFileDir  = "Snapshot/Terminology"
	CodelistFileFmt  = "sct2_%s_*.txt"
	CodelistReleases = [...]string{
		"SnomedCT_UKEditionRF2_PRODUCTION_*",
		"SnomedCT_UKClinicalRF2_PRODUCTION_*",
		"SnomedCT_InternationalRF2_PRODUCTION_*",
	}

	CodelistSchema = "public"
	CodelistParent = "clinicalcode"
	CodelistPrefix = "snomed"

	TableMappings = []TableMap{
		{
			Name:      "Concept",
			Model:     new(Concept),
			Filenames: []string{"Concept"},
		},
		{
			Name:      "Description",
			Model:     new(Description),
			Filenames: []string{"Description"},
		},
		{
			Name:      "Relationship",
			Model:     new(Relationship),
			Filenames: []string{"Relationship", "StatedRelationship"},
		},
	}
)
