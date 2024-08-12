package instrumentation

import (
	"time"

	"github.com/gol-gol/golenv"

	logger "github.com/OpenChaos/ogi/logger"
)

type MetricTxn struct {
	Key       string
	StartedAt time.Time
}

var (
	ENABLE_INSTRUMENTATION = golenv.OverrideIfEnvBool("OGI_INSTRUMENTATION", false)
)

func init() {
	if !ENABLE_INSTRUMENTATION {
		return
	}
}

func StartTransaction(key string) MetricTxn {
	if !ENABLE_INSTRUMENTATION {
		return MetricTxn{}
	}
	return MetricTxn{
		Key:       key,
		StartedAt: time.Now(),
	}
}

func EndTransaction(txn *MetricTxn) {
	if !ENABLE_INSTRUMENTATION {
		return
	}
	logger.Infof("[METRIC] %s (%v)", txn.Key, time.Since(txn.StartedAt))
}
