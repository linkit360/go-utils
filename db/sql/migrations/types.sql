-- CONSTRAINT xmp_retries_status_fk FOREIGN KEY (status) REFERENCES xmp_retry_statuses (name)
-- varchar(127) NOT NULL DEFAULT '',

-- retry status
CREATE TABLE xmp_retry_statuses (name VARCHAR(127) NOT NULL PRIMARY KEY);
INSERT INTO xmp_retry_statuses VALUES (''),('pending'),('script');
ALTER TABLE public.xmp_retries ALTER COLUMN status TYPE VARCHAR(127) USING status::VARCHAR(127);
ALTER TABLE public.xmp_retries ALTER COLUMN status SET DEFAULT '';
ALTER TABLE public.xmp_retries ADD CONSTRAINT xmp_retries_status_fk FOREIGN KEY(status)
  REFERENCES xmp_retry_statuses(name);

-- job status
CREATE TABLE xmp_job_statuses ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_job_statuses VALUES ('ready'),('in progress'),('canceled'), ('done'), ('error');
ALTER TABLE public.xmp_jobs ALTER COLUMN status TYPE VARCHAR(127) USING status::VARCHAR(127);
ALTER TABLE public.xmp_jobs ALTER COLUMN status SET DEFAULT 'ready';

-- job type
CREATE TABLE xmp_job_types ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_job_types VALUES ('injection'),('expired');
ALTER TABLE public.xmp_jobs ALTER COLUMN type TYPE VARCHAR(127) USING type::VARCHAR(127);

-- operator transaction log type
CREATE TABLE xmp_operator_transaction_log_types ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_operator_transaction_log_types VALUES ('mo'),('mt'),('callback'), ('consent'), ('charge');
ALTER TABLE public.xmp_operator_transaction_log ALTER COLUMN type TYPE VARCHAR(127) USING type::VARCHAR(127);
ALTER TABLE public.xmp_operator_transaction_log ALTER COLUMN type SET DEFAULT '';
ALTER TABLE public.xmp_operator_transaction_log
  ADD CONSTRAINT xmp_operator_transaction_log_type_fk FOREIGN KEY(type)
  REFERENCES xmp_operator_transaction_log_types(name);

-- subscription status
CREATE TABLE xmp_subscriptions_statuses ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_subscriptions_statuses VALUES (''), ('failed'), ('paid'), ('blacklisted'), ('postpaid'), ('rejected'), ('canceled'), ('pending');
ALTER TABLE public.xmp_subscriptions ALTER COLUMN result TYPE VARCHAR(127) USING result::VARCHAR(127);
ALTER TABLE public.xmp_subscriptions ALTER COLUMN result SET DEFAULT '';
ALTER TABLE public.xmp_subscriptions
  ADD CONSTRAINT xmp_subscriptions_status_fk FOREIGN KEY(result)
REFERENCES xmp_subscriptions_statuses(name);

-- transaction result
CREATE TABLE xmp_transactions_results ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_transactions_results VALUES
  ('failed'),
  ('sms'),
  ('paid'),
  ('retry_failed'),
  ('retry_paid'),
  ('rejected'),
  ('expired_paid'),
  ('expired_failed'),
  ('injection_paid'),
  ('injection_failed');
ALTER TABLE public.xmp_transactions ALTER COLUMN result TYPE VARCHAR(127) USING result::VARCHAR(127);
ALTER TABLE public.xmp_transactions
  ADD CONSTRAINT xmp_transactions_result_fk FOREIGN KEY(result)
REFERENCES xmp_transactions_results(name);

-- user action
CREATE TABLE xmp_user_actions_actions ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_user_actions_actions VALUES ('access'), ('pull_click'), ('content_get'), ('rejected'), ('redirect'), ('autoclick');
ALTER TABLE public.xmp_user_actions ALTER COLUMN action TYPE VARCHAR(127) USING action::VARCHAR(127);
