CREATE TABLE xmp_content_sent (
  id                serial PRIMARY KEY,
  tid               varchar(127) NOT NULL DEFAULT '',
  sent_at           TIMESTAMP NOT NULL DEFAULT NOW(),
  msisdn            varchar(32) NOT NULL DEFAULT '',
  id_campaign       INT NOT NULL DEFAULT 0,
  id_service        INT NOT NULL DEFAULT 0,
  id_content        INT NOT NULL DEFAULT 0,
  id_subscription   INT NOT NULL DEFAULT 0,
  operator_code     INT NOT NULL DEFAULT 0,
  country_code      INT NOT NULL DEFAULT 0
);

CREATE TYPE user_action AS ENUM ('access', 'pull_click', 'content_get');
CREATE TABLE xmp_user_actions (
  id serial PRIMARY KEY,
  tid  varchar(127) NOT NULL DEFAULT '',
  action user_action NOT NULL,
  error varchar(511) NOT NULL DEFAULT '',
  access_at                   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE xmp_campaigns_access (
  id                          serial PRIMARY KEY,
  tid                         varchar(127) NOT NULL DEFAULT '',
  access_at                   TIMESTAMP NOT NULL DEFAULT NOW(),
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
