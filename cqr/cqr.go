package cqr

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func AddCQRHandler(allReload func(c *gin.Context), r *gin.Engine) {
	rg := r.Group("/cqr")
	rg.GET("", allReload)
}

type CQRConfig struct {
	Table      string
	ReloadFunc func() error
}

func InitCQR(cqrConfigs []CQRConfig) error {
	for _, cqrConfig := range cqrConfigs {
		begin := time.Now()
		err := cqrConfig.ReloadFunc()
		if err != nil {
			log.WithFields(log.Fields{
				"table": cqrConfig.Table,
				"error": err.Error(),
				"took":  time.Since(begin),
			}).Error("reload failed")
		} else {
			err = fmt.Errorf("%s: %s", cqrConfig.Table, err.Error())
			log.WithFields(log.Fields{
				"table": cqrConfig.Table,
				"took":  time.Since(begin),
			}).Error("reload done")
			return err
		}
	}
	return nil
}
func CQRReloadFunc(cqrConfigs []CQRConfig, c *gin.Context) func(*gin.Context) {
	var tableNames []string
	for _, v := range cqrConfigs {
		tableNames = append(tableNames, v.Table)
	}
	log.WithFields(log.Fields{
		"tables": strings.Join(tableNames, ", "),
	}).Debug("cqr request")

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
			if strings.Contains(table, cqrConfig.Table) {
				found = true
				begin := time.Now()
				err := cqrConfig.ReloadFunc()
				if err != nil {
					r.Success = false
					r.Status = http.StatusInternalServerError
					log.WithFields(log.Fields{
						"table": table,
						"error": err.Error(),
						"took":  time.Since(begin),
					}).Error("reload failed")
				} else {
					r.Success = true
					log.WithFields(log.Fields{
						"table": table,
						"took":  time.Since(begin),
					}).Info("reload done")
				}
			}
		}
		if !found {
			r.Success = false
			r.Status = http.StatusInternalServerError
			log.WithFields(log.Fields{
				"table": table,
				"error": err.Error(),
			}).Error("table not fouund")
		}
		render(r, c)
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
