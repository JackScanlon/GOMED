/********************************************************************************
 * Copy CSV file                                                                *
 *                                                                              *
 *  ? Copy a CSV file to pgsql                                                  *
 *                                                                              *
 ********************************************************************************/

--![template] name:"file"
copy {{ .targetName | print }} from '{{ .filePath | print }}'
  with (
    DELIMITER {{ .delimiter | print }},
    FORMAT 'csv',
    ENCODING 'UTF-8',
    HEADER 'on'
  );
--![endtemplate]
