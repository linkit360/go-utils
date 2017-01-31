CREATE OR REPLACE FUNCTION  xmp_user_actions_create_partition_and_insert()
  RETURNS trigger AS
$BODY$
DECLARE
  partition_date TEXT;
  partition TEXT;
BEGIN
  partition_date := to_char(NEW.sent_at,'YYYY_MM_DD');
  partition := TG_TABLE_NAME || '_' || partition_date;
  IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
    RAISE NOTICE 'A partition has been created %',partition;

    EXECUTE 'CREATE TABLE ' || partition || ' ( ' ||
            'check ( date(sent_at) = ''' || partition_date||
            ''' ) ) INHERITS (' || TG_TABLE_NAME || ');';

    EXECUTE 'CREATE INDEX ' || partition || '_sent_at_idx ON ' || partition || '(sent_at);';
    EXECUTE 'CREATE INDEX ' || partition || '_tid_idx ON ' || partition || '(tid);';
    EXECUTE 'CREATE INDEX ' || partition || '_msisdn_idx ON ' || partition || '(msisdn);';
    EXECUTE 'CREATE INDEX ' || partition || '_action_idx ON ' || partition || '(action);';
  END IF;

  EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_TABLE_NAME || ' ' || quote_literal(NEW) || ').* ';
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER xmp_user_actions_insert_trigger
BEFORE INSERT ON xmp_user_actions
FOR EACH ROW EXECUTE PROCEDURE xmp_user_actions_create_partition_and_insert();
