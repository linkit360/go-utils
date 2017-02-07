-- this is too long and also wrong.
-- select to_char(sent_at, 'YYYY-MM-DD'),
--   ( select count(*) FROM xmp_transactions where sent_at =  to_char(sent_at, 'YYYY-MM-DD') ) total_transactions,
-- 
--   ( select count(*) FROM xmp_transactions where sent_at =  to_char(sent_at, 'YYYY-MM-DD') and result in ('paid', 'retry_paid') ) total_paid_transactions,
-- 
--   ( select count(*) from (
--                        SELECT DISTINCT msisdn
--                        FROM xmp_transactions
--                        WHERE to_char(sent_at, 'YYYY-MM-DD')
--                      ) as t ) unique_users_count
-- from xmp_transactions
-- where sent_at > '2017-01-01' AND sent_at <= '2017-01-31'
-- group by to_char(sent_at, 'YYYY-MM-DD');

create type returntype as
(
  dt DATE,
  total_transactions int,
  total_paid_transactions int,
  unique_users_count int
);
-- drop FUNCTION tst_func()
-- drop type returntype CASCADE ;

CREATE OR REPLACE FUNCTION tst_func() returns setof returntype as
$BODY$
DECLARE
  declare
  r returntype%rowtype;
  monthDate DATE;
  monthDateTime timestamp;
BEGIN
  for monthDateTime in select created_at from test where  created_at >= '2017-01-01' AND created_at <= '2017-01-31'
  loop
    monthDate := CAST(monthDateTime as DATE);

    raise NOTICE 'date %s' , monthDate;
    r.dt = monthDate ;
    raise NOTICE 'total transactions.. ';
    select count(*) FROM xmp_transactions where CAST(sent_at AS DATE) = monthDate into r.total_transactions;
    raise NOTICE 'total paid transactions.. ';
    select count(*) FROM xmp_transactions where CAST(sent_at AS DATE) = monthDate and result in ('paid', 'retry_paid') into r.total_paid_transactions;
    raise NOTICE 'total unique msisdn transactions.. ';
    select count(*) from ( SELECT DISTINCT msisdn
                             FROM xmp_transactions
                             WHERE CAST(sent_at AS DATE) = monthDate
                         ) as t into r.unique_users_count;
    RETURN NEXT r;
  END LOOP;
  return;
END;
$BODY$
LANGUAGE plpgsql VOLATILE;


CREATE OR REPLACE FUNCTION tst_date_list_func() returns setof returntype as
$BODY$
DECLARE
  declare
  r returntype%rowtype;
  monthDate varchar;
BEGIN
  for monthDate in select created_at from test where  created_at >= '2017-01-01' AND created_at <= '2017-01-01'
  loop
    raise NOTICE 'date %s' , monthDate;
    r.dt = monthDate;
    RETURN NEXT r;
  END LOOP;
  return;
END;
$BODY$
LANGUAGE plpgsql VOLATILE;

select tst_date_list_func()




CREATE OR REPLACE FUNCTION stat_func(monthDate DATE) returns setof returntype as
$BODY$
DECLARE
  declare
  r returntype%rowtype;
BEGIN
    raise NOTICE 'date %s' , monthDate;
    r.dt = monthDate ;
    raise NOTICE 'total transactions.. ';
    select count(*) FROM xmp_transactions where CAST(sent_at AS DATE) = monthDate into r.total_transactions;
    raise NOTICE 'total paid transactions.. ';
    select count(*) FROM xmp_transactions where CAST(sent_at AS DATE) = monthDate and result in ('paid', 'retry_paid') into r.total_paid_transactions;
    raise NOTICE 'total unique msisdn transactions.. ';
    select count(*) from ( SELECT DISTINCT msisdn
                           FROM xmp_transactions
                           WHERE CAST(sent_at AS DATE) = monthDate
                         ) as t into r.unique_users_count;
    RETURN r;
END;
$BODY$
LANGUAGE plpgsql VOLATILE;
