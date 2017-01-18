package rec

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/vostrok/utils/db"
	m "github.com/vostrok/utils/metrics"
)

// please, do not add any json named field in old field,
// bcz unmarshalling will brake the flow
type Record struct {
	Msisdn                   string    `json:",omitempty"`
	Tid                      string    `json:",omitempty"`
	Result                   string    `json:",omitempty"`
	SubscriptionStatus       string    `json:",omitempty"`
	OperatorCode             int64     `json:",omitempty"`
	CountryCode              int64     `json:",omitempty"`
	ServiceId                int64     `json:",omitempty"`
	SubscriptionId           int64     `json:",omitempty"`
	CampaignId               int64     `json:",omitempty"`
	RetryId                  int64     `json:",omitempty"`
	SentAt                   time.Time `json:",omitempty"`
	CreatedAt                time.Time `json:",omitempty"`
	LastPayAttemptAt         time.Time `json:",omitempty"`
	AttemptsCount            int       `json:",omitempty"`
	KeepDays                 int       `json:",omitempty"`
	DelayHours               int       `json:",omitempty"`
	PaidHours                int       `json:",omitempty"`
	OperatorName             string    `json:",omitempty"`
	OperatorToken            string    `json:",omitempty"`
	OperatorErr              string    `json:",omitempty"`
	Paid                     bool      `json:",omitempty"`
	Price                    int       `json:",omitempty"`
	Pixel                    string    `json:",omitempty"`
	Publisher                string    `json:",omitempty"`
	SMSText                  string    `json:",omitempty"`
	SMSSend                  bool      `json:",omitempty"`
	Periodic                 bool      `json:"periodic,omitempty"`
	PeriodicDays             string    `json:"days,omitempty"`
	PeriodicAllowedFromHours int       `json:"allowed_from,omitempty"`
	PeriodicAllowedToHours   int       `json:"allowed_to,omitempty"`
}

var dbConn *sql.DB
var conf db.DataBaseConfig
var DBErrors m.Gauge
var AddNewSubscriptionDuration prometheus.Summary

func Init(dbC db.DataBaseConfig) {
	log.SetLevel(log.DebugLevel)
	dbConn = db.Init(dbC)
	conf = dbC

	DBErrors = m.NewGauge("", "", "db_errors", "DB errors pverall mt_manager")
	AddNewSubscriptionDuration = m.NewSummary("subscription_add_to_db_duration_seconds", "new subscription add duration")
}
func GenerateTID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.WithField("error", err.Error()).Error("generate uniq id")
	}
	tid := fmt.Sprintf("%d-%s", time.Now().Unix(), u4)
	log.WithField("tid", tid).Debug("generated tid")
	return tid
}
func GetSuspendedRetriesCount() (count int, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("get suspended retries count failed")
			} else {
				fields["count"] = count
				log.WithFields(fields).Debug("get suspended retries")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT count(*) count FROM %sretries "+
		"WHERE status IN ( 'pending', 'script' ) "+
		"AND updated_at < (CURRENT_TIMESTAMP - 4 * INTERVAL '1 hour' ) ",
		conf.TablePrefix,
	)
	rows, err := dbConn.Query(query)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&count,
		); err != nil {
			DBErrors.Inc()

			err = fmt.Errorf("rows.Scan: %s", err.Error())
			return count, err
		}
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("get pending retries: rows.Err: %s", err.Error())
		return count, err
	}
	return count, nil
}

func GetRetriesPeriod() (seconds int64, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("get retries period failed")
			} else {
				log.WithFields(fields).Debug("get retries period")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT coalesce((SELECT extract (epoch from (now() - "+
		"MIN(last_pay_attempt_at))::interval) seconds from %sretries), 0)",
		conf.TablePrefix,
	)
	rows, err := dbConn.Query(query)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&seconds,
		); err != nil {
			DBErrors.Inc()

			err = fmt.Errorf("rows.Scan: %s", err.Error())
			return
		}
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("rows.Err: %s", err.Error())
		return
	}

	return
}

func GetSuspendedSubscriptionsCount() (count int, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("get mo count failed")
			} else {
				fields["count"] = count
				log.WithFields(fields).Debug("get mo count")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT count(*) count FROM %ssubscriptions "+
		"WHERE result = ''"+
		"AND sent_at < (CURRENT_TIMESTAMP - 2 * INTERVAL '1 hour' ) ",
		conf.TablePrefix,
	)
	rows, err := dbConn.Query(query)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&count,
		); err != nil {
			DBErrors.Inc()

			err = fmt.Errorf("rows.Scan: %s", err.Error())
			return count, err
		}
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("get pending subscriptions: rows.Err: %s", err.Error())
		return count, err
	}

	return count, nil
}
func GetRetryTransactions(operatorCode int64, batchLimit int) ([]Record, error) {
	begin := time.Now()
	var retries []Record
	var err error
	var query string
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":          time.Since(begin),
				"operator_code": operatorCode,
				"limit":         batchLimit,
				//"query":         query,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load retries failed")
			} else {
				fields["count"] = len(retries)
				log.WithFields(fields).Debug("load retries")
			}
		}()
	}()

	query = fmt.Sprintf("SELECT "+
		"id, "+
		"tid, "+
		"created_at, "+
		"last_pay_attempt_at, "+
		"attempts_count, "+
		"keep_days, "+
		"delay_hours, "+
		"msisdn, "+
		"price, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		"FROM %sretries "+
		"WHERE "+
		" operator_code = $1 AND "+
		" status = '' "+
		" ORDER BY last_pay_attempt_at ASC "+
		" LIMIT %s", // get the last touched
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)

	rows, err := dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return []Record{}, err
	}
	defer rows.Close()

	for rows.Next() {
		record := Record{}
		if err := rows.Scan(
			&record.RetryId,
			&record.Tid,
			&record.CreatedAt,
			&record.LastPayAttemptAt,
			&record.AttemptsCount,
			&record.KeepDays,
			&record.DelayHours,
			&record.Msisdn,
			&record.Price,
			&record.OperatorCode,
			&record.CountryCode,
			&record.ServiceId,
			&record.SubscriptionId,
			&record.CampaignId,
		); err != nil {
			DBErrors.Inc()
			return []Record{}, fmt.Errorf("Rows.Next: %s", err.Error())
		}

		retries = append(retries, record)
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("GetRetries RowsError: %s", err.Error())
		return []Record{}, err
	}
	return retries, nil
}
func SetSubscriptionStatus(status string, id int64) (err error) {
	if id == 0 {
		log.WithFields(log.Fields{"error": "no subscription id"}).Error("set periodic status")
		return nil
	}
	begin := time.Now()
	defer func() {
		fields := log.Fields{
			"status": status,
			"id":     id,
			"took":   time.Since(begin),
		}
		if err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("set periodic result failed")
		} else {
			log.WithFields(fields).Debug("set periodic result")
		}
	}()

	query := fmt.Sprintf("UPDATE %ssubscriptions SET "+
		"result = $1, "+
		"updated_at = $2 "+
		"WHERE id = $3",
		conf.TablePrefix,
	)

	updatedAt := time.Now().UTC()
	_, err = dbConn.Exec(query, status, updatedAt, id)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("dbConn.Exec: %s, Query: %s", err.Error(), query)
		return err
	}
	return nil
}

func SetRetryStatus(status string, id int64) (err error) {
	if id == 0 {
		log.WithFields(log.Fields{"error": "no retry id"}).Error("set retry status")
		return nil
	}
	begin := time.Now()
	defer func() {
		fields := log.Fields{
			"status": status,
			"id":     id,
			"took":   time.Since(begin),
		}
		if err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("set retry status failed")
		} else {
			log.WithFields(fields).Debug("set retry status")
		}
	}()

	query := fmt.Sprintf("UPDATE %sretries SET "+
		"status = $1, "+
		"updated_at = $2 "+
		"WHERE id = $3", conf.TablePrefix)

	updatedAt := time.Now().UTC()
	_, err = dbConn.Exec(query, status, updatedAt, id)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("dbConn.Exec: %s, Query: %s", err.Error(), query)
		return err
	}
	return nil
}

func LoadScriptRetries(hoursPassed int, operatorCode int64, batchLimit int) (records []Record, err error) {
	var retries []Record
	query := ""
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":  time.Since(begin),
				"hours": hoursPassed,
				"limit": batchLimit,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load retries failed")
			} else {
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load retries")
			}
		}()
	}()
	query = fmt.Sprintf("SELECT "+
		"id, "+
		"tid, "+
		"created_at, "+
		"last_pay_attempt_at, "+
		"attempts_count, "+
		"keep_days, "+
		"delay_hours, "+
		"msisdn, "+
		"price, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		"FROM %sretries "+
		"WHERE "+
		"operator_code = $1 AND "+
		"status IN ( 'pending', 'script' ) AND "+
		"updated_at < (CURRENT_TIMESTAMP - 5 * INTERVAL '1 minute' ) "+
		"ORDER BY last_pay_attempt_at ASC LIMIT %s", // get the last touched
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)
	rows, err := dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()
		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return []Record{}, err
	}
	defer rows.Close()

	for rows.Next() {
		record := Record{}
		if err := rows.Scan(
			&record.RetryId,
			&record.Tid,
			&record.CreatedAt,
			&record.LastPayAttemptAt,
			&record.AttemptsCount,
			&record.KeepDays,
			&record.DelayHours,
			&record.Msisdn,
			&record.Price,
			&record.OperatorCode,
			&record.CountryCode,
			&record.ServiceId,
			&record.SubscriptionId,
			&record.CampaignId,
		); err != nil {
			DBErrors.Inc()

			err = fmt.Errorf("Rows.Next: %s", err.Error())
			return []Record{}, err
		}

		retries = append(retries, record)
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("rows.Error: %s", err.Error())
		return []Record{}, err
	}
	return retries, nil
}

type ActiveSubscription struct {
	Id         int64
	CreatedAt  time.Time
	Msisdn     string
	ServiceId  int64
	CampaignId int64
}

func LoadActiveSubscriptions(operatorCode int64, hours int) (records []ActiveSubscription, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load active subscriptions failed")
			} else {
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load active subscriptions ")
			}
		}()
	}()

	hoursPassed := ""
	if hours > 0 {
		hoursPassed = fmt.Sprintf("(CURRENT_TIMESTAMP - %d * INTERVAL '1 hour' ) < created_at AND ", hours)
	}
	query := fmt.Sprintf("SELECT "+
		"id, "+
		"msisdn, "+
		"id_service, "+
		"id_campaign, "+
		"created_at "+
		"FROM %ssubscriptions "+
		"WHERE "+
		hoursPassed+
		"result IN ('', 'paid', 'failed') AND "+
		"operator_code = $1",
		conf.TablePrefix)

	prev := []ActiveSubscription{}
	rows, err := dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return prev, err
	}
	defer rows.Close()

	for rows.Next() {
		var p ActiveSubscription
		if err := rows.Scan(
			&p.Id,
			&p.Msisdn,
			&p.ServiceId,
			&p.CampaignId,
			&p.CreatedAt,
		); err != nil {
			DBErrors.Inc()

			err = fmt.Errorf("Rows.Next: %s", err.Error())
			return prev, err
		}
		prev = append(prev, p)
	}

	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("Rows.Err: %s", err.Error())
		return prev, err
	}
	return prev, nil
}

func AddNewSubscriptionToDB(r *Record) error {
	if r.SubscriptionId > 0 {
		log.WithFields(log.Fields{
			"tid":    r.Tid,
			"id":     r.SubscriptionId,
			"msisdn": r.Msisdn,
		}).Debug("already has subscription id")
		return nil
	}
	if len(r.Msisdn) > 32 {
		log.WithFields(log.Fields{
			"tid":    r.Tid,
			"msisdn": r.Msisdn,
			"error":  "too long msisdn",
		}).Error("strange msisdn, truncating")
		r.Msisdn = r.Msisdn[:31]
	}
	if r.PeriodicDays == "" {
		r.PeriodicDays = "[]"
	}
	begin := time.Now()
	query := fmt.Sprintf("INSERT INTO %ssubscriptions ( "+
		"sent_at, "+
		"result, "+
		"id_campaign, "+
		"id_service, "+
		"msisdn, "+
		"publisher, "+
		"pixel, "+
		"tid, "+
		"operator_token, "+
		"country_code, "+
		"operator_code, "+
		"paid_hours, "+
		"delay_hours, "+
		"keep_days, "+
		"price, "+
		"periodic, "+
		"days,"+
		"allowed_from,"+
		"allowed_to"+
		") values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,"+
		" $11, $12, $13, $14, $15, $16, $17, $18, $19) "+
		"RETURNING id",
		conf.TablePrefix,
	)

	if err := dbConn.QueryRow(query,
		r.SentAt,
		"",
		r.CampaignId,
		r.ServiceId,
		r.Msisdn,
		r.Publisher,
		r.Pixel,
		r.Tid,
		r.OperatorToken,
		r.CountryCode,
		r.OperatorCode,
		r.PaidHours,
		r.DelayHours,
		r.KeepDays,
		r.Price,
		r.Periodic,
		r.PeriodicDays,
		r.PeriodicAllowedFromHours,
		r.PeriodicAllowedToHours,
	).Scan(&r.SubscriptionId); err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Scan: %s", err.Error())
		log.WithFields(log.Fields{
			"tid":   r.Tid,
			"error": err.Error(),
			"query": query,
			"msg":   "requeue",
		}).Error("add new subscription")
		return err
	}
	AddNewSubscriptionDuration.Observe(time.Since(begin).Seconds())
	log.WithFields(log.Fields{
		"tid":  r.Tid,
		"id":   r.SubscriptionId,
		"took": time.Since(begin).Seconds(),
	}).Info("added new subscription")
	return nil
}

func GetSuspendedSubscriptions(operatorCode int64, hours, limit int) (records []Record, err error) {
	query := fmt.Sprintf("SELECT "+
		"id, "+
		"tid, "+
		"msisdn, "+
		"pixel, "+
		"publisher, "+
		"id_service, "+
		"id_campaign, "+
		"operator_code, "+
		"country_code, "+
		"attempts_count, "+
		"delay_hours, "+
		"paid_hours, "+
		"keep_days, "+
		"price "+
		" FROM %ssubscriptions "+
		" WHERE result = '' AND "+
		"operator_code = $1 AND "+
		" (CURRENT_TIMESTAMP - %d * INTERVAL '1 hour' ) > created_at "+
		" ORDER BY id ASC LIMIT %s",
		conf.TablePrefix,
		hours,
		strconv.Itoa(limit),
	)
	var rows *sql.Rows
	rows, err = dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()
		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return
	}
	defer rows.Close()

	for rows.Next() {
		record := Record{}

		if err = rows.Scan(
			&record.SubscriptionId,
			&record.Tid,
			&record.Msisdn,
			&record.Pixel,
			&record.Publisher,
			&record.ServiceId,
			&record.CampaignId,
			&record.OperatorCode,
			&record.CountryCode,
			&record.AttemptsCount,
			&record.DelayHours,
			&record.PaidHours,
			&record.KeepDays,
			&record.Price,
		); err != nil {
			DBErrors.Inc()
			err = fmt.Errorf("rows.Scan: %s", err.Error())
			return
		}
		records = append(records, record)
	}
	if rows.Err() != nil {
		DBErrors.Inc()
		err = fmt.Errorf("rows.Err: %s", err.Error())
		return
	}
	return
}

func GetPeriodics(operatorCode int64, batchLimit int) (records []Record, err error) {
	begin := time.Now()
	query := ""
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":         time.Since(begin),
				"query":        query,
				"operatorCode": operatorCode,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load periodic failed")
			} else {
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load periodic")
			}
		}()
	}()

	dayName := strings.ToLower(time.Now().Format("Mon"))
	hourNow := time.Now().Format("15")
	inSpecifiedHours := "( allowed_from <= " + hourNow + " AND  allowed_to >= " + hourNow + " ) "

	var periodics []Record

	query = fmt.Sprintf("SELECT "+
		"id, "+
		"sent_at, "+
		"tid , "+
		"operator_token, "+
		"price, "+
		"id_service, "+
		"id_campaign, "+
		"country_code, "+
		"operator_code, "+
		"msisdn, "+
		"keep_days, "+
		"delay_hours, "+
		"paid_hours "+
		"FROM %ssubscriptions "+
		"WHERE "+
		"operator_code = $1 AND periodic = true AND "+
		"( days ? '"+dayName+"' OR days ? 'any' ) AND "+ // today
		inSpecifiedHours+"AND "+
		"result NOT IN ('rejected', 'blacklisted', 'canceled', 'pending') AND "+ // not cancelled, not rejected, not blacklisted
		"last_pay_attempt_at < (CURRENT_TIMESTAMP -  INTERVAL '18 hours' ) "+ // not processed today
		"ORDER BY last_pay_attempt_at ASC LIMIT %s", // get the last touched
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)

	rows, err := dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return []Record{}, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Record{}
		if err := rows.Scan(
			&p.SubscriptionId,
			&p.SentAt,
			&p.Tid,
			&p.OperatorToken,
			&p.Price,
			&p.ServiceId,
			&p.CampaignId,
			&p.OperatorCode,
			&p.CountryCode,
			&p.Msisdn,
			&p.KeepDays,
			&p.DelayHours,
			&p.PaidHours,
		); err != nil {
			DBErrors.Inc()
			return []Record{}, fmt.Errorf("Rows.Next: %s", err.Error())
		}

		periodics = append(periodics, p)
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("GetPeriodic RowsError: %s", err.Error())
		return []Record{}, err
	}
	return periodics, nil
}

func GetPeriodicSubscriptionByToken(token string) (p Record, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load periodic cache failed")
			} else {
				log.WithFields(fields).Debug("load periodic cache")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT "+
		"id, "+
		"sent_at, "+
		"tid , "+
		"operator_token, "+
		"price, "+
		"id_service, "+
		"id_campaign, "+
		"country_code, "+
		"operator_code, "+
		"msisdn, "+
		"keep_days, "+
		"delay_hours, "+
		"paid_hours "+
		"FROM %ssubscriptions "+
		"WHERE operator_token = $1 LIMIT 1",
		conf.TablePrefix,
	)

	rows, err := dbConn.Query(query, token)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return Record{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&p.SubscriptionId,
			&p.SentAt,
			&p.Tid,
			&p.OperatorToken,
			&p.Price,
			&p.ServiceId,
			&p.CampaignId,
			&p.OperatorCode,
			&p.CountryCode,
			&p.Msisdn,
			&p.KeepDays,
			&p.DelayHours,
			&p.PaidHours,
		); err != nil {
			DBErrors.Inc()
			return Record{}, fmt.Errorf("Rows.Next: %s", err.Error())
		}
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("GetPeriodic RowsError: %s", err.Error())
		return Record{}, err
	}
	return p, nil
}
