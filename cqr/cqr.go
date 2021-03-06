package cqr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func AddCQRHandler(allReload func(c *gin.Context), r *gin.Engine) {
	rg := r.Group("/cqr")
	rg.GET("", allReload)
}

type ReloadInterface interface {
	Reload() error
}

type CQRConfig struct {
	Enabled bool
	Tables  []string
	Data    ReloadInterface
	WebHook string
}

func InitCQR(cqrConfigs []CQRConfig) error {
	var tableNames []string
	for _, v := range cqrConfigs {
		tableNames = append(tableNames, v.Tables...)
	}
	log.WithFields(log.Fields{
		"tables": strings.Join(tableNames, ", "),
	}).Debug("init request")

	for _, cqrConfig := range cqrConfigs {
		if !cqrConfig.Enabled {
			log.WithFields(log.Fields{
				"tables": fmt.Sprintf("%#v", cqrConfig.Tables),
			}).Debug("disabled")
			continue
		}
		begin := time.Now()
		log.WithFields(log.Fields{
			"cqr": fmt.Sprintf("%#v", cqrConfig),
		}).Debug("cqr reload...")
		if err := cqrConfig.Data.Reload(); err != nil {
			err = fmt.Errorf("%s: %s", cqrConfig.Tables, err.Error())
			log.WithFields(log.Fields{
				"table": cqrConfig.Tables,
				"error": err.Error(),
				"took":  time.Since(begin),
			}).Error("reload failed")
			return err
		}
		if cqrConfig.WebHook != "" {
			log.WithFields(log.Fields{
				"table": cqrConfig.Tables,
				"hook":  cqrConfig.WebHook,
			}).Debug("found webhook")

			resp, err := http.Get(cqrConfig.WebHook)
			if err != nil || resp.StatusCode != 200 {
				fields := log.Fields{
					"table": cqrConfig.Tables,
					"hook":  cqrConfig.WebHook,
				}
				if resp != nil {
					fields["code"] = resp.Status
				}
				if err != nil {
					fields["error"] = err.Error()
				}
				log.WithFields(fields).Error("hook failed")
			}
		}

		log.WithFields(log.Fields{
			"table": cqrConfig.Tables,
			"took":  time.Since(begin),
		}).Info("reload done")
	}
	return nil
}

// attention: no check for cqrConfig.Enabled for certain table
// it's ok, since if it's broken, we do nothing.
func CQRReloadFunc(cqrConfigs []CQRConfig, c *gin.Context) func(*gin.Context) {
	var tableNames []string
	for _, v := range cqrConfigs {
		tableNames = append(tableNames, v.Tables...)
	}
	fn := func(c *gin.Context) {
		var err error
		r := response{Err: err, Status: http.StatusOK}

		table, exists := c.GetQuery("table")
		if !exists || table == "" {
			table, exists = c.GetQuery("t")
			if !exists || table == "" {
				err := errors.New("Table name required")
				log.WithFields(log.Fields{}).Error(err.Error())
				r.Status = http.StatusBadRequest
				r.Err = err
				render(r, c)
				return
			}
		}
		found := false
		for _, cqrConfig := range cqrConfigs {
			for _, configTableName := range cqrConfig.Tables {

				if strings.Contains(configTableName, table) {
					found = true
					begin := time.Now()
					err := cqrConfig.Data.Reload()
					if err != nil {
						r.Success = false
						r.Err = err
						r.Status = http.StatusInternalServerError
						log.WithFields(log.Fields{
							"table": table,
							"error": err.Error(),
							"took":  time.Since(begin),
						}).Error("reload failed")
						render(r, c)
						return
					} else {
						r.Success = true
						log.WithFields(log.Fields{
							"table":   table,
							"took":    time.Since(begin),
							"ua":      c.Request.UserAgent(),
							"referer": c.Request.Referer(),
						}).Info("reload done")
					}
					if cqrConfig.WebHook != "" {
						log.WithFields(log.Fields{
							"table": cqrConfig.Tables,
							"hook":  cqrConfig.WebHook,
						}).Debug("webhook")

						resp, err := http.Get(cqrConfig.WebHook)
						if err != nil || resp.StatusCode != 200 {
							fields := log.Fields{
								"table": cqrConfig.Tables,
								"hook":  cqrConfig.WebHook,
							}
							if resp != nil {
								fields["code"] = resp.Status
							}
							if err != nil {
								fields["error"] = err.Error()
							}
							log.WithFields(fields).Error("hook failed")
						}
					}
				}
			}
		}
		if !found {
			r.Success = false
			r.Status = http.StatusInternalServerError
			log.WithFields(log.Fields{
				"table":     table,
				"avialable": strings.Join(tableNames, ", "),
			}).Error("table not found")
			r.Err = fmt.Errorf("Table %s not found", table)
		}
		render(r, c)
		return
	}
	return fn
}

type response struct {
	Success bool        `json:"success,omitempty"`
	Err     error       `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  int         `json:"-"`
}

func render(msg response, c *gin.Context) {
	if msg.Err != nil {
		c.Header("Error", msg.Err.Error())
		c.Error(msg.Err)
	}
	c.JSON(msg.Status, msg)
}
