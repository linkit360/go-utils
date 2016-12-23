select count(*) from xmp_subscriptions
where date_trunc('day',  created_at)  = '2016-12-04 00:00:00.000000';
-- 61590

select result, count(*) from xmp_subscriptions
where date_trunc('day',  created_at)  = '2016-12-04 00:00:00.000000' group by result;


select count(*) from xmp_user_actions
where date_trunc('day',  access_at)  = '2016-12-04 00:00:00.000000'
and action = 'pull_click';
-- 66698

select count(*) from xmp_campaigns_access
where
  length(msisdn) > 5 AND tid in (
  SELECT tid
  FROM xmp_user_actions
  WHERE date_trunc('day', access_at) = '2016-12-04 00:00:00.000000' AND action = 'pull_click'
);
-- 61864