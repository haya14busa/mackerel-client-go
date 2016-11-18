package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestFindMonitors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitors" {
			t.Error("request URL should be /api/v0/monitors but :", req.URL.Path)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"monitors": {
				{
					"id":            "2cSZzK3XfmG",
					"type":          "connectivity",
					"scopes":        []string{},
					"excludeScopes": []string{},
				},
				{
					"id":                              "2c5bLca8d",
					"type":                            "external",
					"name":                            "testMonitorExternal",
					"url":                             "https://www.example.com/",
					"maxCheckAttempts":                3,
					"service":                         "someService",
					"notificationInterval":            60,
					"responseTimeCritical":            5000,
					"responseTimeWarning":             10000,
					"responseTimeDuration":            5,
					"certificationExpirationCritical": 15,
					"certificationExpirationWarning":  30,
					"containsString":                  "Foo Bar Baz",
				},
				{
					"id":         "2DujfcR2kA9",
					"name":       "expression test",
					"type":       "expression",
					"expression": "avg(roleSlots('service:role','loadavg5'))",
					"operator":   ">",
					"warning":    20,
					"critical":   30,
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	monitors, err := client.FindMonitors()

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if monitors[0].Type != "connectivity" {
		t.Error("request sends json including type but: ", monitors[0])
	}

	if monitors[1].Type != "external" {
		t.Error("request sends json including type but: ", monitors[1])
	}

	if monitors[1].Service != "someService" {
		t.Error("request sends json including service but: ", monitors[1])
	}

	if monitors[1].NotificationInterval != 60 {
		t.Error("request sends json including notificationInterval but: ", monitors[1])
	}

	if monitors[1].ResponseTimeCritical != 5000 {
		t.Error("request sends json including responseTimeCritical but: ", monitors[1])
	}

	if monitors[1].ResponseTimeWarning != 10000 {
		t.Error("request sends json including responseTimeWarning but: ", monitors[1])
	}

	if monitors[1].ResponseTimeDuration != 5 {
		t.Error("request sends json including responseTimeDuration but: ", monitors[1])
	}

	if monitors[1].CertificationExpirationCritical != 15 {
		t.Error("request sends json including certificationExpirationCritical but: ", monitors[1])
	}

	if monitors[1].CertificationExpirationWarning != 30 {
		t.Error("request sends json including certificationExpirationWarning but: ", monitors[1])
	}

	if monitors[1].ContainsString != "Foo Bar Baz" {
		t.Error("request sends json including containsString but: ", monitors[1])
	}

	if monitors[2].Expression != "avg(roleSlots('service:role','loadavg5'))" {
		t.Error("request sends json including expression but: ", monitors[2])
	}
}

const monitorsjson = `
{
  "monitors": [
    {
      "id": "2cSZzK3XfmA",
      "type": "connectivity",
      "scopes": [],
      "excludeScopes": []
    },
    {
      "id"  : "2cSZzK3XfmB",
      "type": "host",
      "name": "disk.aa-00.writes.delta",
      "duration": 3,
      "metric": "disk.aa-00.writes.delta",
      "operator": ">",
      "warning": 20000.0,
      "critical": 400000.0,
      "scopes": [
        "Hatena-Blog"
      ],
      "excludeScopes": [
        "Hatena-Bookmark: db-master"
      ]
    },
    {
      "id"  : "2cSZzK3XfmC",
      "type": "service",
      "name": "Hatena-Blog - access_num.4xx_count",
      "service": "Hatena-Blog",
      "duration": 1,
      "metric": "access_num.4xx_count",
      "operator": ">",
      "warning": 50.0,
      "critical": 100.0,
      "notificationInterval": 60
    },
    {
      "id"  : "2cSZzK3XfmD",
      "type": "external",
      "name": "example.com",
      "url": "http://www.example.com",
      "service": "Hatena-Blog",
      "headers": [{"name":"Cache-Control", "value":"no-cache"}]
    },
    {
      "id"  : "2cSZzK3XfmE",
      "type": "expression",
      "name": "role average",
      "expression": "avg(roleSlots(\"server:role\",\"loadavg5\"))",
      "operator": ">",
      "warning": 5.0,
      "critical": 10.0,
      "notificationInterval": 60
    }
  ]
}
`

var wantMonitors = []monitorI{
	&MonitorConnectivity{
		ID:                   "2cSZzK3XfmA",
		Name:                 "",
		Type:                 "connectivity",
		IsMute:               false,
		NotificationInterval: 0x0000000000000000,
		Scopes:               []string{},
		ExcludeScopes:        []string{},
	},
	&MonitorHostMetric{
		ID:                   "2cSZzK3XfmB",
		Name:                 "disk.aa-00.writes.delta",
		Type:                 "host",
		IsMute:               false,
		NotificationInterval: 0x0000000000000000,
		Metric:               "disk.aa-00.writes.delta",
		Operator:             ">",
		Warning:              20000.000000,
		Critical:             400000.000000,
		Duration:             0x0000000000000003,
		Scopes: []string{
			"Hatena-Blog",
		},
		ExcludeScopes: []string{
			"Hatena-Bookmark: db-master",
		},
	},
	&MonitorServiceMetric{
		ID:                   "2cSZzK3XfmC",
		Name:                 "Hatena-Blog - access_num.4xx_count",
		Type:                 "service",
		IsMute:               false,
		NotificationInterval: 0x000000000000003c,
		Service:              "Hatena-Blog",
		Metric:               "access_num.4xx_count",
		Operator:             ">",
		Warning:              50.000000,
		Critical:             100.000000,
		Duration:             0x0000000000000001,
	},
	&MonitorExternalHTTP{
		ID:                              "2cSZzK3XfmD",
		Name:                            "example.com",
		Type:                            "external",
		IsMute:                          false,
		NotificationInterval:            0x0000000000000000,
		URL:                             "http://www.example.com",
		MaxCheckAttempts:                0.000000,
		Service:                         "Hatena-Blog",
		ResponseTimeCritical:            0.000000,
		ResponseTimeWarning:             0.000000,
		ResponseTimeDuration:            0.000000,
		ContainsString:                  "",
		CertificationExpirationCritical: 0x0000000000000000,
		CertificationExpirationWarning:  0x0000000000000000,
	},
	&MonitorExpression{
		ID:                   "2cSZzK3XfmE",
		Name:                 "role average",
		Type:                 "expression",
		IsMute:               false,
		NotificationInterval: 0x000000000000003c,
		Expression:           "avg(roleSlots(\"server:role\",\"loadavg5\"))",
		Operator:             ">",
		Warning:              5.000000,
		Critical:             10.000000,
	},
}

func TestMonitor_JSON_mapstructure(t *testing.T) {
	var data struct {
		Monitors []map[string]interface{} `json:"monitors"`
	}
	if err := json.NewDecoder(strings.NewReader(monitorsjson)).Decode(&data); err != nil {
		t.Error(err)
	}
	ms := make([]monitorI, 0, len(data.Monitors))
	for _, mmap := range data.Monitors {
		m, err := decodeMonitorFromMap(mmap)
		if err != nil {
			t.Error(err)
		}
		ms = append(ms, m)
	}
	if !reflect.DeepEqual(ms, wantMonitors) {
		t.Errorf("fail to get correct data: diff: (-got +want)\n%v", pretty.Compare(ms, wantMonitors))
	}
}
