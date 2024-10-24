/********************************************************************************
 * Map table(s)                                                                 *
 *                                                                              *
 *  ? Create base source-target tables                                          *
 *                                                                              *
 ********************************************************************************/
create table public.temp_icd9_sct_map (
  id     serial        primary key,
  source varchar(7)    not null,
  target varchar(18)   not null
);

create table public.temp_mesh_sct_map (
  id     serial        primary key,
  source varchar(7)    not null,
  target varchar(18)   not null
);


/********************************************************************************
 * Build mapping                                                                *
 *                                                                              *
 *  ? Generate final mapping tables...                                          *
 *     -> icd9map: ICD-9 to SNOMED map                                          *
 *     -> meshmap: MeSH  to SNOMED map                                          *
 *                                                                              *
 ********************************************************************************/
do $tx$
begin
  -- drop table if exists
  if exists(
    select 1
      from information_schema.tables
    where table_schema = 'public'
      and table_name in ('temp_icd9map', 'temp_meshmap')
  ) then
    drop table if exists public.temp_icd9map cascade;
    drop table if exists public.temp_meshmap cascade;
  end if;

  -- create icd9->sct map table
  create table public.temp_icd9map (
    id          serial        primary key,
    codes       text[]        default '{}'::text[],
    snomed_code varchar(18)   not null,
    unique (snomed_code)
  );

  -- build final icd-9 map
  with
    codemap as (
        select
            target,
            array_agg(source)::text[] as codes
        from public.temp_icd9_sct_map
       group by target
    )
  insert into public.temp_icd9map (
    codes, snomed_code
  )
    select
          codes,
          target as snomed_code
      from codemap;

  -- create mesh->sct map table
  create table public.temp_meshmap (
    id          serial        primary key,
    codes       text[]        default '{}'::text[],
    snomed_code varchar(18)   not null,
    unique (snomed_code)
  );

  -- build final icd-9 map
  with
    codemap as (
        select
            target,
            array_agg(source)::text[] as codes
        from public.temp_icd9_sct_map
       group by target
    )
  insert into public.temp_icd9map (
    codes, snomed_code
  )
    select
          codes,
          target as snomed_code
      from codemap;
end;
$tx$ language plpgsql;
