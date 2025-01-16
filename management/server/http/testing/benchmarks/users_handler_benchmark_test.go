//go:build benchmark
// +build benchmark

package benchmarks

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/netbirdio/netbird/management/server"
	"github.com/netbirdio/netbird/management/server/http/api"
	"github.com/netbirdio/netbird/management/server/http/testing/testing_tools"
)

// Map to store peers, groups, users, and setupKeys by name
var benchCasesUsers = map[string]testing_tools.BenchmarkCase{
	"Users - XS":     {Peers: 10000, Groups: 10000, Users: 5, SetupKeys: 10000},
	"Users - S":      {Peers: 5, Groups: 5, Users: 10, SetupKeys: 5},
	"Users - M":      {Peers: 100, Groups: 20, Users: 1000, SetupKeys: 1000},
	"Users - L":      {Peers: 5, Groups: 5, Users: 5000, SetupKeys: 5},
	"Peers - L":      {Peers: 10000, Groups: 5, Users: 5000, SetupKeys: 5},
	"Groups - L":     {Peers: 5, Groups: 10000, Users: 5000, SetupKeys: 5},
	"Setup Keys - L": {Peers: 5, Groups: 5, Users: 5000, SetupKeys: 10000},
	"Users - XL":     {Peers: 500, Groups: 50, Users: 25000, SetupKeys: 3000},
}

func BenchmarkUpdateUser(b *testing.B) {
	var expectedMetrics = map[string]testing_tools.PerformanceMetrics{
		"Users - XS":     {MinMsPerOpLocal: 100, MaxMsPerOpLocal: 160, MinMsPerOpCICD: 150, MaxMsPerOpCICD: 300},
		"Users - S":      {MinMsPerOpLocal: 0.3, MaxMsPerOpLocal: 3, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 10},
		"Users - M":      {MinMsPerOpLocal: 1, MaxMsPerOpLocal: 3, MinMsPerOpCICD: 3, MaxMsPerOpCICD: 15},
		"Users - L":      {MinMsPerOpLocal: 0.3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 2, MaxMsPerOpCICD: 20},
		"Peers - L":      {MinMsPerOpLocal: 80, MaxMsPerOpLocal: 150, MinMsPerOpCICD: 80, MaxMsPerOpCICD: 300},
		"Groups - L":     {MinMsPerOpLocal: 10, MaxMsPerOpLocal: 30, MinMsPerOpCICD: 20, MaxMsPerOpCICD: 60},
		"Setup Keys - L": {MinMsPerOpLocal: 0.3, MaxMsPerOpLocal: 3, MinMsPerOpCICD: 2, MaxMsPerOpCICD: 15},
		"Users - XL":     {MinMsPerOpLocal: 5, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 10, MaxMsPerOpCICD: 30},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	recorder := httptest.NewRecorder()

	for name, bc := range benchCasesUsers {
		b.Run(name, func(b *testing.B) {
			apiHandler, am, _ := testing_tools.BuildApiBlackBoxWithDBState(b, "../testdata/users.sql", nil, false)
			testing_tools.PopulateTestData(b, am.(*server.DefaultAccountManager), bc.Peers, bc.Groups, bc.Users, bc.SetupKeys)

			b.ResetTimer()
			start := time.Now()
			for i := 0; i < b.N; i++ {
				groupId := testing_tools.TestGroupId
				if i%2 == 0 {
					groupId = testing_tools.NewGroupId
				}
				requestBody := api.UserRequest{
					AutoGroups: []string{groupId},
					IsBlocked:  false,
					Role:       "admin",
				}

				// the time marshal will be recorded as well but for our use case that is ok
				body, err := json.Marshal(requestBody)
				assert.NoError(b, err)

				req := testing_tools.BuildRequest(b, body, http.MethodPut, "/api/users/"+testing_tools.TestUserId, testing_tools.TestAdminId)
				apiHandler.ServeHTTP(recorder, req)
			}

			testing_tools.EvaluateBenchmarkResults(b, name, time.Since(start), expectedMetrics[name], recorder)
		})
	}
}

func BenchmarkGetOneUser(b *testing.B) {
	b.Skip("Skipping benchmark as endpoint is missing")
	var expectedMetrics = map[string]testing_tools.PerformanceMetrics{
		"Users - XS":     {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Users - S":      {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Users - M":      {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Users - L":      {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Peers - L":      {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Groups - L":     {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Setup Keys - L": {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
		"Users - XL":     {MinMsPerOpLocal: 0.5, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 12},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	recorder := httptest.NewRecorder()

	for name, bc := range benchCasesUsers {
		b.Run(name, func(b *testing.B) {
			apiHandler, am, _ := testing_tools.BuildApiBlackBoxWithDBState(b, "../testdata/users.sql", nil, false)
			testing_tools.PopulateTestData(b, am.(*server.DefaultAccountManager), bc.Peers, bc.Groups, bc.Users, bc.SetupKeys)

			b.ResetTimer()
			start := time.Now()
			for i := 0; i < b.N; i++ {
				req := testing_tools.BuildRequest(b, nil, http.MethodGet, "/api/users/"+testing_tools.TestUserId, testing_tools.TestAdminId)
				apiHandler.ServeHTTP(recorder, req)
			}

			testing_tools.EvaluateBenchmarkResults(b, name, time.Since(start), expectedMetrics[name], recorder)
		})
	}
}

func BenchmarkGetAllUsers(b *testing.B) {
	var expectedMetrics = map[string]testing_tools.PerformanceMetrics{
		"Users - XS":     {MinMsPerOpLocal: 0, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 1, MaxMsPerOpCICD: 10},
		"Users - S":      {MinMsPerOpLocal: 0, MaxMsPerOpLocal: 2, MinMsPerOpCICD: 0, MaxMsPerOpCICD: 10},
		"Users - M":      {MinMsPerOpLocal: 3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 15},
		"Users - L":      {MinMsPerOpLocal: 10, MaxMsPerOpLocal: 20, MinMsPerOpCICD: 20, MaxMsPerOpCICD: 50},
		"Peers - L":      {MinMsPerOpLocal: 15, MaxMsPerOpLocal: 25, MinMsPerOpCICD: 20, MaxMsPerOpCICD: 50},
		"Groups - L":     {MinMsPerOpLocal: 15, MaxMsPerOpLocal: 25, MinMsPerOpCICD: 25, MaxMsPerOpCICD: 55},
		"Setup Keys - L": {MinMsPerOpLocal: 15, MaxMsPerOpLocal: 25, MinMsPerOpCICD: 25, MaxMsPerOpCICD: 55},
		"Users - XL":     {MinMsPerOpLocal: 80, MaxMsPerOpLocal: 120, MinMsPerOpCICD: 100, MaxMsPerOpCICD: 300},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	recorder := httptest.NewRecorder()

	for name, bc := range benchCasesUsers {
		b.Run(name, func(b *testing.B) {
			apiHandler, am, _ := testing_tools.BuildApiBlackBoxWithDBState(b, "../testdata/users.sql", nil, false)
			testing_tools.PopulateTestData(b, am.(*server.DefaultAccountManager), bc.Peers, bc.Groups, bc.Users, bc.SetupKeys)

			b.ResetTimer()
			start := time.Now()
			for i := 0; i < b.N; i++ {
				req := testing_tools.BuildRequest(b, nil, http.MethodGet, "/api/users", testing_tools.TestAdminId)
				apiHandler.ServeHTTP(recorder, req)
			}

			testing_tools.EvaluateBenchmarkResults(b, name, time.Since(start), expectedMetrics[name], recorder)
		})
	}
}

func BenchmarkDeleteUsers(b *testing.B) {
	var expectedMetrics = map[string]testing_tools.PerformanceMetrics{
		"Users - XS":     {MinMsPerOpLocal: 3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 25},
		"Users - S":      {MinMsPerOpLocal: 2, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 20},
		"Users - M":      {MinMsPerOpLocal: 2, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 20},
		"Users - L":      {MinMsPerOpLocal: 2, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 15},
		"Peers - L":      {MinMsPerOpLocal: 3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 25},
		"Groups - L":     {MinMsPerOpLocal: 2, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 20},
		"Setup Keys - L": {MinMsPerOpLocal: 3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 20},
		"Users - XL":     {MinMsPerOpLocal: 3, MaxMsPerOpLocal: 10, MinMsPerOpCICD: 5, MaxMsPerOpCICD: 20},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	recorder := httptest.NewRecorder()

	for name, bc := range benchCasesUsers {
		b.Run(name, func(b *testing.B) {
			apiHandler, am, _ := testing_tools.BuildApiBlackBoxWithDBState(b, "../testdata/users.sql", nil, false)
			testing_tools.PopulateTestData(b, am.(*server.DefaultAccountManager), bc.Peers, bc.Groups, 1000, bc.SetupKeys)

			b.ResetTimer()
			start := time.Now()
			for i := 0; i < b.N; i++ {
				req := testing_tools.BuildRequest(b, nil, http.MethodDelete, "/api/users/"+"olduser-"+strconv.Itoa(i), testing_tools.TestAdminId)
				apiHandler.ServeHTTP(recorder, req)
			}

			testing_tools.EvaluateBenchmarkResults(b, name, time.Since(start), expectedMetrics[name], recorder)
		})
	}
}
