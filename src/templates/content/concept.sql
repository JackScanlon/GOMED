/********************************************************************************
 * Constants                                                                    *
 *                                                                              *
 *  ? Const. values describing identifiers etc                                  *
 *                                                                              *
 *  [!] Ref: https://nhsengland.kahootz.com/t_c_home/view?objectId=45960752     *
 *                                                                              *
 ********************************************************************************/

/*md

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

*/



/********************************************************************************
 * Description identifier                                                       *
 *                                                                              *
 *  ? Compute the descriptor type and assign a identifier                       *
 *                                                                              *
 ********************************************************************************/

--![template] name:"descriptionIdentifier"

do $tx$
declare
  FSN        constant varchar   := '900000000000003001';
  SYN        constant varchar   := '900000000000013009';
  DEF        constant varchar   := '900000000000550004';

  RLRS       constant varchar[] := '{"999001261000000100", "999000691000001104"}'::varchar[];

  PREFERRED  constant varchar   := '900000000000548007';
  ACCEPTABLE constant varchar   := '900000000000549004';
begin
  -- create description identifier
  if not exists(select 1 from pg_catalog.pg_type where typname = 'sctident') then
    create type sctident as enum (
      -- fully specified name
      'F',
      -- preferred synonym
      'P',
      -- synonym
      'S',
      -- textual definition
      'D'
    );
  end if;

  -- drop legacy column
  if exists(
    select 1
      from information_schema.columns
     where table_name = 'clinicalcode_snomed_description'
       and column_name = 'identifier'
  ) then
    alter table public.clinicalcode_snomed_description
     drop column identifier;
  end if;

  -- append term identifier column
  alter table public.clinicalcode_snomed_description
    add identifier sctident default null;

  -- compute the identifier and assign to each
  update public.clinicalcode_snomed_description as d
     set identifier = (
       case
         when r.acceptability_id = PREFERRED  and d.type_id = FSN then 'F'::sctident
         when r.acceptability_id = PREFERRED  and d.type_id = SYN then 'P'::sctident
         when r.acceptability_id = ACCEPTABLE and d.type_id = SYN then 'S'::sctident
         when d.type_id = DEF then 'D'::sctident
         else null::sctident
       end
     )
    from public.clinicalcode_snomed_refset_lang as r
   where d.id = r.referenced_component_id
     and r.active = true
     and r.refset_id = any(RLRS::varchar[]);
end;
$tx$ language plpgsql;

--![endtemplate]



/********************************************************************************
 * Simplified codelist                                                          *
 *                                                                              *
 *  ? Generate a single, simplified table containing human readable terms and   *
 *    associated descriptors                                                    *
 *                                                                              *
 ********************************************************************************/

--![template] name:"simplifyCodelist"

do $tx$
declare
  SIG_CL  constant varchar := '900000000000020002';
  SIG_CS  constant varchar := '900000000000017005';
  SIG_CI  constant varchar := '900000000000448009';

  ENUM_CL constant integer := 0;
  ENUM_CI constant integer := 1;
  ENUM_CS constant integer := 2;
begin
  -- install ext(s) if not available
  create extension if not exists pg_trgm schema public;
  create extension if not exists btree_gin schema public;

  -- create ts agg if not available
  if not exists(select 1 from pg_catalog.pg_proc where proname = 'tsvector_agg' and prokind = 'a') then
    create aggregate tsvector_agg(tsvector) (
      stype = pg_catalog.tsvector,
      sfunc = pg_catalog.tsvector_concat,
      initcond = ''
    );
  end if;

  --[?] Handle table migration if managed service
  {{if .Data.managed}}
    -- drop table if exists
    if exists(
      select 1 from information_schema.tables where table_schema='public' and table_name='clinicalcode_snomed_codes'
    ) then
      drop table public.clinicalcode_snomed_codes;
    end if;

    -- create snomed table
    create table public.clinicalcode_snomed_codes (
      id              serial        primary key,
      code            varchar(18)   not null default '',
      description     varchar(256)  not null default '',
      case_sig        integer       default ENUM_CI,
      active          boolean       default true,
      effective_time  date          not null default now()::date,
      mesh_codes      text[]        default '{}'::text[],
      opcs4_codes     text[]        default '{}'::text[],
      icd9_codes      text[]        default '{}'::text[],
      icd10_codes     text[]        default '{}'::text[],
      readcv2_codes   text[]        default '{}'::text[],
      readcv3_codes   text[]        default '{}'::text[],
      search_vector   tsvector      default '',
      synonyms_vector tsvector      default '',
      relation_vector tsvector      default '',
      unique (code)
    );
  {{else}}
    -- truncate (cascade not enforced to ensure an error is thrown)
    truncate table public.clinicalcode_snomed_codes;
  {{end}}

  -- insert `clinicalcode_snomed_codes` rows
  with
    --> collect valid concepts
    concepts as (
      select
            row_number() over (order by c.id::bigint) as id,
            c.id as code,
            d.term as description,
            (
              case
                when d.case_significance_id = SIG_CL then ENUM_CL
                when d.case_significance_id = SIG_CS then ENUM_CS
                when d.case_significance_id = SIG_CI then ENUM_CI
                else ENUM_CI
              end
            ) as case_sig,
            c.active,
            c.effective_time
        from public.clinicalcode_snomed_concept as c
        join public.clinicalcode_snomed_description as d
          on c.id = d.concept_id
       where c.active = true
         and d.active = true
         and d.identifier = 'F'::sctident
    ),
    --> collect the synonyms of each valid concept
    synonyms as (
      select
          c.code,
          tsvector_agg(
            to_tsvector('pg_catalog.english', coalesce(d.term, '')::text)
          ) as vec
        from concepts as c
        join public.clinicalcode_snomed_description as d
          on c.code = d.concept_id
       where d.active = true
         and d.identifier in ('P'::sctident, 'S'::sctident)
       group by c.code
    ),
    --> collect all code mappings from: OPCS4/ICD-10/ReadCV2/ReadCV3
    map_components as (
      select
          rank() over (partition by snomed_code, source, code order by priority) as nrow,
          *
        from (
          select
                c.code as snomed_code,
                r.map_origin as source,
                r.map_target as code,
                r.map_priority as priority
            from concepts as c
            join public.clinicalcode_snomed_refset_map as r
              on c.code = r.referenced_component_id
           where r.active = true
          union
          select
                c.code as snomed_code,
                r.map_origin as source,
                r.map_target as code,
                (
                  case
                    when r.is_assured then 1
                    else 2
                  end
                ) as priority
            from concepts as c
            join public.clinicalcode_snomed_ctv_map as r
              on c.code = r.concept_id
           where r.active = true
      ) as m
    ),
    --> deduplicate code map(s) and aggregate
    codemap as (
      select
            snomed_code,
            source,
            array_agg(code order by priority asc)::text[] as codes
        from map_components
       where nrow = 1
       group by snomed_code, source
    )
  insert into public.clinicalcode_snomed_codes (
              id,             code,      description,
        case_sig,           active,   effective_time,
      mesh_codes,      opcs4_codes,       icd9_codes,
     icd10_codes,    readcv2_codes,    readcv3_codes,
    search_vector,  synonyms_vector, relation_vector
  )
    select
          c.id,
          c.code,
          c.description,
          c.case_sig,
          c.active,
          c.effective_time,
          mesh.codes as mesh_codes,
          opcs.codes as opcs4_codes,
          icd9.codes as icd9_codes,
          icd10.codes as icd10_codes,
          cv2.codes as readcv2_codes,
          cv3.codes as readcv3_codes,
          (
            setweight(to_tsvector('pg_catalog.english', coalesce(c.code, '')), 'A') ||
            setweight(to_tsvector('pg_catalog.english', coalesce(c.description, '')), 'A') ||
            setweight(s.vec, 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(mesh.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(opcs.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(icd9.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(icd10.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(cv2.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(cv3.codes, '{}'::text[]), ' ')), 'B')
          ) as search_vector,
          s.vec as synonyms_vector,
          (
            setweight(to_tsvector('pg_catalog.english', coalesce(c.code, '')), 'A') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(mesh.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(opcs.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(icd9.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(icd10.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(cv2.codes, '{}'::text[]), ' ')), 'B') ||
            setweight(to_tsvector('pg_catalog.english', array_to_string(coalesce(cv3.codes, '{}'::text[]), ' ')), 'B')
          ) as relation_vector
      from concepts as c
      left join synonyms as s
        using (code)
      left join public.temp_meshmap as mesh
        on c.code = mesh.snomed_code
      left join codemap as opcs
        on c.code = opcs.snomed_code and opcs.source = 'OPCS4'
      left join public.temp_icd9map as icd9
        on c.code = icd9.snomed_code
      left join codemap as icd10
        on c.code = icd10.snomed_code and icd10.source = 'ICD-10'
      left join codemap as cv2
        on c.code = cv2.snomed_code and cv2.source = 'ReadCodeV2'
      left join codemap as cv3
        on c.code = cv3.snomed_code and cv3.source = 'ReadCodeV3'
    on conflict (code)
    do update
      set
        code            = excluded.code,
        description     = excluded.description,
        case_sig        = excluded.case_sig,
        active          = excluded.active,
        mesh_codes      = excluded.mesh_codes,
        opcs4_codes     = excluded.opcs4_codes,
        icd9_codes      = excluded.icd9_codes,
        icd10_codes     = excluded.icd10_codes,
        readcv2_codes   = excluded.readcv2_codes,
        readcv3_codes   = excluded.readcv3_codes,
        search_vector   = excluded.search_vector,
        synonyms_vector = excluded.synonyms_vector,
        relation_vector = excluded.relation_vector
     where excluded.active
       and excluded.effective_time > clinicalcode_snomed_codes.effective_time;

  --[?] Handle index creation if managed
  {{if .Data.managed}}
    -- create index
    create index sct_cd_trgm_idx   on public.clinicalcode_snomed_codes using gin (code          gin_trgm_ops);
    create index sct_desc_trgm_idx on public.clinicalcode_snomed_codes using gin (description   gin_trgm_ops);

    create index sct_mesh_txt_idx  on public.clinicalcode_snomed_codes using gin (mesh_codes    array_ops);
    create index sct_opcs_txt_idx  on public.clinicalcode_snomed_codes using gin (opcs4_codes   array_ops);
    create index sct_icd9_txt_idx  on public.clinicalcode_snomed_codes using gin (icd9_codes    array_ops);
    create index sct_icd10_txt_idx on public.clinicalcode_snomed_codes using gin (icd10_codes   array_ops);
    create index sct_cv2_txt_idx   on public.clinicalcode_snomed_codes using gin (readcv2_codes array_ops);
    create index sct_cv3_txt_idx   on public.clinicalcode_snomed_codes using gin (readcv3_codes array_ops);

    -- create sv index
    create index sct_sv_gin_idx    on public.clinicalcode_snomed_codes using gin (search_vector  );
    create index sct_syn_gin_idx   on public.clinicalcode_snomed_codes using gin (synonyms_vector);
    create index sct_rel_gin_idx   on public.clinicalcode_snomed_codes using gin (relation_vector);
  {{end}}

end;
$tx$ language plpgsql;

--![endtemplate]
