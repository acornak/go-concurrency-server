package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func Test_InvalidPort(t *testing.T) {
	// set env vars
	os.Setenv("PORT", "not-valid")
	defer os.Unsetenv("PORT")

	// set up testing logger
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	// run main func
	main()

	// assert logs
	require.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "failed to get port from env vars: {error 26 0  strconv.Atoi: parsing \"not-valid\": invalid syntax}", observedLogs.All()[0].Message)
}
