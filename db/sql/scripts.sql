-- access campaign
alter table xmp_campaigns_access add column sent_at                     TIMESTAMP NOT NULL DEFAULT NOW();
update xmp_campaigns_access set sent_at = access_at ;

-- content sent ok
