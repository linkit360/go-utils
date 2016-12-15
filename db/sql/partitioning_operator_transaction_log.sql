CREATE OR REPLACE FUNCTION xmp_operator_transaction_log_create_partition_and_insert() RETURNS trigger AS
$BODY$
DECLARE
  partition_date TEXT;
  partition TEXT;
BEGIN
  partition_date := to_char(NEW.date,'YYYY-MM-DD');
  partition := TG_RELNAME || '_' || partition_date;
  IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
    RAISE NOTICE 'A partition has been created %',partition;

    EXECUTE 'CREATE TABLE ' || partition ||
            ' (check (date = ''' || NEW.sent_at || ''')) INHERITS (' || TG_RELNAME || ');';
    EXECUTE 'CREATE INDEX ' || partition || '_created_at_idx ON ' || partition || '(created_at);';

  END IF;
  EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER xmp_operator_transaction_log_insert_trigger
BEFORE INSERT ON xmp_operator_transaction_log
FOR EACH ROW EXECUTE PROCEDURE xmp_operator_transaction_log_create_partition_and_insert();
