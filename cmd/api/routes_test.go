package main

import (
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type route struct {
	path    string
	methods []string
}

var routes = []route{
	{
		path:    "/v1/api/smart",
		methods: []string{"GET"},
	},
}

var testApp application

func init() {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	testApp = application{
		config: config{
			port: 1234,
			env:  develop,
		},
		logger:  logger,
		version: version,
		timeout: 300,
	}
}

func routeExists(t *testing.T, r *mux.Router, routeToCheck string, methodToCheck []string) {
	found := false

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			t.Error("failed to get path from route ", routeToCheck, ": ", err)
			return err
		}

		if tpl != routeToCheck {
			t.Errorf("routes don't match. Expected: %s, got: %s", routeToCheck, tpl)
			return err
		}

		mtd, err := route.GetMethods()
		if err != nil {
			t.Error("failed to get method of route ", routeToCheck, ": ", err)
			return err
		}

		if !compareSlices(methodToCheck, mtd) {
			t.Errorf("methods don't match. Expected: %s, got %s", methodToCheck, mtd)
			return err
		}

		found = true
		return nil
	})

	if err != nil {
		t.Error("test failed: ", err)
	}

	if !found {
		t.Errorf("did not find %s in registered routes", routeToCheck)
	}
}

func Test_Routes_Exist(t *testing.T) {
	testRoutes := testApp.routes()

	for _, rt := range routes {
		routeExists(t, testRoutes, rt.path, rt.methods)
	}
}
