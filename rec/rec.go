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
	Type                     string    `json:"type,omitempty"`
	Msisdn                   string    `json:"msisdn,omitempty"`
	Tid                      string    `json:"tid,omitempty"`
	Result                   string    `json:"result,omitempty"`
	SubscriptionStatus       string    `json:"subscription_status,omitempty"`
	OperatorCode             int64     `json:"operator_code,omitempty"`
	CountryCode              int64     `json:"country_code,omitempty"`
	ServiceId                int64     `json:"service_id,omitempty"`
	SubscriptionId           int64     `json:"subscription_is,omitempty"`
	CampaignId               int64     `json:"campaign_id,omitempty"`
	RetryId                  int64     `json:"retry_id,omitempty"`
	SentAt                   time.Time `json:"sent_at,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	LastPayAttemptAt         time.Time `json:"last_pay_attempt_at,omitempty"`
	AttemptsCount            int       `json:"attempts_count,omitempty"`
	KeepDays                 int       `json:"keep_days,omitempty"`
	DelayHours               int       `json:"delay_hours,omitempty"`
	PaidHours                int       `json:"paid_hours,omitempty"`
	OperatorName             string    `json:"operator_name,omitempty"`
	OperatorToken            string    `json:"operator_token,omitempty"`
	OperatorErr              string    `json:"opertor_err,omitempty"`
	Notice                   string    `json:"notice,omitempty"`
	Paid                     bool      `json:"paid,omitempty"`
	Price                    int       `json:"price,omitempty"`
	Pixel                    string    `json:"pixel,omitempty"`
	Publisher                string    `json:"publisher,omitempty"`
	SMSText                  string    `json:"sms_text,omitempty"`
	SMSSend                  bool      `json:"sms_send,omitempty"`
	Periodic                 bool      `json:"periodic,omitempty"`
	PeriodicDays             string    `json:"days,omitempty"`
	PeriodicAllowedFromHours int       `json:"allowed_from,omitempty"`
	PeriodicAllowedToHours   int       `json:"allowed_to,omitempty"`
}

func (r Record) TransactionOnly() bool {
	return r.Type == "injection" || r.Type == "expired"
}

var dbConn *sql.DB
var conf db.DataBaseConfig
var DBErrors m.Gauge
var Warn m.Gauge
var AddNewSubscriptionDuration prometheus.Summary

func Init(dbC db.DataBaseConfig) {
	log.SetLevel(log.DebugLevel)
	dbConn = db.Init(dbC)
	conf = dbC

	DBErrors = m.NewGauge("", "", "db_errors", "DB errors overall")
	Warn = m.NewGauge("", "", "warnings", "warnings overall")
	go func() {
		for range time.Tick(time.Minute) {
			DBErrors.Update()
			Warn.Update()
		}
	}()

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
func GetRetryTransactions(operatorCode int64, batchLimit int, paidOnceHours int) ([]Record, error) {
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
				"query":         query,
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

	notPaidInHours := ""
	if paidOnceHours > 0 {
		notPaidInHours = fmt.Sprintf(" AND msisdn NOT IN ("+
			" SELECT DISTINCT msisdn "+
			" FROM %stransactions "+
			" WHERE sent_at > (CURRENT_TIMESTAMP -  INTERVAL '%d hours' ) AND "+
			"       ( result = 'paid' OR result = 'retry_paid') )",
			conf.TablePrefix,
			paidOnceHours)
	}

	query = fmt.Sprintf("SELECT "+
		"msisdn, "+
		"id, "+
		"tid, "+
		"created_at, "+
		"last_pay_attempt_at, "+
		"attempts_count, "+
		"keep_days, "+
		"delay_hours, "+
		"price, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		"FROM %sretries "+
		"WHERE "+
		" operator_code = $1 AND "+
		" status = ''  "+notPaidInHours+
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
			&record.Msisdn,
			&record.RetryId,
			&record.Tid,
			&record.CreatedAt,
			&record.LastPayAttemptAt,
			&record.AttemptsCount,
			&record.KeepDays,
			&record.DelayHours,
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
		"updated_at = $2, "+
		"attempts_count = attempts_count + 1 "+ // for retry sent consent
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

func LoadActiveSubscriptions(hours int) (records []ActiveSubscription, err error) {
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
		"result IN ('', 'paid', 'failed') ",
		conf.TablePrefix)

	prev := []ActiveSubscription{}
	rows, err := dbConn.Query(query)
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
	if r.CampaignId == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no campaign id")
		Warn.Inc()
	}
	if r.ServiceId == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no service id")
		Warn.Inc()
	}
	if r.DelayHours == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no delay hours")
		Warn.Inc()
	}
	if r.KeepDays == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no keep days")
		Warn.Inc()
	}
	if r.PaidHours == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no paid hours")
		Warn.Inc()
	}
	if r.OperatorCode == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no operator code")
		Warn.Inc()
	}
	if r.CountryCode == 0 {
		log.WithFields(log.Fields{
			"tid": r.Tid,
		}).Warn("no country code")
		Warn.Inc()
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

// bare periodic
func GetPeriodics(batchLimit, repeaIntervalMinutes int, intervalType string, loc *time.Location) (records []Record, err error) {
	begin := time.Now()
	query := ""
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":         time.Since(begin),
				"intervalType": intervalType,
				"loc":          loc.String(),
				"query":        query,
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
	var interval string
	if intervalType == "hour" {
		interval = time.Now().In(loc).Format("15")
	} else if intervalType == "min" {
		now := time.Now().In(loc)
		interval = strconv.Itoa(60*now.Hour() + now.Minute())
	} else {
		err = fmt.Errorf("Unknown interval Type: %s", intervalType)
		return
	}
	log.WithFields(log.Fields{
		"interval": interval,
		"day":      dayName,
	}).Debug("time params")

	inSpecifiedTime := "( allowed_from <= " + interval + " AND  allowed_to >= " + interval + " ) "

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
		"periodic = true AND "+inSpecifiedTime+"AND ("+
		// paid not processed today
		"  (  ( days ? '"+dayName+"' OR days ? 'any' ) AND "+
		"     result = 'paid' AND "+
		"     last_pay_attempt_at < (CURRENT_TIMESTAMP -  INTERVAL '24 hours' ) "+
		"  ) "+
		"  OR  "+
		//not paid processed today (including '' and failed)
		"  (  ( days ? '"+dayName+"' OR days ? 'any' ) AND "+
		"     result NOT IN ('rejected', 'paid', 'postpaid', 'pending' ) AND "+
		"     last_pay_attempt_at < (CURRENT_TIMESTAMP -  %d * INTERVAL '1 minute' ) "+
		"  ) "+
		" )"+ // close inspecified time range AND
		"ORDER BY last_pay_attempt_at ASC LIMIT %s", // get the last touched
		conf.TablePrefix,
		repeaIntervalMinutes,
		strconv.Itoa(batchLimit),
	)

	rows, err := dbConn.Query(query)
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

// retries for periodic?
func GetNotPaidPeriodics(delay_minutes, batchLimit int) (records []Record, err error) {
	begin := time.Now()
	query := ""
	defer func() {
		defer func() {
			fields := log.Fields{
				"took": time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load not paid periodic failed")
			} else {
				//fields["query"] = query
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load not paid periodic")
			}
		}()
	}()

	dayName := strings.ToLower(time.Now().Format("Mon"))

	var periodics []Record
	matchedToday := "( days ? '" + dayName + "' OR days ? 'any' ) "

	earlierMatchedToday := "( last_pay_attempt_at + INTERVAL '24 hours' < NOW() AND " +
		"result NOT IN ('rejected', 'blacklisted', 'canceled', 'pending') )"

	notPaidAtAllMatchedTime := "( last_pay_attempt_at + INTERVAL '24 hours' > NOW() AND " +
		"result NOT IN ('rejected', 'blacklisted', 'canceled', 'pending', 'paid') " +
		fmt.Sprintf("AND last_pay_attempt_at + %d * INTERVAL '1 minute' < NOW() )", delay_minutes)

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
		matchedToday+" AND periodic = true AND "+
		"("+earlierMatchedToday+" OR "+notPaidAtAllMatchedTime+")"+
		"ORDER BY last_pay_attempt_at ASC LIMIT %s",
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)

	rows, err := dbConn.Query(query)
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
func GetSubscriptionByToken(token string) (p Record, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":  time.Since(begin),
				"token": token,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("get subscription by token failed")
			} else {
				log.WithFields(fields).Debug("get subscription by token")
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
func GetSubscriptionByMsisdn(msisdn string) (p Record, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"took":   time.Since(begin),
				"msisdn": msisdn,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("get subscription by msisdn failed")
			} else {
				log.WithFields(fields).Debug("get subscription by msisdn")
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
		"WHERE msisdn = $1 LIMIT 1",
		conf.TablePrefix,
	)

	if err = dbConn.QueryRow(query, msisdn).Scan(
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
		if err == sql.ErrNoRows {
			return
		}
		DBErrors.Inc()
		return Record{}, fmt.Errorf("Rows.Next: %s", err.Error())
	}
	return p, nil
}

func GetRetryByMsisdn(msisdn, status string) (r Record, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"msisdn": msisdn,
				"took":   time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load retry failed")
			} else {
				fields["att_count"] = r.AttemptsCount
				fields["tid"] = r.Tid
				log.WithFields(fields).Debug("loaded retry")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT "+
		"msisdn, "+
		"id, "+
		"tid, "+
		"created_at, "+
		"last_pay_attempt_at, "+
		"attempts_count, "+
		"keep_days, "+
		"delay_hours, "+
		"price, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		"FROM %sretries "+
		"WHERE "+
		" msisdn = $1 AND status = $2"+
		" ORDER BY id "+
		" LIMIT 1", // get the oldest retry
		conf.TablePrefix,
	)

	if err := dbConn.QueryRow(query, msisdn, status).Scan(
		&r.Msisdn,
		&r.RetryId,
		&r.Tid,
		&r.CreatedAt,
		&r.LastPayAttemptAt,
		&r.AttemptsCount,
		&r.KeepDays,
		&r.DelayHours,
		&r.Price,
		&r.OperatorCode,
		&r.CountryCode,
		&r.ServiceId,
		&r.SubscriptionId,
		&r.CampaignId,
	); err != nil {
		if err != sql.ErrNoRows {
			DBErrors.Inc()
		}
		return Record{}, err // do not change type of error, please, it's being checked further
	}

	return
}

func GetBufferPixelByCampaignId(campaignId int64) (r Record, err error) {
	begin := time.Now()
	defer func() {
		defer func() {
			fields := log.Fields{
				"campaign_id": campaignId,
				"took":        time.Since(begin),
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load buffer pixel failed")
			} else {
				fields["tid"] = r.Tid
				log.WithFields(fields).Debug("loaded buffer pixel")
			}
		}()
	}()

	query := fmt.Sprintf("SELECT "+
		"sent_at, "+
		"id_campaign, "+
		"tid, "+
		"pixel "+
		"FROM %spixel_buffer "+
		"WHERE "+
		" id_campaign = $1 "+
		" ORDER BY id "+
		" LIMIT 1", // get the oldest retry
		conf.TablePrefix,
	)

	if err := dbConn.QueryRow(query, campaignId).Scan(
		&r.SentAt,
		&r.CampaignId,
		&r.Tid,
		&r.Pixel,
	); err != nil {
		if err != sql.ErrNoRows {
			DBErrors.Inc()
		}
		return Record{}, err // do not change type of error, please, it's being checked further
	}

	return
}

func GetRepeatSentConsent(operatorCode int64, delayMinutes, batchLimit int) (records []Record, err error) {
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
				log.WithFields(fields).Error("load sent consent repeat failed")
			} else {
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load sent consent repeat")
			}
		}()
	}()

	dayName := strings.ToLower(time.Now().Format("Mon"))

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
		"( days ? '"+dayName+"' OR days ? 'any' ) AND "+
		"result = 'consent' AND attempts_count = 0 AND "+
		"updated_at < (CURRENT_TIMESTAMP -  %d * INTERVAL '1 minutes' ) "+
		"ORDER BY last_pay_attempt_at ASC LIMIT %s",
		conf.TablePrefix,
		delayMinutes,
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
			return []Record{}, fmt.Errorf("rows.Scan: %s", err.Error())
		}

		periodics = append(periodics, p)
	}
	if rows.Err() != nil {
		DBErrors.Inc()

		err = fmt.Errorf("rows.Err: %s", err.Error())
		return []Record{}, err
	}
	return periodics, nil
}

func GetNotSentPixels(hours, limit int) (records []Record, err error) {
	defer func() {
		defer func() {
			fields := log.Fields{
				"hours": hours,
				"limit": limit,
			}
			if err != nil {
				fields["error"] = err.Error()
				log.WithFields(fields).Error("load not sent pixels failed")
			} else {
				fields["count"] = len(records)
				log.WithFields(fields).Debug("load not sent pixels")
			}
		}()
	}()

	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"took": time.Since(begin),
		}).Debug("get pixels")
	}()
	query := fmt.Sprintf("SELECT "+
		"tid, "+
		"msisdn, "+
		"id_campaign, "+
		"id, "+
		"operator_code, "+
		"country_code, "+
		"pixel, "+
		"publisher "+
		" FROM %ssubscriptions "+
		" WHERE pixel != '' "+
		" AND pixel_sent = false "+
		"AND result NOT IN ('', 'postpaid', 'blacklisted', 'rejected', 'canceled')",
		conf.TablePrefix)

	if hours > 0 {
		query = query +
			fmt.Sprintf(" AND (CURRENT_TIMESTAMP - %d * INTERVAL '1 hour' ) > sent_at ", hours)
	}
	query = query + fmt.Sprintf(" ORDER BY id ASC LIMIT %d", limit)

	rows, err := dbConn.Query(query)
	if err != nil {
		err = fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
		return
	}
	defer rows.Close()

	for rows.Next() {
		record := Record{}

		if err = rows.Scan(
			&record.Tid,
			&record.Msisdn,
			&record.CampaignId,
			&record.SubscriptionId,
			&record.OperatorCode,
			&record.CountryCode,
			&record.Pixel,
			&record.Publisher,
		); err != nil {
			err = fmt.Errorf("rows.Scan: %s", err.Error())
			return
		}
		records = append(records, record)
	}
	if rows.Err() != nil {
		err = fmt.Errorf("row.Err: %s", err.Error())
		return
	}
	return
}
