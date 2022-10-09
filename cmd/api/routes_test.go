package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type route struct {
	path    string
	methods []string
}

var testApp application

func init() {
	observedZapCore, _ := observer.New(zap.InfoLevel)
	logger = zap.New(observedZapCore).Sugar()

	testApp = application{
		config: config{
			port: 4001,
			env:  develop,
		},
		logger:  logger,
		version: version,
	}
}

func routeExists(t *testing.T, r *mux.Router, routeToCheck string, methodToCheck []string) (found bool, err error) {
	found = false

	err = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		if tpl != routeToCheck {
			return err
		}

		mtd, err := route.GetMethods()
		if err != nil {
			return err
		}

		if !compareSlices(methodToCheck, mtd) {
			return errors.New("methods do not match")
		}

		found = true
		return nil
	})

	return
}

func Test_RouteExist(t *testing.T) {
	testRoutes := testApp.routes()
	routes := []route{
		{
			path:    "/v1/api/smart",
			methods: []string{http.MethodGet},
		},
	}

	for _, rt := range routes {
		exists, err := routeExists(t, testRoutes, rt.path, rt.methods)
		if err != nil {
			t.Errorf("test failed: %v", err)
		}
		if !exists {
			t.Errorf("did not find %s in registered routes", rt.path)
		}
	}
}

func Test_RouteDoesNotExist(t *testing.T) {
	testRoutes := testApp.routes()
	routes := []route{
		{
			path:    "/non/existing/route",
			methods: []string{http.MethodGet},
		},
	}

	for _, rt := range routes {
		exists, err := routeExists(t, testRoutes, rt.path, rt.methods)
		if err != nil {
			t.Errorf("test failed: %v", err)
		}
		if exists {
			t.Errorf("should not have found %s in registered routes", rt.path)
		}
	}
}

func Test_MethodDoesNotExist(t *testing.T) {
	testRoutes := testApp.routes()
	routes := []route{
		{
			path:    "/v1/api/smart",
			methods: []string{http.MethodGet, http.MethodPost},
		},
	}

	for _, rt := range routes {
		_, err := routeExists(t, testRoutes, rt.path, rt.methods)
		if err == nil {
			t.Error("test failed: methods should not be found")
		}
	}
}
