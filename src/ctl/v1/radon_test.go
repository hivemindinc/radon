/*
 * Radon
 *
 * Copyright 2018 The Radon Authors.
 * Code is licensed under the GPLv3.
 *
 */

package v1

import (
	"testing"

	"proxy"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
	"github.com/xelabs/go-mysqlstack/driver"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func TestCtlV1RadonConfig(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/config", RadonConfigHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		type radonParams1 struct {
			MaxConnections      int      `json:"max-connections"`
			MaxResultSize       int      `json:"max-result-size"`
			MaxJoinRows         int      `json:"max-join-rows"`
			DDLTimeout          int      `json:"ddl-timeout"`
			QueryTimeout        int      `json:"query-timeout"`
			TwoPCEnable         bool     `json:"twopc-enable"`
			LoadBalance         int      `json:"load-balance"`
			AllowIP             []string `json:"allowip,omitempty"`
			AuditMode           string   `json:"audit-mode"`
			StreamBufferSize    int      `json:"stream-buffer-size"`
			Blocks              int      `json:"blocks-readonly"`
			LowerCaseTableNames int      `json:"lower-case-table-names"`
		}

		// 200.
		{
			// client
			p := &radonParams1{
				MaxConnections:      1023,
				MaxResultSize:       1073741823,
				MaxJoinRows:         32767,
				QueryTimeout:        33,
				TwoPCEnable:         true,
				LoadBalance:         1,
				AllowIP:             []string{"127.0.0.1", "127.0.0.2"},
				AuditMode:           "A",
				StreamBufferSize:    16777216,
				Blocks:              128,
				LowerCaseTableNames: 1,
			}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/config", p))
			recorded.CodeIs(200)

			radonConf := proxy.Config()
			assert.Equal(t, 1023, radonConf.Proxy.MaxConnections)
			assert.Equal(t, 1073741823, radonConf.Proxy.MaxResultSize)
			assert.Equal(t, 32767, radonConf.Proxy.MaxJoinRows)
			assert.Equal(t, 0, radonConf.Proxy.DDLTimeout)
			assert.Equal(t, 33, radonConf.Proxy.QueryTimeout)
			assert.Equal(t, true, radonConf.Proxy.TwopcEnable)
			assert.Equal(t, 1, radonConf.Proxy.LoadBalance)
			assert.Equal(t, []string{"127.0.0.1", "127.0.0.2"}, radonConf.Proxy.IPS)
			assert.Equal(t, "A", radonConf.Audit.Mode)
			assert.Equal(t, 16777216, radonConf.Proxy.StreamBufferSize)
			assert.Equal(t, 128, radonConf.Router.Blocks)
			assert.Equal(t, 1, radonConf.Proxy.LowerCaseTableNames)
		}

		// Unset AllowIP.
		{
			// client
			p := &radonParams1{
				MaxConnections:   1023,
				MaxResultSize:    1073741824,
				MaxJoinRows:      32768,
				QueryTimeout:     33,
				TwoPCEnable:      true,
				AuditMode:        "A",
				StreamBufferSize: 67108864,
			}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/config", p))
			recorded.CodeIs(200)

			radonConf := proxy.Config()
			assert.Equal(t, 1023, radonConf.Proxy.MaxConnections)
			assert.Equal(t, 1073741824, radonConf.Proxy.MaxResultSize)
			assert.Equal(t, 32768, radonConf.Proxy.MaxJoinRows)
			assert.Equal(t, 0, radonConf.Proxy.DDLTimeout)
			assert.Equal(t, 33, radonConf.Proxy.QueryTimeout)
			assert.Equal(t, true, radonConf.Proxy.TwopcEnable)
			assert.Nil(t, radonConf.Proxy.IPS)
			assert.Equal(t, "A", radonConf.Audit.Mode)
			assert.Equal(t, 67108864, radonConf.Proxy.StreamBufferSize)
		}
	}
}

func TestCtlV1RadonConfigError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/config", RadonConfigHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 405.
		{
			p := &radonParams{}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/radon/config", p))
			recorded.CodeIs(405)
		}

		// 500.
		{
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/config", nil))
			recorded.CodeIs(500)
		}
	}
}

func TestCtlV1RadonReadOnly(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/readonly", ReadonlyHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 200.
		{
			// client
			p := &readonlyParams{
				ReadOnly: true,
			}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/readonly", p))
			recorded.CodeIs(200)
		}
	}
}

func TestCtlV1ReadOnlyError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/readonly", ReadonlyHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 405.
		{
			p := &readonlyParams{}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/radon/readonly", p))
			recorded.CodeIs(405)
		}

		// 500.
		{
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/readonly", nil))
			recorded.CodeIs(500)
		}
	}
}

func TestCtlV1RadonTwopc(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/twopc", TwopcHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 200.
		{
			// client
			p := &twopcParams{
				Twopc: true,
			}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/twopc", p))
			recorded.CodeIs(200)
		}
	}
}

func TestCtlV1TwopcError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/twopc", ReadonlyHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 405.
		{
			p := &twopcParams{}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/radon/twopc", p))
			recorded.CodeIs(405)
		}

		// 500.
		{
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/twopc", nil))
			recorded.CodeIs(500)
		}
	}
}

func TestCtlV1RadonThrottle(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/throttle", ThrottleHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 200.
		{
			// client
			p := &throttleParams{
				Limits: 100,
			}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/throttle", p))
			recorded.CodeIs(200)
		}
	}
}

func TestCtlV1RadonThrottleError(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		// server
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Put("/v1/radon/throttle", ThrottleHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		// 405.
		{
			p := &throttleParams{}
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "http://localhost/v1/radon/throttle", p))
			recorded.CodeIs(405)
		}

		// 500.
		{
			recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("PUT", "http://localhost/v1/radon/throttle", nil))
			recorded.CodeIs(500)
		}
	}
}

func TestCtlV1RadonStatus(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	fakedbs, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()
	address := proxy.Address()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create .*", &sqltypes.Result{})
	}

	// create database.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create database test"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	// create test table.
	{
		client, err := driver.NewConn("mock", "mock", address, "", "utf8")
		assert.Nil(t, err)
		query := "create table test.t1(id int, b int) partition by hash(id)"
		_, err = client.FetchAll(query, -1)
		assert.Nil(t, err)
	}

	{
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Get("/v1/radon/status", StatusHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/v1/radon/status", nil))
		recorded.CodeIs(200)

		want := "{\"readonly\":false}"
		got := recorded.Recorder.Body.String()
		assert.Equal(t, want, got)
	}
}

func TestCtlV1RadonApiAddress(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	_, proxy, cleanup := proxy.MockProxy(log)
	defer cleanup()

	{
		api := rest.NewApi()
		router, _ := rest.MakeRouter(
			rest.Get("/v1/radon/restapiaddress", RestAPIAddressHandler(log, proxy)),
		)
		api.SetApp(router)
		handler := api.MakeHandler()

		recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/v1/radon/restapiaddress", nil))
		recorded.CodeIs(200)
	}
}
