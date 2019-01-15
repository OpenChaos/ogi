package ogiproducer

import (
	"net/http"
	"testing"

	"bou.ke/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/OpenChaos/ogi/instrumentation"
	logger "github.com/OpenChaos/ogi/logger"
)

func TestValdiateConfig(t *testing.T) {
	var guard *monkey.PatchGuard
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()

		panic("mocked")
	})

	setTestConfig()
	assert.NotPanics(t, func() { validateConfig() })
	unsetTestConfig()
	assert.Panicsf(t, func() { validateConfig() }, "mocked")
}

func TestNewProducer(t *testing.T) {
	setTestConfig()

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*Kafka).NewProducer, func(*Kafka) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return
	})

	NewProducer()
	assert.True(t, guardB)
}

func TestProduce(t *testing.T) {
	setTestConfig()

	mp := &MockProducer{}
	var nr, nrEnd, producerGuard *monkey.PatchGuard
	var nrB, nrEndB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		nrB = true
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		nrEndB = true
		return
	})
	producerGuard = monkey.Patch(NewProducer, func() Producer {
		producerGuard.Unpatch()
		defer producerGuard.Restore()
		return mp
	})

	mp.On("Produce").Return()
	mp.On("Close").Return()
	Produce("topik", []byte{}, "key")
	assert.True(t, nrB)
	assert.True(t, nrEndB)
	mp.Mock.AssertExpectations(t)
}
