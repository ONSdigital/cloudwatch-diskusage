package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/cloudwatch-diskusage/metric"
	"github.com/ONSdigital/go-ns/log"
)

var (
	filesystems   = strings.Split(os.Getenv("FILESYSTEMS"), ":")
	instanceID    = os.Getenv("INSTANCE_ID")
	namespace     = os.Getenv("NAMESPACE")
	persistent, _ = strconv.ParseBool(os.Getenv("PERSISTENT"))
	region        = os.Getenv("REGION")
)

func main() {
	m, err := metric.New("AvailableStorage", namespace, region)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	if len(filesystems) < 1 {
		log.Error(errors.New("no filesystems provided"), nil)
		os.Exit(1)
	}
	if len(instanceID) < 1 {
		log.Error(errors.New("no instance id provided"), nil)
		os.Exit(1)
	}

	log.Namespace = fmt.Sprintf("cloudwatch-diskusage-%s", instanceID)

	if !persistent {
		publishTo(m)
		return
	}

	for range time.Tick(1 * time.Minute) {
		publishTo(m)
	}
}

func publishTo(m *metric.Metric) {
	for _, fs := range filesystems {
		d := &metric.Identifier{
			Filesystem: fs,
			InstanceID: instanceID,
		}

		v, err := m.Publish(d)
		if err != nil {
			log.Error(err, nil)
			continue
		}

		log.Debug("published datapoint", log.Data{"identifier": d, "value": v})
	}
}
