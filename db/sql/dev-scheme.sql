CREATE TABLE xmp_billing_partner_operators
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_bp INTEGER,
  operator_code INTEGER,
  revenue_sharing INTEGER,
  created_at TIMESTAMP
);
CREATE TABLE xmp_billing_partners
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(256),
  country_code INTEGER,
  bp_operators VARCHAR(256),
  description VARCHAR(512),
  revenue_sharing VARCHAR(256),
  created_at TIMESTAMP DEFAULT now() NOT NULL
);
CREATE TABLE xmp_campaign_ratio_counter
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_campaign INTEGER,
  counter INTEGER
);
CREATE INDEX index_for_xmp_campaign_ratio_counter_field_id_campaign
  ON xmp_campaign_ratio_counter (id_campaign);

CREATE TABLE xmp_campaigns
(
  id SERIAL PRIMARY KEY NOT NULL,
  link VARCHAR(512),
  name VARCHAR(128) NOT NULL,
  id_publisher INTEGER NOT NULL,
  priority INTEGER,
  status INTEGER,
  page_welcome VARCHAR(32) DEFAULT '' NOT NULL,
  page_msisdn VARCHAR(32) DEFAULT '' NOT NULL,
  page_pin VARCHAR(32) DEFAULT '' NOT NULL,
  page_thank_you VARCHAR(32) DEFAULT '' NOT NULL,
  page_error VARCHAR(2047) DEFAULT '' NOT NULL,
  page_success VARCHAR(2047) not null default '',
  page_banner VARCHAR(32) DEFAULT '' NOT NULL,
  description VARCHAR(250) DEFAULT '' NOT NULL,
  blacklist INTEGER DEFAULT 0 NOT NULL,
  whitelist INTEGER DEFAULT 0 NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  service_id INTEGER NOT NULL,
  capping INTEGER,
  capping_currency INTEGER,
  cpa INTEGER,
  hash VARCHAR(32) DEFAULT '' NOT NULL,
  type INTEGER DEFAULT 1 NOT NULL,
  created_by INTEGER NOT NULL,
  autoclick_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  autoclick_ratio INT NOT NULL DEFAULT 1
);

CREATE TABLE xmp_campaigns_access (
  id                          serial PRIMARY KEY,
  tid                         varchar(127) NOT NULL DEFAULT '',
  created_at                  TIMESTAMP NOT NULL DEFAULT NOW(),
  sent_at                     TIMESTAMP NOT NULL DEFAULT NOW(),
  msisdn                      varchar(32) NOT NULL DEFAULT '',
  ip                          varchar(32) NOT NULL DEFAULT '',
  os                          varchar(127) NOT NULL DEFAULT '',
  device                      varchar(127) NOT NULL DEFAULT '',
  browser                     varchar(127) NOT NULL DEFAULT '',
  operator_code               INT NOT NULL DEFAULT 0,
  country_code                INT NOT NULL DEFAULT 0,
  supported                   boolean NOT NULL DEFAULT FALSE,
  user_agent                  varchar(4091) NOT NULL DEFAULT '',
  referer                     varchar(4091) NOT NULL DEFAULT '',
  url_path                    varchar(4091) NOT NULL DEFAULT '',
  method                      varchar(127) NOT NULL DEFAULT '',
  headers                     VARCHAR(4091) NOT NULL DEFAULT '',
  error                       varchar(4091) NOT NULL DEFAULT '',
  id_campaign                 INT NOT NULL DEFAULT 0,
  id_service                  INT NOT NULL DEFAULT 0,
  id_content                  INT NOT NULL DEFAULT 0,
  geoip_country               varchar(127) NOT NULL DEFAULT '',
  geoip_iso                   varchar(127) NOT NULL DEFAULT '',
  geoip_city                  varchar(127) NOT NULL DEFAULT '',
  geoip_timezone              varchar(127) NOT NULL DEFAULT '',
  geoip_latitude              DOUBLE PRECISION NOT NULL DEFAULT .0,
  geoip_longitude             DOUBLE PRECISION NOT NULL DEFAULT .0,
  geoip_metro_code            int NOT NULL DEFAULT 0,
  geoip_postal_code           varchar(127) NOT NULL DEFAULT '',
  geoip_subdivisions          varchar(511) NOT NULL DEFAULT '',
  geoip_is_anonymous_proxy    boolean NOT NULL DEFAULT FALSE,
  geoip_is_satellite_provider boolean NOT NULL DEFAULT FALSE,
  geoip_accuracy_radius       int NOT NULL DEFAULT 0
);

create index xmp_campaigns_access_sent_at_idx
  on xmp_campaigns_access(sent_at);
create index xmp_campaigns_access_id_service_idx
  on xmp_campaigns_access(id_service);
create index xmp_campaigns_access_id_campaign_idx
  on xmp_campaigns_access(id_campaign);
create index xmp_campaigns_access_tid_idx
  on xmp_campaigns_access(tid);

CREATE TABLE xmp_cheese_dynamic_url_log
(
  id SERIAL PRIMARY KEY NOT NULL,
  url VARCHAR(512),
  id_service INTEGER,
  created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE xmp_content
(
  id SERIAL PRIMARY KEY NOT NULL,
  content_name VARCHAR(256),
  id_category INTEGER,
  id_sub_category INTEGER,
  publisher_name VARCHAR(256),
  id_platform INTEGER,
  id_uploader INTEGER,
  id_publisher INTEGER,
  status INTEGER,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  name VARCHAR(256) DEFAULT ''::character varying,
  object VARCHAR(32) DEFAULT ''::character varying,
  id_content_provider INTEGER NOT NULL
);
CREATE TABLE xmp_content_blacklist
(
  id SERIAL PRIMARY KEY NOT NULL,
  category VARCHAR(32),
  id_unit INTEGER,
  id_country INTEGER
);
CREATE TABLE xmp_campaigns_keywords
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_campaign INTEGER NOT NULL,
  keyword varchar(255) NOT NULL
);
CREATE TABLE xmp_content_category
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64),
  icon VARCHAR(64)
);
CREATE TABLE xmp_content_links
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_content INTEGER,
  link VARCHAR(512),
  created_at TIMESTAMP DEFAULT now(),
  ttl_hours INTEGER,
  id_subscription INTEGER,
  counter INTEGER,
  status INTEGER
);

CREATE TABLE xmp_content_platforms
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64)
);
CREATE TABLE xmp_content_providers
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32) NOT NULL
);
CREATE TABLE xmp_content_publishers
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32) NOT NULL,
  description VARCHAR(512)
);
CREATE TABLE xmp_content_sent
(
  id SERIAL PRIMARY KEY NOT NULL,
  tid VARCHAR(127) NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  sent_at TIMESTAMP DEFAULT now() NOT NULL,
  msisdn VARCHAR(32) DEFAULT ''::character varying NOT NULL,
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  id_service INTEGER DEFAULT 0 NOT NULL,
  id_content INTEGER DEFAULT 0 NOT NULL,
  id_subscription INTEGER DEFAULT 0 NOT NULL,
  operator_code INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL
);
create index xmp_content_sent_sent_at_idx
  on xmp_content_sent(sent_at);
create index xmp_content_sent_id_campaign_idx
  on xmp_content_sent(id_campaign);
create index xmp_content_sent_id_service_idx
  on xmp_content_sent(id_service);
create index xmp_content_sent_tid_idx
  on xmp_content_sent(tid);


CREATE TABLE xmp_content_unique_urls (
  id SERIAL PRIMARY KEY NOT NULL,
  tid VARCHAR(127) NOT NULL,
  sent_at TIMESTAMP DEFAULT now() NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  msisdn VARCHAR(32) DEFAULT ''::character varying NOT NULL,
  id_campaign INTEGER NOT NULL,
  id_service INTEGER NOT NULL,
  id_content INTEGER NOT NULL,
  id_subscription INTEGER NOT NULL,
  operator_code INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL,
  content_path VARCHAR (255) NOT NULL,
  content_name VARCHAR (255) NOT NULL,
  unique_url VARCHAR (255) NOT NULL
);
create index xmp_content_unique_urls_sent_at_idx
  on xmp_content_unique_urls(sent_at);


CREATE TABLE xmp_content_sub_category
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64),
  description VARCHAR(512)
);
CREATE TABLE xmp_conversion_report
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_report VARCHAR(32),
  lp INTEGER,
  msisdn INTEGER,
  pin INTEGER,
  subscribed INTEGER,
  first_charge INTEGER,
  subscribers_rate DOUBLE PRECISION,
  msisdn_rate DOUBLE PRECISION,
  pin_rate DOUBLE PRECISION,
  tariffication_rate DOUBLE PRECISION,
  date DATE,
  link VARCHAR(64),
  thank_you INTEGER,
  unsubscribed INTEGER,
  uniq_subscribed INTEGER,
  country_code VARCHAR(8),
  operator_code VARCHAR(8)
);
CREATE TABLE xmp_countries
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(250),
  code INTEGER,
  status INTEGER,
  iso VARCHAR(32),
  priority INTEGER
);
CREATE TABLE xmp_cqrs_log
(
  id SERIAL PRIMARY KEY NOT NULL,
  service_name VARCHAR(64),
  service_method VARCHAR(64),
  data VARCHAR(256),
  created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE xmp_currency
(
  id SERIAL PRIMARY KEY NOT NULL,
  code VARCHAR(3)
);
CREATE TABLE xmp_inject_jobs
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64),
  description VARCHAR(256),
  wording VARCHAR(256),
  country_code INTEGER,
  operator_code INTEGER,
  arpu INTEGER,
  status_msisdn INTEGER,
  id_platform INTEGER,
  total_msisdn INTEGER,
  status_job INTEGER,
  created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE xmp_mobilink_queue
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  campaign VARCHAR(64),
  headers JSON
);
CREATE TABLE xmp_mobilink_revenue_report
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  total_mo INTEGER,
  total_mo_uniq INTEGER,
  total_mo_uniq_success_charge INTEGER,
  total_mo_success_charge INTEGER,
  total_mo_failed_charge INTEGER,
  total_mo_revenue INTEGER,
  total_retry INTEGER,
  total_uniq_retry INTEGER,
  total_success_retry INTEGER,
  total_failed_retry INTEGER,
  total_mo_success_rate DOUBLE PRECISION,
  total_mo_uniq_success_rate DOUBLE PRECISION,
  total_retry_success_rate DOUBLE PRECISION,
  id_report INTEGER,
  report_date DATE,
  total_lp_hits INTEGER,
  total_mo_uniq_30_days INTEGER,
  total_retry_revenue INTEGER,
  total_revenue INTEGER
);
CREATE TABLE xmp_mobilink_transactions
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  campaign VARCHAR(64),
  msisdn VARCHAR(16),
  status INTEGER,
  id_task INTEGER,
  type INTEGER
);
CREATE TABLE xmp_msisdn_postpaid
(
  id SERIAL PRIMARY KEY NOT NULL,
  msisdn VARCHAR(32)
);
CREATE TABLE xmp_msisdn_blacklist
(
  id SERIAL PRIMARY KEY NOT NULL,
  msisdn VARCHAR(32)
);
CREATE TABLE xmp_msisdn_whitelist
(
  id SERIAL PRIMARY KEY NOT NULL,
  msisdn VARCHAR(32)
);
CREATE TABLE xmp_operator_ip
(
  id SERIAL PRIMARY KEY NOT NULL,
  ip_from VARCHAR(32),
  ip_to VARCHAR(32),
  operator_code INTEGER,
  country_code INTEGER DEFAULT 0 NOT NULL
);
CREATE TABLE xmp_operator_msisdn_prefix
(
  id SERIAL PRIMARY KEY NOT NULL,
  operator_code INTEGER,
  prefix VARCHAR(8)
);


CREATE TABLE xmp_operators
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64),
  country_code INTEGER,
  isp VARCHAR(250),
  msisdn_prefix VARCHAR(128),
  mcc VARCHAR(8),
  mnc VARCHAR(8),
  created_at TIMESTAMP DEFAULT now(),
  status INTEGER DEFAULT 1,
  code INTEGER,
  mt_url VARCHAR(255) DEFAULT ''::character varying NOT NULL,
  settings JSONB DEFAULT '{}'::jsonb NOT NULL,
  rps INTEGER DEFAULT 0 NOT NULL,
  msisdn_headers JSONB DEFAULT '[]'::jsonb NOT NULL
);
CREATE TABLE xmp_payment_type
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32),
  description VARCHAR(64)
);
CREATE TABLE xmp_publishers
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(256),
  contact_person VARCHAR(256),
  regex VARCHAR(2047) NOT NULL DEFAULT '{23}',
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  publisher_code INTEGER,
  status INTEGER
);
CREATE TABLE xmp_publishers_attributes
(
  id SERIAL PRIMARY KEY NOT NULL,
  country_code INTEGER,
  operator_code INTEGER,
  id_service INTEGER,
  price INTEGER,
  id_currency INTEGER,
  id_publisher INTEGER,
  callback_url VARCHAR(512),
  publisher_code INTEGER,
  id_campaign INTEGER,
  cpa_ratio DOUBLE PRECISION,
  status INTEGER
);

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

CREATE TABLE xmp_retries_expired
(
  id SERIAL PRIMARY KEY NOT NULL,
  status retry_status NOT NULL DEFAULT  '',
  tid varchar(127) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  updated_at TIMESTAMP DEFAULT now() NOT NULL,
  expired_at TIMESTAMP DEFAULT now() NOT NULL,
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
  id_campaign INTEGER NOT NULL
);
create index xmp_retries_expired_last_pay_attempt_at_idx on xmp_retries_expired (last_pay_attempt_at);
create index xmp_retries_expired_created_at_idx on xmp_retries_expired(created_at);


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

CREATE TABLE xmp_revenue_report
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_report INTEGER,
  date DATE,
  total_subs INTEGER,
  new_subs INTEGER,
  revenue INTEGER,
  country_code INTEGER,
  operator_code INTEGER,
  id_service INTEGER,
  transactions_count INTEGER,
  unsuccess_transactions_count INTEGER,
  unsubscribed_count INTEGER,
  success_revenue_mobilink INTEGER,
  success_charged_mobilink INTEGER,
  mo_new_mobilink INTEGER,
  mo_unique_mobilink INTEGER,
  total_failed_hits_mobilink INTEGER,
  total_retry_hits_mobilink INTEGER,
  unique_retry_hits_mobilink INTEGER
);

CREATE TABLE xmp_roles
(
  id SERIAL PRIMARY KEY NOT NULL,
  role_name VARCHAR(128),
  created_at TIMESTAMP DEFAULT now(),
  created_by INTEGER
);
CREATE TABLE xmp_sections
(
  id SERIAL PRIMARY KEY NOT NULL,
  section_name VARCHAR(128)
);
CREATE TABLE xmp_service_content
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_service INTEGER,
  id_content INTEGER,
  status INTEGER
);
CREATE TABLE xmp_service_country_settings
(
  id SERIAL PRIMARY KEY NOT NULL,
  country_code INTEGER,
  operator_code INTEGER,
  id_service INTEGER
);

CREATE TABLE xmp_services
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  status INTEGER NOT NULL DEFAULT 1,
  id_payment_type INTEGER, -- reg pull chechelan
  id_currency INTEGER NOT NULL,
  price DOUBLE PRECISION NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL,
  name VARCHAR(32) NOT NULL,
  description VARCHAR(32) NOT NULL DEFAULT '',
  short_number VARCHAR(255) NOT NULL DEFAULT '',
  paid_hours INTEGER DEFAULT 0 NOT NULL,
  delay_hours INT NOT NULL DEFAULT 10,
  keep_days INTEGER NOT NULL,
  not_paid_text VARCHAR(255) NOT NULL DEFAULT '',
  send_not_paid_text_enabled bool not null default false,
  days JSONB DEFAULT '[]'::jsonb NOT NULL, -- ['','any','sun','mon','tue','wed','thu','fri','sat']
  allowed_from INT NOT NULL not null default 0,
  allowed_to INT NOT NULL not null default 0,
  send_content_text_template VARCHAR(255) NOT NULL DEFAULT '%v'
);

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

CREATE TABLE xmp_subscription_type
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32),
  description VARCHAR(64)
);

CREATE TABLE xmp_subscriptions_active
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_service INTEGER,
  country_code INTEGER,
  created_at TIMESTAMP DEFAULT now(),
  msisdn VARCHAR(32),
  status INTEGER,
  operator_code INTEGER
);

CREATE TABLE public.xmp_pixel_buffer (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  sent_at TIMESTAMP NOT NULL DEFAULT now(),
  id_campaign INT NOT NULL,
  tid CHARACTER VARYING(127) NOT NULL DEFAULT '',
  pixel VARCHAR(511) NOT NULL DEFAULT ''
);
create index xmp_pixel_buffer_id_campaign_idx
  on xmp_pixel_buffer(id_campaign);
create index xmp_pixel_buffer_id_sent_at_idx
  on xmp_pixel_buffer(sent_at);

CREATE TABLE public.xmp_pixel_settings (
  id SERIAL PRIMARY KEY ,
  id_campaign INT NOT NULL,
  operator_code INTEGER NOT NULL DEFAULT 0,
  country_code INTEGER NOT NULL DEFAULT 0,
  publisher VARCHAR(511) NOT NULL DEFAULT '',
  endpoint VARCHAR(2047) NOT NULL DEFAULT '',
  timeout INT NOT NULL DEFAULT 30,
  enabled BOOLEAN NOT NULL DEFAULT false,
  ratio INT NOT NULL DEFAULT 2
);

CREATE TABLE public.xmp_pixel_transactions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  sent_at TIMESTAMP NOT NULL DEFAULT now(),
  tid CHARACTER VARYING(127) NOT NULL DEFAULT '',
  msisdn CHARACTER VARYING(32) NOT NULL DEFAULT '',
  id_campaign INTEGER NOT NULL DEFAULT 0,
  operator_code INTEGER NOT NULL DEFAULT 0,
  country_code INTEGER NOT NULL DEFAULT 0,
  pixel VARCHAR(511) NOT NULL DEFAULT '',
  endpoint VARCHAR(511) NOT NULL DEFAULT '',
  publisher VARCHAR(511) NOT NULL DEFAULT '',
  response_code INT NOT NULL DEFAULT 0
);
create index xmp_pixel_transactions_sent_at_idx
  on xmp_pixel_transactions(sent_at);
create index xmp_pixel_transactions_id_campaign_idx
  on xmp_pixel_transactions(id_campaign);
create index xmp_pixel_transactions_pixel_idx
  on xmp_pixel_transactions(pixel);
create index xmp_pixel_transactions_publisher_idx
  on xmp_pixel_transactions(publisher);

CREATE INDEX index_xmp_subscriptions_active_id_service ON xmp_subscriptions_active (id_service);
CREATE INDEX index_xmp_subscriptions_id_service ON xmp_subscriptions_active (id_service);
CREATE INDEX index_xmp_subscriptions_active_country_code ON xmp_subscriptions_active (country_code);
CREATE INDEX index_xmp_subscriptions_country_code ON xmp_subscriptions_active (country_code);
CREATE INDEX index_xmp_subscriptions_active_on_created_at ON xmp_subscriptions_active (created_at);
CREATE INDEX index_xmp_subscriptions_created_at ON xmp_subscriptions_active (created_at);
CREATE INDEX index_xmp_subscriptions_active_operator_code ON xmp_subscriptions_active (operator_code);
CREATE INDEX index_xmp_subscriptions_operator_code ON xmp_subscriptions_active (operator_code);

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

CREATE TABLE xmp_transactions_dr (
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP,
  tran_type INTEGER,
  msisdn VARCHAR(32),
  country_code INTEGER,
  id_service INTEGER,
  status INTEGER,
  operator_code INTEGER
);
CREATE INDEX index_xmp_transactions_dr_on_created_at ON xmp_transactions_dr (created_at);
CREATE INDEX index_xmp_transactions_dr_on_tran_type ON xmp_transactions_dr (tran_type);
CREATE INDEX index_xmp_transactions_dr_on_country_code ON xmp_transactions_dr (country_code);
CREATE INDEX index_xmp_transactions_dr_on_id_service ON xmp_transactions_dr (id_service);
CREATE INDEX index_xmp_transactions_dr_on_operator_code ON xmp_transactions_dr (operator_code);
CREATE TABLE xmp_transactions_dr_test
(
  id SERIAL PRIMARY KEY NOT NULL,
  msisdn VARCHAR(32),
  id_service INTEGER,
  created_at TIMESTAMP
);
CREATE TABLE xmp_uniq_url
(
  id SERIAL PRIMARY KEY NOT NULL,
  file_src VARCHAR(128),
  file_uniq_linq VARCHAR(256),
  created_at TIMESTAMP DEFAULT now(),
  expired INTEGER
);

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


CREATE TABLE xmp_user_activity_logs
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  id_xmp_user INTEGER,
  activity VARCHAR(64),
  id_record INTEGER,
  ip_user VARCHAR(32)
);
CREATE TABLE xmp_user_data_reports
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  created_by INTEGER,
  report_date DATE,
  report_type VARCHAR(32),
  report_status INTEGER,
  file_link VARCHAR(128),
  id_report INTEGER,
  updated_at TIMESTAMP,
  campaign VARCHAR(64)
);
CREATE TABLE xmp_user_role_permissions
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_role INTEGER,
  id_section INTEGER,
  action_name VARCHAR(256),
  permission_name VARCHAR(256)
);
CREATE TABLE xmp_user_roles
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_user INTEGER,
  id_role INTEGER
);
CREATE TABLE xmp_user_roles_orig
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_user INTEGER,
  id_role INTEGER
);
-- (

CREATE TABLE xmp_users_types
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32),
  description VARCHAR(256)
);
CREATE TABLE xmp_users
(
  id SERIAL PRIMARY KEY NOT NULL,
  username VARCHAR(32),
  password VARCHAR(32),
  email VARCHAR(64),
  type INTEGER,
  active INTEGER,
  CONSTRAINT xmp_users_type_fkey FOREIGN KEY (type) REFERENCES xmp_users_types (id)
);
CREATE TABLE xmp_users_transaction_type
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32)
);
CREATE TABLE xmp_users_transactions
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_xmp_user INTEGER,
  id_xmp_users_transaction_type INTEGER,
  created_at TIMESTAMP DEFAULT now()
);


--
create schema tr;
CREATE TABLE tr.partners
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  name VARCHAR(127) NOT NULL
);

CREATE TABLE tr.partners_destinations
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_partner int NOT NULL,
  active bool not null default false,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  amount_limit INTEGER DEFAULT 0 NOT NULL,
  destination VARCHAR(2047) DEFAULT ''::character varying NOT NULL,
  rate_limit INT NOT NULL DEFAULT 0,
  price_per_hit DOUBLE PRECISION NOT NULL DEFAULT 0,
  operator_code INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL,
  score INT NOT NULL DEFAULT 0
);

CREATE TABLE tr.destinations_hits
(
  id SERIAL PRIMARY KEY NOT NULL,
  id_partner int NOT NULL,
  id_destination int NOT NULL,
  tid VARCHAR(127) NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  sent_at TIMESTAMP DEFAULT now() NOT NULL,
  destination VARCHAR(2048) DEFAULT ''::character varying NOT NULL,
  msisdn VARCHAR(32) DEFAULT ''::character varying NOT NULL,
  price_per_hit DOUBLE PRECISION NOT NULL DEFAULT 0,
  operator_code INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL
);

create index destinations_hits_sent_at_idx
  on tr.destinations_hits(sent_at);
create index destinations_hits_id_partner_idx
  on tr.destinations_hits(id_partner);
create index destinations_hits_id_destination_idx
  on tr.destinations_hits(id_destination);

-- pixel settings
-- insert INTO xmp_pixel_settings
-- (operator_code, country_code, publisher, endpoint, enabled, ratio, id_campaign)
-- VALUES (41001, 41, 'Mobusi', 'http://wap.singiku.com/pass/jojokucpaga.php?aff_sub=%pixel%', true, 1, 290);
-- insert INTO xmp_pixel_settings
-- (operator_code, country_code, publisher, endpoint, enabled, ratio, id_campaign)
-- VALUES (41001, 41, 'Kimia', 'http://wap.singiku.com/pass/jojokucpaga.php?aff_sub=%pixel%', true, 2, 290);
--
-- insert into subscriptions(msisdn, tid, price, delay_hours,days, allowed_from, allowed_to )
-- values (
--  '12312323', '1484727568-3ecf5094-b608-4ace-6a84-8837b921f58d',
--   90, 10, '[]', 10, 12
-- )