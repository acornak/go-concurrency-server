package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

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

func Test_WorkingMainFunc(t *testing.T) {
	// set env vars
	port := "4001"
	env := develop

	os.Setenv("PORT", port)
	defer os.Unsetenv("PORT")

	os.Setenv("ENV", env)
	defer os.Unsetenv("ENV")

	// set up testing logger
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	// run main func
	go main()

	time.Sleep(1 * time.Second)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal("failed to find process: ", err)
	}

	err = p.Signal(syscall.SIGINT)
	if err != nil {
		t.Fatal("failed to kill process: ", err)
	}

	// assert logs
	assert.Equal(t, fmt.Sprintf("starting server in %s mode on port %s", env, port), observedLogs.All()[0].Message)
}
