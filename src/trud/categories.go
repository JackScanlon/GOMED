package trud

type Category uint8

const (
	SNOMED_NONE     Category = 0x1 << 0
	SNOMED_RELEASE  Category = 0x1 << 1
	SNOMED_READ_MAP Category = 0x1 << 2
	SNOMED_ALL      Category = SNOMED_RELEASE | SNOMED_READ_MAP
)

var (
	categoryUrl = "https://isd.digital.nhs.uk/trud/api/v1/keys/%s/items/%d/releases?latest"
	categoryIds = map[Category]uint16{
		SNOMED_RELEASE:  101,
		SNOMED_READ_MAP: 9,
	}
)

func (category Category) Has(flag Category) bool {
	return (category & flag) == flag
}

func (category Category) GetIds() []uint16 {
	var categoryId []uint16
	for i := uint8(1); i < uint8(3); i++ {
		comp := Category(0x1 << i)
		if category.Has(comp) {
			categoryId = append(categoryId, categoryIds[comp])
		}
	}

	return categoryId
}
