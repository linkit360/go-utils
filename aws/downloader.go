package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

type S3 interface {
	ShouldDownload(path string, reloadIfExists bool) (bool, error)
	Download(bucket, key string) (content []byte, contentLength int64, err error)
}

type s3downloader struct {
	s3   *s3.S3
	s3dl *s3manager.Downloader
	conf Config
}

type Config struct {
	Region              string        `yaml:"region"`
	Id                  string        `yaml:"access_key_id"`
	Secret              string        `yaml:"secret_access_key"`
	DownloadTimeout     time.Duration `yaml:"download_timeout"`                 // 2 minutes
	DownloadConcurrency int           `yaml:"download_concurrency" default:"1"` // 2 minutes
}

func New(s3Conf Config) S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Conf.Region),
		Credentials: credentials.NewStaticCredentials(s3Conf.Id, s3Conf.Secret, ""),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("aws load")
	}

	s3dl := &s3downloader{
		s3: s3.New(sess),
		s3dl: s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
			//d.PartSize = 2 * 1024 * 1024 // 2MB per part
			d.PartSize = 1024 * 1024 // 1MB per part
			d.Concurrency = s3Conf.DownloadConcurrency
		}),
		conf: s3Conf,
	}
	log.WithFields(log.Fields{
		"concurrrency": s3dl.s3dl.Concurrency,
	}).Info("aws init ok")
	return s3dl
}

func (s *s3downloader) ShouldDownload(path string, reloadIfExists bool) (bool, error) {
	fs, err := os.Stat(filepath.Dir(path))

	if err != nil && os.IsNotExist(err) {
		return true, nil
	}

	// if the error occured, try to clean up and try again
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("stat is nok, try to cleanup")

		if err = os.RemoveAll(path); err != nil {
			err = fmt.Errorf("os.RemoveAll: %s", err.Error())
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("campaign cleanup failed")
			return false, err
		}
		return true, err
	}

	// if we configured to reload, - redownload it
	if reloadIfExists {
		if err = os.RemoveAll(path); err != nil {
			err = fmt.Errorf("os.RemoveAll: %s", err.Error())
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("campaign remove failed")
			return false, err
		}

		log.WithFields(log.Fields{}).Debug("campaign cleaned")
		return true, nil
	}

	// if we haven't configured re-download it's ok if folder is > 0
	if fs.Size() > 0 {
		log.WithFields(log.Fields{
			"size": fs.Size(),
		}).Info("campaign already downloaded")
		return true, nil
	}

	// if folder is == 0, then try to download
	log.WithFields(log.Fields{
		"size": fs.Size(),
	}).Warn("campaign folder exists but empty")
	if err = os.RemoveAll(path); err != nil {
		err = fmt.Errorf("os.RemoveAll: %s", err.Error())
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("campaign remove failed")
		return false, err
	}
	return true, nil
}

func (s *s3downloader) Download(bucket, key string) (content []byte, contentLength int64, err error) {

	ctx := context.Background()
	var cancelFn func()
	if s.conf.DownloadTimeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, s.conf.DownloadTimeout)
	}
	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	defer cancelFn()

	//buff := &aws.WriteAtBuffer{}
	//contentLength, err = s.s3dl.DownloadWithContext(ctx, buff, &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(key),
	//})
	//
	result, err := s.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		err = fmt.Errorf("Download: %s, error: %s", key, err.Error())
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			log.WithFields(log.Fields{
				"key":     key,
				"timeout": s.conf.DownloadTimeout,
				"error":   err.Error(),
			}).Error("download canceled due to timeout")
		} else if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			log.WithFields(log.Fields{
				"key":   key,
				"error": err.Error(),
			}).Error("no such object")

		} else {
			log.WithFields(log.Fields{
				"key":   key,
				"error": err.Error(),
			}).Error("failed to download object")
		}
		return
	}

	defer result.Body.Close()
	//content = buff.Bytes()
	contentLength = *result.ContentLength

	content, err = ioutil.ReadAll(result.Body)
	if err != nil {
		err = fmt.Errorf("result.Body.Read: %s", err.Error())
		log.WithFields(log.Fields{
			"key":   key,
			"error": err.Error(),
		}).Error("failed to download object")
		return
	}

	log.WithFields(log.Fields{
		"content_len": len(content),
		"len":         contentLength,
		"timeout":     s.conf.DownloadTimeout,
	}).Info("download done")
	return content, contentLength, nil
}
