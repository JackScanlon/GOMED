-- e.g. array_agg() / string_agg() but for tsvec
create aggregate tsvector_agg(tsvector) (
  stype = pg_catalog.tsvector,
  sfunc = pg_catalog.tsvector_concat,
  initcond = ''
);


-- examine codelist collection
with
  concepts as (
    select
          row_number() over (order by c.id::bigint) as id,
          c.id as code,
          d.term as description,
          (
            case
              when d.case_significance_id = '900000000000020002' then 'CL'::sctcase
              when d.case_significance_id = '900000000000017005' then 'CS'::sctcase
              when d.case_significance_id = '900000000000448009' then 'CI'::sctcase
              else 'CI'::sctcase
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
    )
  ),
  codemap as (
    select
          snomed_code,
          source,
          array_agg(code order by priority asc)::text[] as codes
      from map_components
     where nrow = 1
     group by snomed_code, source
  )
select *
  from codemap;


-- view hierarchy at depth
with
  recursive traversal(child_id, parent_id, depth, path) as (
    select
          first.child_id,
          first.parent_id,
          1 as depth,
          array[first.child_id] as path
      from public.clinicalcode_ontologytagedge as first
     union all
    select
          first.child_id,
          first.parent_id,
          second.depth + 1 as depth,
          path || first.child_id as path
      from public.clinicalcode_ontologytagedge as first,
           traversal as second
     where first.child_id = second.parent_id
       and first.child_id <> ALL(second.path)
  )
select *
  from traversal
 where depth < 2;


-- view root nodes
select *
  from public.clinicalcode_ontologytag as tag
  left join public.clinicalcode_ontologytagedge as edge
    on edge.child_id = tag.id
 where edge.child_id is null;
