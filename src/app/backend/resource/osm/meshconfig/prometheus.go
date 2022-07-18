package meshconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	"gopkg.in/ffmt.v1"

	osmconfigclientset "github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned"
	client "k8s.io/client-go/kubernetes"
)

type QueryInfo struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func GetPromResult(url string, result interface{}) error {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(result)
	if err != nil {
		fmt.Printf("%s", debug.Stack())
		debug.PrintStack()
		return err
	}
	return nil
}

// query metric by prom api
func QueryMetric(endpoint string, query string, method string) (*QueryInfo, error) {
	info := &QueryInfo{}
	ustr := endpoint + "/api/v1/" + method + "?query=" + query
	fmt.Println(ustr)
	u, err := url.Parse(ustr)
	if err != nil {
		return info, err
	}
	u.RawQuery = u.Query().Encode()

	err = GetPromResult(u.String(), &info)
	if err != nil {
		ffmt.Puts(info)
		return info, err
	}
	return info, nil
}

// ProxyPrometheus returns detailed information about an query
func ProxyPrometheus(osmConfigClient osmconfigclientset.Interface, client client.Interface, namespace, name, query, method string) (*QueryInfo, error) {
	log.Printf("Getting details of %s proxy Prometheus in %s namespace", name, namespace)
	// TODO
	url := "http://osm-prometheus." + namespace + ".svc:7070"
	url = "http://192.168.10.35:31001"
	promResult, err := QueryMetric(url, query, method)

	if err != nil {
		return nil, err
	}

	return promResult, nil
}
