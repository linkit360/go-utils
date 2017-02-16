drop table xmp_jobs;
CREATE TABLE xmp_job_statuses ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_job_statuses VALUES ('ready'),('in progress'),('canceled'), ('done'), ('error');
CREATE TABLE xmp_job_types ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_job_types VALUES ('injection'),('expired');
CREATE TABLE xmp_jobs
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_user INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  run_at TIMESTAMP DEFAULT now() NOT NULL,
  finished_at TIMESTAMP DEFAULT '1970-01-01 00:00:00'::timestamp without time zone NOT NULL,
  type varchar(127) NOT NULL,
  status varchar(127) NOT NULL DEFAULT 'ready',
  file_name varchar(127) NOT NULL DEFAULT '',
  log_path  varchar(1023) NOT NULL DEFAULT '',
  params JSONB DEFAULT '{}'::jsonb NOT NULL,
  skip INTEGER NOT NULL DEFAULT 0,
  CONSTRAINT xmp_jobs_status_fk FOREIGN KEY (status) REFERENCES xmp_job_statuses (name),
  CONSTRAINT xmp_jobs_type_fk FOREIGN KEY (type) REFERENCES xmp_job_types (name)
);
create index xmp_jobs_created_at_idx on xmp_jobs(created_at);


-- -----------------------

drop table xmp_retries;
-- drop table xmp_retry_statuses;
CREATE TABLE xmp_retry_statuses (name VARCHAR(127) NOT NULL PRIMARY KEY);
INSERT INTO xmp_retry_statuses VALUES (''),('pending'),('script');
CREATE TABLE xmp_retries
(
  id SERIAL PRIMARY KEY NOT NULL,
  status varchar(127) NOT NULL DEFAULT '',
  tid varchar(127) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL,
  last_pay_attempt_at TIMESTAMP DEFAULT now() NOT NULL,
  attempts_count INTEGER DEFAULT 1 NOT NULL,
  price INTEGER NOT NULL,
  keep_days INTEGER NOT NULL,
  delay_hours INTEGER NOT NULL,
  msisdn VARCHAR(32) NOT NULL,
  operator_code INTEGER NOT NULL,
  country_code INTEGER NOT NULL,
  id_service INTEGER NOT NULL,
  id_subscription INTEGER NOT NULL,
  id_campaign INTEGER NOT NULL,
  CONSTRAINT xmp_retries_status_fk FOREIGN KEY (status) REFERENCES xmp_retry_statuses (name)
);
create index xmp_retries_last_pay_attempt_at_idx on xmp_retries (last_pay_attempt_at);
create index xmp_retries_status_idx on xmp_retries(status);
create index xmp_retries_operator_code_idx on xmp_retries(operator_code);

-- -----------------------


drop table xmp_transactions CASCADE ;
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
CREATE TABLE xmp_transactions
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  sent_at TIMESTAMP NOT NULL DEFAULT NOW(),
  tid CHARACTER VARYING(127) NOT NULL DEFAULT '',
  msisdn VARCHAR(32) NOT NULL,
  country_code INTEGER NOT NULL DEFAULT 0,
  id_service INTEGER NOT NULL DEFAULT 0,
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  operator_code INTEGER NOT NULL DEFAULT 0,
  id_subscription INTEGER NOT NULL DEFAULT 0,
  id_content INTEGER NOT NULL DEFAULT 0,
  operator_token VARCHAR(511) NOT NULL,
  price INTEGER NOT NULL,
  result varchar(127) NOT NULL,
  CONSTRAINT xmp_transactions_result_fk FOREIGN KEY (result) REFERENCES xmp_transactions_results (name)
);
create index xmp_transactions_sent_at_idx
  on xmp_transactions(sent_at);
create index xmp_transactions_msisdn_idx
  on xmp_transactions(msisdn);
create index xmp_transactions_result_idx
  on xmp_transactions(result);
create index xmp_transactions_operator_token_idx
  on xmp_transactions(operator_token);

-- -----------------------

drop table xmp_operator_transaction_log CASCADE ;
CREATE TABLE xmp_operator_transaction_log_types ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_operator_transaction_log_types VALUES ('mo'),('mt'),('callback'), ('consent'), ('charge');
CREATE TABLE xmp_operator_transaction_log (
  id serial PRIMARY KEY,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  sent_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  tid  varchar(127) NOT NULL DEFAULT '',
  msisdn VARCHAR(32) NOT NULL,
  operator_code INTEGER NOT NULL,
  operator_time timestamp not null default now(),
  country_code INTEGER NOT NULL,
  operator_token CHARACTER VARYING(511) NOT NULL,
  error varchar(511) NOT NULL DEFAULT '',
  price INTEGER NOT NULL,
  id_service INTEGER NOT NULL,
  id_subscription INTEGER NOT NULL,
  id_campaign INTEGER NOT NULL,
  request_body varchar(16391) NOT NULL DEFAULT '',
  response_body varchar(16391) NOT NULL DEFAULT '',
  response_decision varchar(511) NOT NULL DEFAULT '',
  response_code INT NOT NULL DEFAULT 0,
  notice varchar(2047) NOT NULL DEFAULT '',
  type varchar(127) not null,
  CONSTRAINT xmp_operator_transaction_log_type_fk FOREIGN KEY (type) REFERENCES xmp_operator_transaction_log_types (name)
);

create index xmp_operator_transaction_log_sent_at_idx
  on xmp_operator_transaction_log(sent_at);
create index xmp_operator_transaction_log_type_idx
  on xmp_operator_transaction_log(type);
create index xmp_operator_transaction_log_notice_idx
  on xmp_operator_transaction_log(notice);
create index xmp_operator_transaction_log_msisdn_idx
  on xmp_operator_transaction_log(msisdn);
create index xmp_operator_transaction_log_id_service_idx
  on xmp_operator_transaction_log(id_service);
create index xmp_operator_transaction_log_id_campaign_idx
  on xmp_operator_transaction_log(id_campaign);
create index xmp_operator_transaction_log_id_operator_token_idx
  on xmp_operator_transaction_log(operator_token);

-- -----------------------

drop table xmp_subscriptions CASCADE ;
CREATE TABLE xmp_subscriptions_statuses ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_subscriptions_statuses VALUES (''), ('failed'), ('paid'), ('blacklisted'), ('postpaid'), ('rejected'), ('canceled'), ('pending');
CREATE TABLE xmp_subscriptions
(
  id SERIAL PRIMARY KEY NOT NULL,
  tid VARCHAR(127) DEFAULT ''::character varying NOT NULL,
  id_service INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  sent_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  msisdn VARCHAR(32),
  operator_code INTEGER DEFAULT 0 NOT NULL,
  operator_token VARCHAR(511) NOT NULL default '',
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  attempts_count INTEGER DEFAULT 0 NOT NULL,
  price INTEGER NOT NULL,
  result varchar(127) NOT NULL DEFAULT '',
  CONSTRAINT xmp_subscriptions_result_fk FOREIGN KEY (result) REFERENCES xmp_subscriptions_statuses (name),
  keep_days INTEGER DEFAULT 10 NOT NULL,
  last_pay_attempt_at TIMESTAMP DEFAULT now() NOT NULL,
  delay_hours INTEGER NOT NULL,
  paid_hours INTEGER DEFAULT 0 NOT NULL,
  pixel VARCHAR(511) NOT NULL DEFAULT '',
  publisher VARCHAR(511)NOT NULL DEFAULT '',
  pixel_sent boolean NOT NULL DEFAULT false,
  pixel_sent_at TIMESTAMP WITHOUT TIME ZONE,
  periodic bool not null default false,
  days JSONB NOT NULL not null default '[]',
  allowed_from INT NOT NULL not null default 11,
  allowed_to INT NOT NULL not null default 13
);
create index xmp_subscriptions_last_pay_attempt_at_idx
  on xmp_subscriptions (last_pay_attempt_at);
create index xmp_subscriptions_sent_at_idx
  on xmp_subscriptions (sent_at);
create index xmp_subscriptions_result_idx
  on xmp_subscriptions(result);
create index xmp_subscriptions_periodic_idx
  on xmp_subscriptions(periodic);
create index xmp_subscriptions_tid_idx
  on xmp_subscriptions(tid);

-- -----------------------

drop table xmp_user_actions CASCADE ;
CREATE TABLE xmp_user_actions_actions ( name VARCHAR(127) NOT NULL PRIMARY KEY );
INSERT INTO xmp_user_actions_actions VALUES ('access'), ('pull_click'), ('content_get'), ('rejected'), ('redirect'), ('autoclick');
CREATE TABLE xmp_user_actions (
  id serial PRIMARY KEY,
  tid  varchar(127) NOT NULL DEFAULT '',
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  msisdn    varchar(32) NOT NULL DEFAULT '',
  action varchar(127) NOT NULL,
  CONSTRAINT xmp_user_actions_action_fk FOREIGN KEY (action) REFERENCES xmp_user_actions_actions (name),
  error varchar(511) NOT NULL DEFAULT '',
  sent_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
create index xmp_user_actions_sent_at_idx
  on xmp_user_actions(sent_at);
create index xmp_user_actions_msisdn_idx
  on xmp_user_actions(msisdn);
create index xmp_user_actions_action_idx
  on xmp_user_actions(action);
create index xmp_user_actions_tid_idx
  on xmp_user_actions(tid);

