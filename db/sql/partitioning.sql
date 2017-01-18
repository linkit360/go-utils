-- access_campaign
-- user_actions
-- operator_transaction_log
-- subscriptions - created at result id campaigs
-- content_sent
-- pixel_transactions
-- transactions - по дате создания и msisdin составной индекс

CREATE OR REPLACE FUNCTION create_partition_and_insert() RETURNS trigger AS
  $BODY$
    DECLARE
      partition_date TEXT;
      partition TEXT;
    BEGIN
      partition_date := to_char(NEW.sent_at,'YYYY-MM-DD');
      partition := TG_RELNAME || '_' || partition_date;
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
        EXECUTE 'CREATE TABLE ' || partition ||
                ' (check (date = ''' || NEW.sent_at || ''')) INHERITS (' || TG_RELNAME || ');';
      END IF;
      EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
      RETURN NULL;
    END;
  $BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE TRIGGER testing_partition_insert_trigger
BEFORE INSERT ON testing_partition
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

--
-- CREATE VIEW show_partitions AS
--   SELECT nmsp_parent.nspname AS parent_schema,
--          parent.relname AS parent,
--          nmsp_child.nspname AS child_schema,
--          child.relname AS child
--   FROM pg_inherits
--     JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
--     JOIN pg_class child ON pg_inherits.inhrelid = child.oid
--     JOIN pg_namespace nmsp_parent ON nmsp_parent.oid = parent.relnamespace
--     JOIN pg_namespace nmsp_child ON nmsp_child.oid = child.relnamespace
--   WHERE parent.relname='testing_partition' ;