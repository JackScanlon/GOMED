package models

type SctIdentifier string

const (
	SctFullySpecified SctIdentifier = "F"
	SctPreferred      SctIdentifier = "P"
	SctSynonym        SctIdentifier = "S"
	SctTextual        SctIdentifier = "D"
)

type SctCase int

const (
	SctCL SctCase = 0
	SctCS SctCase = 1
	SctCI SctCase = 2
)

type OntologyType int

const (
	CLINICAL_DISEASE            OntologyType = 0
	CLINICAL_DOMAIN             OntologyType = 1
	CLINICAL_FUNCTIONAL_ANATOMY OntologyType = 2
)
