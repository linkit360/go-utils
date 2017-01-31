CREATE OR REPLACE FUNCTION  xmp_subscriptions_create_partition_and_insert() RETURNS trigger AS
$BODY$
DECLARE
  partition_date TEXT;
  partition TEXT;
  r xmp_subscriptions%rowtype;
BEGIN
  partition_date := to_char(NEW.sent_at,'YYYY_MM_DD');
  partition := TG_RELNAME || '_' || partition_date;
  IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
    RAISE NOTICE 'A partition has been created %',partition;

    EXECUTE 'CREATE TABLE ' || partition || ' ( ' ||
            'check ( date(sent_at) = ''' || partition_date||
            ''' ) ) INHERITS (' || TG_RELNAME || ');';

    EXECUTE 'CREATE INDEX ' || partition || '_sent_at_idx ON ' || partition || '(sent_at);';
    EXECUTE 'CREATE INDEX ' || partition || '_last_pay_attempt_at_idx ON ' || partition || '(last_pay_attempt_at);';
  END IF;

  EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING * ' INTO r;
  RETURN r;
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

-- https://wiki.postgresql.org/wiki/INSERT_RETURNING_vs_Partitioning article.
-- Check all variables inside the trigger and replace all empty variables with implicit NULL-s and check if it works.

CREATE TRIGGER xmp_subscriptions_insert_trigger
BEFORE INSERT ON xmp_subscriptions
FOR EACH ROW EXECUTE PROCEDURE xmp_subscriptions_create_partition_and_insert();
