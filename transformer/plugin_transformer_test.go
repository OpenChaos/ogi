package ogitransformer

import (
	"testing"

	logger "github.com/OpenChaos/ogi/logger"
	"github.com/gol-gol/golenv"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestNewTransformerPluginOnSuccess(t *testing.T) {
	TransformerType = "plugin"
	TransformerPluginPath = golenv.OverrideIfEnv("TRANSFORMER_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.NotPanics(t, func() { NewTransformerPlugin() })
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")
}

func TestTransformerPluginOnPluginFuncMissing(t *testing.T) {
	TransformerType = "plugin"
	TransformerPluginPath = golenv.OverrideIfEnv("TRANSFORMER_BAD_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.Panics(t, func() { NewTransformerPlugin() })
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")
}

func TestTransformerPluginOnPluginFileMissing(t *testing.T) {
	TransformerType = "plugin"
	TransformerPluginPath = ""
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.Panics(t, func() { NewTransformerPlugin() })
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "transparent")
}
