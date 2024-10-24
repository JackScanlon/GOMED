# ICD-9-CM

## 1. Overview

### 1.1. Hierarchy

```text
  ┌──────────────────────┐
  │                      │
  │  Root ICD-9-CM Node  │
  │                      │
  └───────────┬──────────┘
              │
              └──────────────────┐
                                 │
                 ┌───────────────▼──────────────┐
                 │                              │
                 │ "C" -> Cardiovascular System │
                 │                              │
                 └───────────────┬──────────────┘
                                 │
                                 │
                        ┌────────▼───────────┐
                        │                    │
                        │ "C03" -> Diuretics │
                        │                    │
                        └────────┬───────────┘
                                 │
                                 └───────────────┐
                                                 │
                              ┌──────────────────▼───────────────┐
                              │                                  │
                              │ "C03C" -> High-ceiling diuretics │
                              │                                  │
                              └──────────────────┬───────────────┘
                                                 │
                                                 │
                                                 │
                                 ┌───────────────▼───────────────────────┐
                                 │                                       │
                                 │ "C03CA" -> High-ceiling diuretics and │
                                 │            potassium-sparing agents   │
                                 │            combination                │
                                 │                                       │
                                 └───────────────────────────────────────┘
```

### 1.2. Codelist dataset

| Column | Summary          |
|--------|------------------|
| `$1`   | ICD-9 code       |
| `$2`   | Code description |

### 1.3. Mapping dataset

| Column       | Type           | Summary            |
|--------------|----------------|--------------------|
| `ICD_CODE`   | `varchar(7)`   | ICD-9 code         |
| `ICD_NAME`   | `varchar(255)` | Code description   |
| `SNOMED_CID` | `varchar(18)`  | Snomed mapped code |

## 2. Availability

See:
- Source: https://www.cms.gov/medicare/coding-billing/icd-10-codes/icd-9-cm-diagnosis-procedure-codes-abbreviated-and-full-code-titles
- Download: https://www.cms.gov/medicare/coding/icd9providerdiagnosticcodes/downloads/icd-9-cm-v32-master-descriptions.zip

## 3. Processing

### 3.1. Generate ICD-9-CM codelist

1. Download release from:
    - URL @ [Version 32 Full and Abbreviated Code Titles - Effective October 1, 2014 (ZIP)](https://www.cms.gov/medicare/coding/icd9providerdiagnosticcodes/downloads/icd-9-cm-v32-master-descriptions.zip)

2. Upload all `[Diagnosis Code, Product Code]` from...
    - Diagnosis codes: `CMS32_DESC_LONG_DX.txt`
    - Product codes: `CMS32_DESC_LONG_SG.txt`

### 3.2. Object Ref

Internal codelist object reference:

| Column         | Type          |
|----------------|---------------|
| `code`         | `varchar(50)` |
| `description`  | `varchar(255)`|
| `import_date`  | `timestamp`   |
| `created_date` | `timestamp`   |
