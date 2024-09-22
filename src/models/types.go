package models

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

type OntologyType int

const (
	CLINICAL_DISEASE            = 0
	CLINICAL_DOMAIN             = 1
	CLINICAL_FUNCTIONAL_ANATOMY = 2
)
