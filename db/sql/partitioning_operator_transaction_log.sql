CREATE OR REPLACE FUNCTION xmp_operator_transaction_log_create_partition_and_insert() RETURNS trigger AS
$BODY$
DECLARE
  partition_date TEXT;
  partition TEXT;
BEGIN
  partition_date := to_char(NEW.sent_at,'YYYY_MM_DD');
  partition := TG_RELNAME || '_' || partition_date;
  IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
    RAISE NOTICE 'A partition has been created %',partition;

    EXECUTE 'CREATE TABLE ' || partition || ' ( ' ||
            'check ( date(sent_at) = ''' || partition_date||
            ''' ) ) INHERITS (' || TG_RELNAME || ');';

    EXECUTE 'CREATE INDEX ' || partition || '_sent_at_idx ON ' || partition || '(sent_at);';

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
