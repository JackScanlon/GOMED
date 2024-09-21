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

## Reference
| Table          | Column                 | Value                | Summary                                                           |
|----------------|------------------------|----------------------|-------------------------------------------------------------------|
| `relationship` | `type_id`              | `116680003`          | Describes the parent-child (`IS_A`) relationship                  |
| `description`  | `type_id`              | `900000000000003001` | Fully Specified Name: Unique description of a Concept             |
| `description`  | `type_id`              | `900000000000013009` | Synonym: One or more additional, non-unique descriptor            |
| `description`  | `type_id`              | `900000000000550004` | Textual definition: narrative text (may exceed 256 char)          |
| `description`  | `case_significance_id` | `900000000000020002` | First character is case insensitive; rest is sensitive (`cl`)     |
| `description`  | `case_significance_id` | `900000000000017005` | Entire term case sensitive (`CS`)                                 |
| `description`  | `case_significance_id` | `900000000000448009` | Entire term case insensitive (`ci`)                               |
| `refset`       | `refset_id`            | `999001261000000100` | NHS RLRS: Clinical realm language reference set descriptor        |
| `refset`       | `refset_id`            | `999000691000001104` | NHS RLRS: Pharmacological realm language reference set descriptor |
| `refset`       | `acceptability_id`     | `900000000000548007` | Realm language reference set descriptor acceptability indicator   |
| `refset`       | `acceptability_id`     | `900000000000549004` | Realm language reference set descriptor acceptability indicator   |

## Processing

### 1. Collect SNOMED CT full release

1. Download release from:
    - URL @ [SNOMED CT UK Clinical Edition, RF2](https://isd.digital.nhs.uk/trud/users/authenticated/group/0/pack/26/subpack/101/releases)

2. Upload all `[Concept, Description, Relationship]` from `[release]/Full/Terminology`:
    - `SnomedCT_InternationalRF2_PRODUCTION_*/Full/Terminology`
    - `SnomedCT_UKClinicalRF2_PRODUCTION_*/Full/Terminology`
    - `SnomedCT_UKEditionRF2_PRODUCTION_*/Full/Terminology`

3. Upload ICD-10, OPCS4 and ~~simple~~ ReadCodeV3 code mapping from `SnomedCT_UKClinicalRF2_PRODUCTION_*/Full/Refset/Map`:
    - Filename: `der2_iisssciRefset_ExtendedMapUKCLFull_GB1000000_20240410.txt`
    - Filename: `der2_sRefset_SimpleMapMONOSnapshot_GB_20240828.txt`

### 2. Collect ReadCode map(s)

1. Download release from:
    - URL @ [NHS Data Migration](https://isd.digital.nhs.uk/trud/users/authenticated/filters/0/categories/9/items/9/releases)

2. Upload the following ReadCV2 maps:
    - Assured: `Mapping Tables/Updated/Clinically Assured/rcsctmap2_uk_20200401000001.txt`
    - Unassured: `Mapping Tables/Updated/Not Clinically Assured/Not Clinically Assured/rcsctmap_uk_20200401000001.txt`

3. Upload the following ReadCV3 maps:
    - Assured: `Mapping Tables/Updated/Clinically Assured/ctv3rctmap_uk_20200401000002.txt`

### 3. Generate simplified tables

1. Create `public.clinicalcode_snomed_codes` table
    - Shape:
        | Column           | Type                      | Description                           |
        |------------------|---------------------------|---------------------------------------|
        | `id`             | `integer [seq, pk]`       | Internal id                           |
        | `code`           | `varchar(18)`             | SNOMED Code (`SCTID`)                 |
        | `description`    | `varchar(256)`            | Fully specified name                  |
        | `case_sig`       | `sctsig [user-def, enum]` | Case significance of name             |
        | `active`         | `boolean`                 | Status of the code                    |
        | `effective_time` | `date`                    | When this code came into existence    |
        | `opcs4_codes`    | `text[]`                  | OPCS4 mapping                         |
        | `icd10_codes`    | `text[]`                  | ICD-10 mapping                        |
        | `readcv2_codes`  | `text[]`                  | ReadCV2 mapping                       |
        | `readcv3_codes`  | `text[]`                  | ReadCV3 mapping                       |
        | `search_vector`  | `tsvector [weighted]`     | Weighted search vector incl. synonyms |
        | `synonyms`       | `tsvector [bare, P&S]`    | One or more synonymous description(s) |

2. ...

### 4. Mapping coding systems to SNOMED CT

1. Generate mapping table(s)

    | Codelist(s)          | Filepath                                              | Step     |
    |----------------------|-------------------------------------------------------|----------|
    | ICD-10/OPCS4/ReadCV3 | `SnomedCT_UKClinicalRF2_PRODUCTION_*/Full/Refset/Map` | [Step 1](#1-collect-snomed-ct-full-release) |
    | ReadCV2              | _See `Step` column_                                   | [Step 2](#2-collect-readcode-map) |
    | ReadCV3              | _See `Step` column_                                   | [Step 2](#2-collect-readcode-map) |

2. Simplify map
    - Generate table composed of each map
    - Sieve duplicates
