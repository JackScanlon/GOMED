package trud

import (
	"strconv"
	"strings"
)

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
		SNOMED_RELEASE:  1799, // 101
		SNOMED_READ_MAP: 9,
	}
	categoryNames = map[string]Category{
		"SNOMED_NONE":     SNOMED_NONE,
		"SNOMED_RELEASE":  SNOMED_RELEASE,
		"SNOMED_READ_MAP": SNOMED_READ_MAP,
		"SNOMED_ALL":      SNOMED_ALL,
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

func IsCategoryId(catId uint16, cat Category) bool {
	switch cat {
	case SNOMED_ALL:
		for _, id := range categoryIds {
			if id == catId {
				return true
			}
		}
		return false
	default:
		break
	}

	id, ok := categoryIds[cat]
	if !ok {
		return false
	}

	return catId == id
}

func ParseCategory(str string) (bool, Category) {
	cat, ok := categoryNames[strings.ToUpper(str)]
	if ok {
		return true, cat
	}

	val, err := strconv.ParseUint(str, 10, 8)
	if err == nil {
		return true, Category(val)
	}

	return false, SNOMED_NONE
}
