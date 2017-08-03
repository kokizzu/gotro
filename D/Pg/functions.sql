
-- trigger 
-- 2017-067-25 Prayogo
CREATE OR REPLACE FUNCTION timestamp_changer() RETURNS trigger AS $$
DECLARE  
  changed BOOLEAN  := FALSE;  
  log_table TEXT := quote_ident('_log_' || TG_TABLE_NAME);  
  info TEXT := '';   
  mod_time TIMESTAMP := CURRENT_TIMESTAMP;   
  actor BIGINT;   
  query TEXT := '';  
BEGIN    
  IF (OLD.unique_id <> NEW.unique_id) THEN   
    NEW.updated_at := mod_time;  
    actor := NEW.updated_by;  
    changed := TRUE;    
    IF info <> '' THEN info := info || chr(10); END IF;  
    info := info || 'unique' || E'\t' || OLD.unique_id || E'\t' || NEW.unique_id;   
  END IF;   
  IF (OLD.is_deleted = TRUE) AND (NEW.is_deleted = FALSE) THEN    
    NEW.restored_at := mod_time;    
    actor := NEW.restored_by;    
    IF info <> '' THEN info := info || chr(10); END IF;  
    info := info || 'restore';   
    changed := TRUE;    
  END IF;   
  IF (OLD.is_deleted = FALSE) AND (NEW.is_deleted = TRUE) THEN    
    NEW.deleted_at := mod_time;  
    actor := NEW.deleted_by;  
    IF info <> '' THEN info := info || chr(10); END IF;  
    info := info || 'delete';    
    changed := TRUE;    
  END IF;   
  IF (OLD.data <> NEW.data) THEN    
    NEW.updated_at := mod_time;  
    IF info <> '' THEN info := info || chr(10); END IF;  
    info := info || 'update';    
    query := 'INSERT INTO ' || log_table || '( record_id, user_id, date, info, data_before, data_after )' || ' VALUES(' || OLD.id || ',' || NEW.updated_by || ',' || quote_literal(mod_time) || ',' || quote_literal(info) || ',' || quote_literal(NEW.data) || ',' || quote_literal(OLD.data) || ')';
    EXECUTE query;   
    changed := TRUE;    
  ELSEIF changed THEN   
    query := 'INSERT INTO ' || log_table || '( record_id, user_id, date, info )' || ' VALUES(' || OLD.id || ',' || actor || ',' || quote_literal(mod_time) || ',' || quote_literal(info) || ')';
    EXECUTE query;   
  END IF;   
  IF changed THEN NEW.modified_at := mod_time; END IF;   
  RETURN NEW;  
END;
$$ language plpgsql;

-- merge json
-- 2015-02-26 Prayogo
CREATE OR REPLACE FUNCTION jsonb_merge(JSONB, JSONB)
RETURNS JSONB AS $$
WITH json_union AS (
SELECT * FROM JSONB_EACH($1)
UNION ALL
SELECT * FROM JSONB_EACH($2)
) SELECT JSON_OBJECT_AGG(key, value)::JSONB
FROM json_union
WHERE key NOT IN (SELECT key FROM json_union WHERE value ='null');
$$ LANGUAGE SQL;

-- english numerals
-- 2015-07-26 Prayogo
CREATE OR REPLACE FUNCTION eng_num(num INT) RETURNS TEXT AS $$
BEGIN
RETURN CASE num % 10
WHEN 1 THEN num || 'st'
WHEN 2 THEN num || 'nd'
WHEN 3 THEN num || 'rd'
ELSE num || 'th'
END;
END
$$ LANGUAGE plpgsql;

-- alphanumeric string
-- 2015-06-23 Prayogo
CREATE OR REPLACE FUNCTION alnum_str(str TEXT) RETURNS TEXT AS $$
BEGIN
RETURN CASE
WHEN str IS NULL THEN ''
ELSE LOWER(REGEXP_REPLACE(str,'[^A-Za-z0-9]','','g'))
END;
END
$$ LANGUAGE plpgsql;


-- function to trim dates, returns NULL if dates is considered empty
-- 2015-03-26 Prayogo
CREATE OR REPLACE FUNCTION trim_dates(str TEXT) RETURNS TEXT AS $$
BEGIN
RETURN CASE WHEN (str IS NULL)
OR str = '' OR TRIM(BOTH FROM str) = ''
OR str = '0001-01-01T00:00:00Z' OR str = '0001-01-01' OR str = '0000-01-01'
THEN NULL
ELSE str
END;
END
$$ LANGUAGE plpgsql;

-- function to trim a single date, returns NULL if date is considered empty, get first 10 characters (remove time)
-- 2015-03-26 Prayogo TODO: consider to remove this, since it not used anymore
CREATE OR REPLACE FUNCTION trim_date(str TEXT) RETURNS TEXT AS $$
BEGIN
RETURN CASE WHEN (str IS NULL)
OR str = '' OR TRIM(BOTH FROM str) = ''
OR str = '0001-01-01T00:00:00Z' OR str = '0001-01-01' OR str = '0000-01-01'
THEN NULL
ELSE LEFT(TRIM(BOTH FROM str),10)
END;
END
$$ LANGUAGE plpgsql;

-- function to trim a single date, returns NULL if date is considered empty, get first 7 characters (remove time and day)
-- 2015-03-26 Prayogo
CREATE OR REPLACE FUNCTION yymm_date(str TEXT) RETURNS TEXT AS $$
BEGIN
RETURN CASE WHEN (str IS NULL)
OR str = '' OR TRIM(BOTH FROM str) = ''
OR str = '0001-01-01T00:00:00Z' OR str = '0001-01-01' OR str = '0000-01-01'
THEN NULL
ELSE LEFT(TRIM(BOTH FROM str),7)
END;
END
$$ LANGUAGE plpgsql;

-- convert integer to letter code for class in feeder/epsbed
-- 2016-01-22 Prayogo
CREATE OR REPLACE FUNCTION to_code(bigint) returns text
AS $$
DECLARE
	ov BIGINT;
	nv TEXT;
BEGIN
	ov := $1;
	nv := '';
	WHILE ov > 0 LOOP
		nv := nv || chr( (65+ov % 26)::INT );
		ov := ov / 26;
	END LOOP;
	RETURN nv;
END
	$$ LANGUAGE plpgsql;

-- calculate distance between 2 positions based on lat & long
CREATE OR REPLACE FUNCTION distance(lat1 float, long1 float, lat2 float, long2 float) RETURNS float AS $$
DECLARE
	p FLOAT = 0.017453292519943295;
	d FLOAT ;
BEGIN
	d = 12742 * asin(sqrt((0.5 - COS((lat2 - (lat1)) * p)/2 + COS(lat1 * p) * COS(lat2 * p) *
	(1 - COS((long2 - (long1)) * p))/2 )));
	RETURN d;
END
$$ LANGUAGE plpgsql;

-- get emails column
CREATE OR REPLACE FUNCTION emails_join(data JSONB) RETURNS TEXT AS $$
DECLARE
	emails TEXT;
BEGIN
	emails = ARRAY_TO_STRING(ARRAY[
		ARRAY[
			(CASE WHEN data->>'gmail' = '' THEN NULL ELSE data->>'gmail' END),
			(CASE WHEN data->>'yahoo' = '' THEN NULL ELSE data->>'yahoo' END),
			(CASE WHEN data->>'email' = '' THEN NULL ELSE data->>'email' END),
			(CASE WHEN data->>'outlook' = '' THEN NULL ELSE data->>'outlook' END),
			(CASE WHEN data->>'office_mail' = '' THEN NULL ELSE data->>'office_mail' END)
		]
	],', ');
	RETURN emails;
END
$$ LANGUAGE plpgsql;

-- check number
CREATE OR REPLACE FUNCTION is_num(text) RETURNS BOOLEAN AS $$
DECLARE x NUMERIC;
BEGIN
    x = $1::NUMERIC;
    RETURN TRUE;
EXCEPTION WHEN others THEN
    RETURN FALSE;
END;
$$
STRICT
LANGUAGE plpgsql IMMUTABLE;