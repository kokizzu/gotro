package Pg

import (
	"github.com/kokizzu/gotro/S"
	// TODO: replace with faster one `github.com/jackc/pgx/stdlib`, see https://github.com/jackc/pgx/issues/81
	// https://jmoiron.github.io/sqlx/
	// https://github.com/jmoiron/sqlx
	// https://sourcegraph.com/github.com/jmoiron/sqlx
	"github.com/kokizzu/gotro/D"
)

var Z func(string) string
var ZZ func(string) string
var ZJ func(string) string
var ZI func(int64) string
var ZLIKE func(string) string
var ZS func(string) string

var DEBUG bool

func init() {
	Z = S.Z
	ZZ = S.ZZ
	ZJ = S.ZJJ
	ZI = S.ZI
	ZLIKE = S.ZLIKE
	ZS = S.ZS
	DEBUG = D.DEBUG
}

const SQL_FUNCTIONS = `
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
LANGUAGE plpgsql IMMUTABLE;`

func InitFunctions(conn *RDBMS) {
	conn.InitTrigger()
	conn.DoTransaction(func(tx *Tx) string {
		tx.DoExec(SQL_FUNCTIONS)
		return ``
	})
}
