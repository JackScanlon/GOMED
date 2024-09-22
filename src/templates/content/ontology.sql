/********************************************************************************
 * Ontology                                                                     *
 *                                                                              *
 *  ? Build ontology tags & their relationship(s)                               *
 *                                                                              *
 ********************************************************************************/

--![template] name:"network"

do $tx$
declare
  TYPE_DISEASE     constant integer := 0;
  TYPE_DOMAIN      constant integer := 1;
  TYPE_ANATOMICAL  constant integer := 2;

  SNOMED_CODING_ID constant integer := 9;

  ROOT_CONCEPT_ID  constant varchar := '138875005';
  RELATIONSHIP_ID  constant varchar := '116680003';
begin
  -- drop table if exists
  if exists(
    select 1
      from information_schema.tables
     where table_schema = 'public'
       and (table_name = 'clinicalcode_ontologytag' or table_name = 'clinicalcode_ontologytagedge')
  ) then
    drop table if exists public.clinicalcode_ontologytag cascade;
    drop table if exists public.clinicalcode_ontologytagedge cascade;
  end if;

  -- create ontologytag table
  create table public.clinicalcode_ontologytag (
    id             bigserial     primary key,
    name           varchar(256)  not null,
    type_id        integer       not null,
    reference_id   bigint        default null,
    properties     jsonb         default '{}'::jsonb,
    search_vector  tsvector      default ''
  );

  -- create ontologytagedge table
  create table public.clinicalcode_ontologytagedge (
    id        bigserial primary key,
    child_id  bigint    references clinicalcode_ontologytag (id) not null,
    parent_id bigint    references clinicalcode_ontologytag (id) not null,
    unique (child_id, parent_id)
  );

  -- generate snomed ontology
  insert into public.clinicalcode_ontologytag (name, type_id, reference_id, properties, search_vector)
    select
          description,
          TYPE_DISEASE as type_id,
          null as reference_id,
          json_build_object(
                        'code', code,
            'coding_system_id', SNOMED_CODING_ID
          ) as properties,
          setweight(
            (
              to_tsvector('pg_catalog.english', coalesce(description, '')) ||
              to_tsvector('pg_catalog.english', coalesce(code, ''))
            ),
            'A'
          ) as search_vector
      from public.clinicalcode_snomed_codes
     where code != ROOT_CONCEPT_ID;

  -- generate ontology relationship(s)
  with
    ontology as (
      select
            id,
            properties::json->>'code'::varchar as code
      from public.clinicalcode_ontologytag
    ),
    relationships as (
      select *
        from public.clinicalcode_snomed_relationship
       where active = true
         and RELATIONSHIP_ID = type_id
         and ROOT_CONCEPT_ID not in (source_id, destination_id)
    )
  insert into public.clinicalcode_ontologytagedge (child_id, parent_id)
    select
          c.id as child_id,
          p.id as parent_id
      from relationships as r
      join ontology as c
        on r.source_id = c.code
      join ontology as p
        on r.destination_id = p.code;

end;
$tx$ language plpgsql;

--![endtemplate]
