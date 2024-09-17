# SNOMED

## Overview

```text

 ┌───────────────────────────────────┐                                 ┌───────────────────────────────┐
 │ SCT_Description                   │                                 │                               │
 │                                   │                                 │ SCT_Relationship              │
 │ `typeid` defines whether this     ┼──────────────────┐       ┌──────►                               │
 │ element is a fully specified name │                  │       │      │ Describes attribute & concept │
 │ or if it's a preferred synonym    │                  │       │      │ relationships                 │
 └────────────────────────▲──────────┘                  │       │      │                               │
                          │                             │       │      └───────────────┬───────────────┘
                          │                             │       │                      │
                          │          ┌──────────────────▼───────┼─────────┐            │
                          │          │ SCT_Concept                        │            │
                          │          │                                    │            │
                          │          │ The numerical codes that identify  ◄────────────┘
                          │          │ clinical terms, attributes +/-     │
                          │          │ primitives etc                     │
                          │          └─────▲───┬──────────────────────────┘
 ┌────────────────────────┴───────┐        │   │
 │ SCT_ReferenceSets              │        │   │
 │                                ┼────────┘   │
 │ Groups Concepts & Descriptions │            │
 │ into sets                      ◄────────────┘
 └────────────────────────────────┘

```

## Processing

### 1. Collect SNOMED CT full release

1. Download release from:
    - URL @ [SNOMED CT UK Clinical Edition, RF2](https://isd.digital.nhs.uk/trud/users/authenticated/group/0/pack/26/subpack/101/releases)

2. Upload all `[Concept, Description, Relationship]` from `[release]/Full/Terminology`:
    - `SnomedCT_InternationalRF2_PRODUCTION_*`
    - `SnomedCT_UKClinicalRF2_PRODUCTION_*`
    - `SnomedCT_UKEditionRF2_PRODUCTION_*`

3. Upload ICD-10 mapping from `SnomedCT_UKClinicalRF2_PRODUCTION_*/Full/Refset/Map`:
    - Filename: `der2_iisssciRefset_ExtendedMapUKCLFull_GB1000000_20240410.txt`

### 2. Collect ReadCode map

1. Download release from:
    - URL @ [NHS Data Migration](https://isd.digital.nhs.uk/trud/users/authenticated/filters/0/categories/9/items/9/releases)

2. Upload the following ReadCode CV2 maps:
    - Primary care: `Primary Care Refsets/der2_cssRefset_PrimaryCareFull_GB1000000_20200401.txt`
    - Assured: `Mapping Tables/Updated/Clinically Assured/rcsctmap2_uk_20200401000001.txt`
    - Unassured: `Mapping Tables/Updated/Not Clinically Assured/Not Clinically Assured/rcsctmap_uk_20200401000001.txt`

3. Upload the following ReadCode CV3 maps:
    - Assured: `Mapping Tables/Updated/Clinically Assured/ctv3rctmap_uk_20200401000002.txt`


### 3. Mapping coding systems to SNOMED CT

1. Map ICD-10 4th edition to SNOMED CT:
    - Use `SnomedCT_UKClinicalRF2_PRODUCTION_*/Full/Refset/Map` from [Step 1](#1-collect-snomed-ct-full-release)

2. Map ReadCode CV2 to SNOMED CT:
    - Use files found in part 2 of [Step 2](#2-collect-readcode-map)

3. Map ReadCode CV3 to SNOMED CT:
    - Use files found in part 3 of [Step 2](#2-collect-readcode-map)
