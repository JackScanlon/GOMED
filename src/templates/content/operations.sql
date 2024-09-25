/********************************************************************************
 * Cleanup                                                                      *
 *                                                                              *
 *  ? Cleanup managed tables & types                                            *
 *                                                                              *
 ********************************************************************************/

--![template] name:"cleanup"

do $tx$
declare
begin

  {{if or (.Data.all) (.Data.releases) }}
    drop table if exists clinicalcode_snomed_concept;
    drop table if exists clinicalcode_snomed_description;
    drop table if exists clinicalcode_snomed_relationship;
    drop table if exists clinicalcode_snomed_refset_lang;
    drop table if exists clinicalcode_snomed_refset_map;
    drop table if exists clinicalcode_snomed_ctv_map;
  {{end}}

  {{if or (.Data.all) (.Data.codelist) }}
    drop table if exists clinicalcode_snomed_codes cascade;
    drop type if exists sctident;
  {{end}}

  {{if or (.Data.all) (.Data.ontology) }}
    drop table if exists public.clinicalcode_ontologytag cascade;
    drop table if exists public.clinicalcode_ontologytagedge cascade;
  {{end}}

end;
$tx$ language plpgsql;

--![endtemplate]
