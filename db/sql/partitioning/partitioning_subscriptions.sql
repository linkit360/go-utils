CREATE OR REPLACE FUNCTION  xmp_subscriptions_create_partition_and_insert() RETURNS trigger AS
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
    EXECUTE 'CREATE INDEX ' || partition || '_sent_at_result_id_campaign_idx ON ' || partition || '(sent_at, result, id_campaign);';
    EXECUTE 'CREATE INDEX ' || partition || '_created_at_result_id_campaign_idx ON ' || partition || '(created_at, result, id_campaign);';
  END IF;
  EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER xmp_subscriptions_insert_trigger
BEFORE INSERT ON xmp_subscriptions
FOR EACH ROW EXECUTE PROCEDURE xmp_subscriptions_create_partition_and_insert();
