package rec

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/vostrok/utils/db"
	m "github.com/vostrok/utils/metrics"
	"strings"
)

type Record struct {
	Msisdn             string    `json:",omitempty"`
	Tid                string    `json:",omitempty"`
	Result             string    `json:",omitempty"`
	SubscriptionStatus string    `json:",omitempty"`
	OperatorCode       int64     `json:",omitempty"`
	CountryCode        int64     `json:",omitempty"`
	ServiceId          int64     `json:",omitempty"`
	SubscriptionId     int64     `json:",omitempty"`
	CampaignId         int64     `json:",omitempty"`
	RetryId            int64     `json:",omitempty"`
	CreatedAt          time.Time `json:",omitempty"`
	LastPayAttemptAt   time.Time `json:",omitempty"`
	AttemptsCount      int       `json:",omitempty"`
	KeepDays           int       `json:",omitempty"`
	DelayHours         int       `json:",omitempty"`
	OperatorName       string    `json:",omitempty"`
	OperatorToken      string    `json:",omitempty"`
	OperatorErr        string    `json:",omitempty"`
	Paid               bool      `json:",omitempty"`
	Price              int       `json:",omitempty"`
	Pixel              string    `json:",omitempty"`
	Publisher          string    `json:",omitempty"`
	SMSText            string    `json:",omitempty"`
}

var dbConn *sql.DB
var conf db.DataBaseConfig
var DBErrors m.Gauge

func Init(dbC db.DataBaseConfig) {
	log.SetLevel(log.DebugLevel)
	dbConn = db.Init(dbC)
	conf = dbC

	DBErrors = m.NewGauge("", "", "db_errors", "DB errors pverall mt_manager")
}

func GetNotPaidSubscriptions(batchLimit int) ([]Record, error) {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"took": time.Since(begin),
		}).Debug("get notpaid subscriptions")
	}()
	var subscr []Record
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
		"attempts_count "+
		" FROM %ssubscriptions "+
		" WHERE result = '' "+
		" ORDER BY id DESC LIMIT %s",
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)
	rows, err := dbConn.Query(query)
	if err != nil {
		DBErrors.Inc()
		return subscr, fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
	}
	defer rows.Close()

	for rows.Next() {
		record := Record{}

		if err := rows.Scan(
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
		); err != nil {
			DBErrors.Inc()
			return subscr, err
		}
		subscr = append(subscr, record)
	}
	if rows.Err() != nil {
		DBErrors.Inc()
		return subscr, fmt.Errorf("row.Err: %s", err.Error())
	}

	return subscr, nil
}

func GetRetryTransactions(operatorCode int64, batchLimit int) ([]Record, error) {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"took": time.Since(begin),
		}).Debug("get retry transactions")
	}()

	var retries []Record
	query := fmt.Sprintf("SELECT "+
		"id, "+
		"tid, "+
		"created_at, "+
		"last_pay_attempt_at, "+
		"attempts_count, "+
		"keep_days, "+
		"msisdn, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		"FROM %sretries "+
		"WHERE (CURRENT_TIMESTAMP - delay_hours * INTERVAL '1 hour' ) > last_pay_attempt_at AND "+
		"operator_code = $1 AND "+
		"status = '' "+
		"ORDER BY last_pay_attempt_at ASC LIMIT %s", // get the last touched
		conf.TablePrefix,
		strconv.Itoa(batchLimit),
	)
	rows, err := dbConn.Query(query, operatorCode)
	if err != nil {
		DBErrors.Inc()
		return []Record{}, fmt.Errorf("db.Query: %s, query: %s", err.Error(), query)
	}
	defer rows.Close()

	retryIds := []interface{}{}
	for rows.Next() {
		record := Record{}
		if err := rows.Scan(
			&record.RetryId,
			&record.Tid,
			&record.CreatedAt,
			&record.LastPayAttemptAt,
			&record.AttemptsCount,
			&record.KeepDays,
			&record.Msisdn,
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
		retryIds = append(retryIds, record.RetryId)
	}
	if rows.Err() != nil {
		DBErrors.Inc()
		return []Record{}, fmt.Errorf("GetRetries RowsError: %s", err.Error())
	}

	if err := setPendingStatus(retryIds); err != nil {
		return []Record{}, err
	}
	return retries, nil
}
func setPendingStatus(ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"took":  time.Since(begin),
			"count": len(ids),
		}).Debug("update pending retries")
	}()

	paramStr := []string{}
	for i, _ := range ids {
		j := i + 1
		paramStr = append(paramStr, "$"+strconv.Itoa(j))
	}
	query := fmt.Sprintf("UPDATE %sretries SET "+
		"status = 'pending' "+
		"WHERE id IN ("+strings.Join(paramStr, ", ")+")", conf.TablePrefix)

	_, err := dbConn.Exec(query, ids...)
	if err != nil {
		DBErrors.Inc()

		err = fmt.Errorf("dbConn.Exec: %s", err.Error())
		log.WithFields(log.Fields{
			"error ": err.Error(),
			"query":  query,
		}).Error("update retry failed")
		return err
	}
	return nil
}

type PreviuosSubscription struct {
	Id        int64
	CreatedAt time.Time
}

// CREATE INDEX xmp_subscriptions_prev_idx ON xmp_subscriptions (id, msisdn, id_service, created_at);
func (t Record) GetPreviousSubscription(paidHours int) (PreviuosSubscription, error) {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":       t.Tid,
			"paidHours": paidHours,
			"took":      time.Since(begin),
		}).Debug("get previous subscription")
	}()

	// get the very old first, but not elder than paidHours ago
	query := fmt.Sprintf("SELECT "+
		"id, "+
		"created_at "+
		"FROM %ssubscriptions "+
		"WHERE id < $1 AND "+
		"msisdn = $2 AND id_service = $3 AND "+
		"(CURRENT_TIMESTAMP - "+strconv.Itoa(paidHours)+" * INTERVAL '1 hour' ) < created_at "+
		"ORDER BY created_at ASC LIMIT 1",
		conf.TablePrefix)

	var p PreviuosSubscription
	if err := dbConn.QueryRow(query,
		t.SubscriptionId,
		t.Msisdn,
		t.ServiceId,
	).Scan(
		&p.Id,
		&p.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return p, err
		}
		DBErrors.Inc()
		return p, fmt.Errorf("db.QueryRow: %s, query: %s", err.Error(), query)
	}
	return p, nil
}
func (t Record) WriteTransaction() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":    t.Tid,
			"took":   time.Since(begin),
			"result": t.Result,
		}).Debug("write transaction")
	}()
	query := fmt.Sprintf("INSERT INTO %stransactions ("+
		"tid, "+
		"msisdn, "+
		"result, "+
		"operator_code, "+
		"country_code, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign, "+
		"operator_token, "+
		"price "+
		") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		conf.TablePrefix)

	_, err := dbConn.Exec(
		query,
		t.Tid,
		t.Msisdn,
		t.Result,
		t.OperatorCode,
		t.CountryCode,
		t.ServiceId,
		t.SubscriptionId,
		t.CampaignId,
		t.OperatorToken,
		int(t.Price),
	)
	if err != nil {
		DBErrors.Inc()
		log.WithFields(log.Fields{
			"error ": err.Error(),
			"query":  query,
			"tid":    t.Tid,
		}).Error("record transaction failed")
		return fmt.Errorf("db.Exec: %s, Query: %s", err.Error(), query)
	}
	return nil
}

func (s Record) WriteSubscriptionStatus() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":    s.Tid,
			"took":   time.Since(begin),
			"status": s.SubscriptionStatus,
		}).Debug("write subscription status")
	}()
	query := fmt.Sprintf("UPDATE %ssubscriptions SET "+
		"result = $1, "+
		"attempts_count = attempts_count + 1, "+
		"last_pay_attempt_at = $2 "+
		"where id = $3", conf.TablePrefix)

	lastPayAttemptAt := time.Now()
	_, err := dbConn.Exec(query,
		s.SubscriptionStatus,
		lastPayAttemptAt,
		s.SubscriptionId,
	)
	if err != nil {
		DBErrors.Inc()
		log.WithFields(log.Fields{
			"error ": err.Error(),
			"query":  query,
			"tid":    s.Tid,
		}).Error("notify paid subscription failed")
		return fmt.Errorf("db.Exec: %s, Query: %s", err.Error(), query)
	}
	return nil
}

func (r Record) RemoveRetry() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":  r.Tid,
			"took": time.Since(begin),
		}).Debug("remove retry")
	}()
	query := fmt.Sprintf("DELETE FROM %sretries WHERE id = $1", conf.TablePrefix)

	_, err := dbConn.Exec(query, r.RetryId)
	if err != nil {
		DBErrors.Inc()
		log.WithFields(log.Fields{
			"error ": err.Error(),
			"query":  query,
			"tid":    r.Tid,
		}).Error("delete retry failed")
		return fmt.Errorf("db.Exec: %s, query: %s", err.Error(), query)
	}
	return nil
}

func (r Record) TouchRetry() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":  r.Tid,
			"took": time.Since(begin),
		}).Debug("touch retry")
	}()
	query := fmt.Sprintf("UPDATE %sretries SET "+
		"status = '', "+
		"last_pay_attempt_at = $1, "+
		"attempts_count = attempts_count + 1 "+
		"WHERE id = $2", conf.TablePrefix)

	lastPayAttemptAt := time.Now()

	_, err := dbConn.Exec(query, lastPayAttemptAt, r.RetryId)
	if err != nil {
		DBErrors.Inc()
		log.WithFields(log.Fields{
			"error ": err.Error(),
			"query":  query,
			"retry":  fmt.Sprintf("%#v", r),
		}).Error("update retry failed")
		return fmt.Errorf("db.Exec: %s, query: %s", err.Error(), query)
	}
	return nil
}

func (r Record) StartRetry() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":  r.Tid,
			"took": time.Since(begin),
		}).Debug("add retry")
	}()
	if r.KeepDays == 0 {
		return fmt.Errorf("Retry Keep Days required, service id: %s", r.ServiceId)
	}
	if r.DelayHours == 0 {
		return fmt.Errorf("Retry Delay Hours required, service id: %s", r.ServiceId)
	}

	query := fmt.Sprintf("INSERT INTO  %sretries ("+
		"tid, "+
		"keep_days, "+
		"delay_hours, "+
		"msisdn, "+
		"operator_code, "+
		"country_code, "+
		"pixel, "+
		"publisher, "+
		"id_service, "+
		"id_subscription, "+
		"id_campaign "+
		") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		conf.TablePrefix)
	_, err := dbConn.Exec(query,
		&r.Tid,
		&r.KeepDays,
		&r.DelayHours,
		&r.Msisdn,
		&r.OperatorCode,
		&r.CountryCode,
		&r.Pixel,
		&r.Publisher,
		&r.ServiceId,
		&r.SubscriptionId,
		&r.CampaignId)
	if err != nil {
		DBErrors.Inc()
		return fmt.Errorf("db.Exec: %s, query: %s", err.Error(), query)
	}
	return nil
}

func (r Record) AddBlacklistedNumber() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":    r.Tid,
			"msisdn": r.Msisdn,
			"took":   time.Since(begin),
		}).Debug("add blacklisted")
	}()

	query := fmt.Sprintf("INSERT INTO  %smsisdn_blacklist ( msisdn ) VALUES ($1)",
		conf.TablePrefix)
	if _, err := dbConn.Exec(query, &r.Msisdn); err != nil {
		DBErrors.Inc()
		return fmt.Errorf("db.Exec: %s, query: %s", err.Error(), query)
	}
	return nil
}

func (r Record) AddPostPaidNumber() error {
	begin := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"tid":    r.Tid,
			"msisdn": r.Msisdn,
			"took":   time.Since(begin),
		}).Debug("add postpaid")
	}()

	query := fmt.Sprintf("INSERT INTO %smsisdn_postpaid ( msisdn ) VALUES ($1)", conf.TablePrefix)
	if _, err := dbConn.Exec(query, &r.Msisdn); err != nil {
		DBErrors.Inc()
		return fmt.Errorf("db.Exec: %s, query: %s", err.Error(), query)
	}
	return nil
}
