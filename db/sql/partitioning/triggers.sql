
CREATE TRIGGER xmp_user_actions_insert_trigger
BEFORE INSERT ON xmp_user_actions
FOR EACH ROW EXECUTE PROCEDURE xmp_user_actions_create_partition_and_insert();

CREATE TRIGGER xmp_subscriptions_insert_trigger
BEFORE INSERT ON xmp_subscriptions
FOR EACH ROW EXECUTE PROCEDURE xmp_subscriptions_create_partition_and_insert();

CREATE TRIGGER xmp_pixel_transactions_insert_trigger
BEFORE INSERT ON xmp_pixel_transactions
FOR EACH ROW EXECUTE PROCEDURE xmp_pixel_transactions_create_partition_and_insert();

CREATE TRIGGER xmp_operator_transaction_log_insert_trigger
BEFORE INSERT ON xmp_operator_transaction_log
FOR EACH ROW EXECUTE PROCEDURE xmp_operator_transaction_log_create_partition_and_insert();

CREATE TRIGGER xmp_content_sent_insert_trigger
BEFORE INSERT ON xmp_content_sent
FOR EACH ROW EXECUTE PROCEDURE xmp_content_sent_create_partition_and_insert();

CREATE TRIGGER xmp_campaigns_access_insert_trigger
BEFORE INSERT ON xmp_campaigns_access
FOR EACH ROW EXECUTE PROCEDURE xmp_campaigns_access_create_partition_and_insert();

CREATE TRIGGER xmp_transaction_insert_trigger
BEFORE INSERT ON xmp_transactions
FOR EACH ROW EXECUTE PROCEDURE xmp_transactions_create_partition_and_insert();


alter table xmp_user_actions disable trigger xmp_user_actions_insert_trigger;
alter table xmp_subscriptions disable trigger xmp_subscriptions_insert_trigger;
alter table xmp_pixel_transactions disable trigger  xmp_pixel_transactions_insert_trigger;
alter table xmp_operator_transaction_log disable trigger xmp_operator_transaction_log_insert_trigger;
alter table xmp_content_sent disable trigger xmp_content_sent_insert_trigger;
alter table xmp_campaigns_access disable trigger xmp_campaigns_access_insert_trigger;



drop trigger xmp_user_actions_insert_trigger on xmp_user_actions ;
drop trigger xmp_subscriptions_insert_trigger on xmp_subscriptions ;
drop trigger xmp_pixel_transactions_insert_trigger on xmp_pixel_transactions ;
drop trigger xmp_operator_transaction_log_insert_trigger on xmp_operator_transaction_log ;
drop trigger xmp_campaigns_access_insert_trigger on xmp_campaigns_access ;
drop trigger xmp_transaction_insert_trigger on xmp_transactions ;
