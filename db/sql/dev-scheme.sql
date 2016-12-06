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
CREATE INDEX index_for_xmp_campaign_ratio_counter_field_id_campaign ON xmp_campaign_ratio_counter (id_campaign);
CREATE TABLE xmp_campaigns
(
  id SERIAL PRIMARY KEY NOT NULL,
  link VARCHAR(512),
  name VARCHAR(128) NOT NULL,
  id_publisher INTEGER NOT NULL,
  priority INTEGER,
  status INTEGER,
  page_welcome VARCHAR(32),
  page_msisdn VARCHAR(32),
  page_pin VARCHAR(32),
  page_thank_you VARCHAR(32),
  page_error VARCHAR(32),
  page_banner VARCHAR(32),
  description VARCHAR(250),
  blacklist INTEGER,
  whitelist INTEGER,
  created_at TIMESTAMP DEFAULT now(),
  service_id INTEGER,
  capping INTEGER,
  capping_currency INTEGER,
  cpa INTEGER,
  hash VARCHAR(32),
  type INTEGER DEFAULT 1 NOT NULL,
  created_by INTEGER,
  autoclick_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  autoclick_ratio INT NOT NULL DEFAULT 1
);

CREATE TABLE xmp_campaigns_access (
  id                          serial PRIMARY KEY,
  tid                         varchar(127) NOT NULL DEFAULT '',
  created_at                  TIMESTAMP NOT NULL DEFAULT NOW(),
  access_at                   TIMESTAMP NOT NULL DEFAULT NOW(),
  sent_at                     TIMESTAMP NOT NULL DEFAULT NOW(),
  msisdn                      varchar(32) NOT NULL DEFAULT '',
  ip                          varchar(32) NOT NULL DEFAULT '',
  os                          varchar(127) NOT NULL DEFAULT '',
  device                      varchar(127) NOT NULL DEFAULT '',
  browser                     varchar(127) NOT NULL DEFAULT '',
  operator_code               INT NOT NULL DEFAULT 0,
  country_code                INT NOT NULL DEFAULT 0,
  supported                   boolean NOT NULL DEFAULT FALSE,
  user_agent                  varchar(511) NOT NULL DEFAULT '',
  referer                     varchar(511) NOT NULL DEFAULT '',
  url_path                    varchar(511) NOT NULL DEFAULT '',
  method                      varchar(127) NOT NULL DEFAULT '',
  headers                     VARCHAR(2047) NOT NULL DEFAULT '',
  error                       varchar(511) NOT NULL DEFAULT '',
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
CREATE EXTENSION btree_gist;
CREATE INDEX xmp_campaigns_access_long_lat_gistidx ON xmp_campaigns_access USING gist(geoip_longitude, geoip_latitude);

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
CREATE UNIQUE INDEX xmp_content_links_id_uindex ON xmp_content_links (id);
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
CREATE TYPE retry_status AS ENUM ('', 'pending', 'script');
CREATE TABLE xmp_retries
(
  id SERIAL PRIMARY KEY NOT NULL,
  status retry_status NOT NULL DEFAULT  '',
  tid varchar(127) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  last_pay_attempt_at TIMESTAMP DEFAULT now() NOT NULL,
  attempts_count INTEGER DEFAULT 1 NOT NULL,
  keep_days INTEGER NOT NULL,
  delay_hours INTEGER NOT NULL,
  msisdn VARCHAR(32) NOT NULL,
  operator_code INTEGER NOT NULL,
  country_code INTEGER NOT NULL,
  id_service INTEGER NOT NULL,
  id_subscription INTEGER NOT NULL,
  id_campaign INTEGER NOT NULL
);

CREATE TABLE xmp_operator_transaction_log (
  id serial PRIMARY KEY,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  sent_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  tid  varchar(127) NOT NULL DEFAULT '',
  msisdn VARCHAR(32) NOT NULL,
  operator_code INTEGER NOT NULL,
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
  response_code INT NOT NULL DEFAULT 0
);


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
  name VARCHAR(32) NOT NULL,
  description VARCHAR(32),
  keyword VARCHAR(32),
  url VARCHAR(128),
  price DOUBLE PRECISION,
  id_payment_type INTEGER,
  id_subscription_type INTEGER,
  retry_days INTEGER,
  wording TEXT,
  status INTEGER,
  id_currency INTEGER,
  created_at TIMESTAMP DEFAULT now() NOT NULL,
  channel_sms INTEGER,
  channel_wap INTEGER,
  channel_web INTEGER,
  start_date TIMESTAMP,
  price_option VARCHAR(32),
  link VARCHAR(128),
  pull_msisdn_ttr INTEGER,
  pull_retry_delay INTEGER,
  sms_send INTEGER,
  paid_hours INTEGER DEFAULT 0 NOT NULL,
  delay_hours INT NOT NULL DEFAULT 10
);

CREATE TABLE xmp_subscription_type
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(32),
  description VARCHAR(64)
);
CREATE TYPE subscription_status AS ENUM
  ('', 'failed', 'paid', 'blacklisted', 'postpaid', 'rejected', 'past', 'canceled');

CREATE TABLE xmp_subscriptions
(
  id SERIAL PRIMARY KEY NOT NULL,
  tid VARCHAR(127) DEFAULT ''::character varying NOT NULL,
  last_success_date TIMESTAMP DEFAULT now(),
  id_service INTEGER DEFAULT 0 NOT NULL,
  country_code INTEGER DEFAULT 0 NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  msisdn VARCHAR(32),
  operator_code INTEGER DEFAULT 0 NOT NULL,
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  attempts_count INTEGER DEFAULT 0 NOT NULL,
  price INTEGER NOT NULL,
  result SUBSCRIPTION_STATUS DEFAULT ''::subscription_status NOT NULL,
  keep_days INTEGER DEFAULT 10 NOT NULL,
  last_pay_attempt_at TIMESTAMP DEFAULT now() NOT NULL,
  delay_hours INTEGER NOT NULL,
  paid_hours INTEGER DEFAULT 0 NOT NULL,
  pixel VARCHAR(511) NOT NULL DEFAULT '',
  publisher VARCHAR(511)NOT NULL DEFAULT '',
  pixel_sent boolean NOT NULL DEFAULT false,
  pixel_sent_at TIMESTAMP WITHOUT TIME ZONE
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

CREATE INDEX index_xmp_subscriptions_active_id_service ON xmp_subscriptions_active (id_service);
CREATE INDEX index_xmp_subscriptions_id_service ON xmp_subscriptions_active (id_service);
CREATE INDEX index_xmp_subscriptions_active_country_code ON xmp_subscriptions_active (country_code);
CREATE INDEX index_xmp_subscriptions_country_code ON xmp_subscriptions_active (country_code);
CREATE INDEX index_xmp_subscriptions_active_on_created_at ON xmp_subscriptions_active (created_at);
CREATE INDEX index_xmp_subscriptions_created_at ON xmp_subscriptions_active (created_at);
CREATE INDEX index_xmp_subscriptions_active_operator_code ON xmp_subscriptions_active (operator_code);
CREATE INDEX index_xmp_subscriptions_operator_code ON xmp_subscriptions_active (operator_code);
CREATE TABLE xmp_transaction_types
(
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(64)
);
CREATE TYPE transaction_result AS ENUM (
  'failed', 'sms', 'paid', 'retry_failed', 'retry_paid', 'rejected', 'past');

CREATE TABLE xmp_transactions
(
  id SERIAL PRIMARY KEY NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  tid CHARACTER VARYING(127) NOT NULL DEFAULT '',
  msisdn VARCHAR(32) NOT NULL,
  country_code INTEGER NOT NULL DEFAULT 0,
  id_service INTEGER NOT NULL DEFAULT 0,
  operator_code INTEGER NOT NULL DEFAULT 0,
  id_subscription INTEGER NOT NULL DEFAULT 0,
  id_content INTEGER NOT NULL DEFAULT 0,
  operator_token VARCHAR(511) NOT NULL,
  price INTEGER NOT NULL,
  result TRANSACTION_RESULT NOT NULL,
  id_campaign INTEGER DEFAULT 0 NOT NULL
);

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

CREATE TYPE user_action AS ENUM ('access', 'pull_click', 'content_get');
CREATE TABLE xmp_user_actions (
  id serial PRIMARY KEY,
  tid  varchar(127) NOT NULL DEFAULT '',
  id_campaign INTEGER DEFAULT 0 NOT NULL,
  msisdn    varchar(32) NOT NULL DEFAULT '',
  action user_action NOT NULL,
  error varchar(511) NOT NULL DEFAULT '',
  access_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

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



-- pixel settings
-- insert INTO xmp_pixel_settings
-- (operator_code, country_code, publisher, endpoint, enabled, ratio, id_campaign)
-- VALUES (41001, 41, 'Mobusi', 'http://wap.singiku.com/pass/jojokucpaga.php?aff_sub=%pixel%', true, 1, 290);
-- insert INTO xmp_pixel_settings
-- (operator_code, country_code, publisher, endpoint, enabled, ratio, id_campaign)
-- VALUES (41001, 41, 'Kimia', 'http://wap.singiku.com/pass/jojokucpaga.php?aff_sub=%pixel%', true, 2, 290);
--
