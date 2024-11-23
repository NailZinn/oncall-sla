package sla

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func fetchMetrics(query string, timestamp int32, defaultValue float32) float32 {
	endpoint := viper.GetString("prometheus.endpoint")
	request, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	queryParams := request.URL.Query()
	queryParams.Set("query", query)
	queryParams.Set("time", fmt.Sprintf("%d", timestamp))
	request.URL.RawQuery = queryParams.Encode()

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Printf("[ERROR]: failed to fetch metrics from Prometheus; %s\n", err.Error())
		return defaultValue
	}

	if response.StatusCode != 200 {
		fmt.Printf("[ERROR]: failed to fetch metrics from Prometheus; %d\n", response.StatusCode)
		return defaultValue
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	data := gjson.Get(string(body), "data.result.0.value")
	result := data.Array()

	if len(result) == 0 {
		return defaultValue
	}

	value, _ := strconv.ParseFloat(result[1].Str, 32)

	return float32(value)
}
