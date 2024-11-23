package prober

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"oncall-sla/metrics"
	"time"

	"github.com/spf13/viper"
)

const (
	eventStartBasis = 1741467600
	oneDayInSeconds = 60 * 60 * 24
)

type Event struct {
	Start int32  `json:"start"`
	End   int32  `json:"end"`
	User  string `json:"user"`
	Team  string `json:"team"`
	Role  string `json:"role"`
}

func Run() {
	for {
		probe()
		time.Sleep(30 * time.Second)
	}
}

func probe() {
	metrics.CreateEventRequestsTotal.Inc()

	proberCfg := viper.GetStringMapString("prober")

	delta := int32(rand.Intn(1 << 7))
	eventStart := eventStartBasis + oneDayInSeconds*delta
	eventEnd := eventStart + oneDayInSeconds

	json, _ := json.Marshal(&Event{
		Start: eventStart,
		End:   eventEnd,
		User:  proberCfg["user"],
		Team:  proberCfg["team"],
		Role:  proberCfg["role"],
	})
	body := bytes.NewReader(json)
	createRequest, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s", proberCfg["appurl"], proberCfg["eventurl"]),
		body,
	)
	createRequest.Header.Add("Content-Type", "application/json")
	createRequest.Header.Add(
		"Authorization",
		getApplicationAuth(http.MethodPost,
			proberCfg["eventurl"],
			string(json),
			proberCfg["appkey"],
			proberCfg["appname"],
		),
	)

	fmt.Printf("[INFO]: create event on %s\n", time.Unix(int64(eventStart), 0))

	beforeRequest := time.Now().UnixMilli()
	createResponse, createErr := http.DefaultClient.Do(createRequest)
	afterRequest := time.Now().UnixMilli()

	if createErr != nil || createResponse.StatusCode != 201 {
		if createErr != nil {
			fmt.Printf("[ERROR]: failed to create event; %s\n", createErr.Error())
		} else {
			fmt.Printf("[ERROR]: failed to create event; %d\n", createResponse.StatusCode)
		}
		metrics.CreateEventRequstFailTotal.Inc()
		metrics.CreateEventRequestDurationInMs.Set(float64(afterRequest - beforeRequest))
		return
	}

	defer createResponse.Body.Close()
	buf, _ := io.ReadAll(createResponse.Body)
	eventId := string(buf)

	fmt.Printf("[INFO]: created event %s\n", eventId)
	metrics.CreateEventRequestSuccessTotal.Inc()

	deleteRequest, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s%s/%s", proberCfg["appurl"], proberCfg["eventurl"], eventId),
		nil,
	)
	deleteRequest.Header.Add(
		"Authorization",
		getApplicationAuth(
			http.MethodDelete,
			fmt.Sprintf("%s/%s", proberCfg["eventurl"], eventId),
			"",
			proberCfg["appkey"],
			proberCfg["appname"],
		),
	)
	deleteResponse, deleteErr := http.DefaultClient.Do(deleteRequest)

	if deleteErr == nil && deleteResponse.StatusCode == 200 {
		fmt.Printf("[INFO]: deleted event %s\n", eventId)
	} else {
		defer deleteResponse.Body.Close()
		fmt.Printf("[ERROR]: failed to delete event %s; %s\n", eventId, deleteErr.Error())
	}

	metrics.CreateEventRequestDurationInMs.Set(float64(afterRequest - beforeRequest))
}

func getApplicationAuth(method, path, body, appKey, appName string) string {
	window := time.Now().UTC().UnixMilli() / 30000
	text := fmt.Sprintf("%d %s %s %s", window, method, path, body)
	HMAC := hmac.New(sha512.New, []byte(appKey))
	HMAC.Write([]byte(text))
	digest := base64.RawURLEncoding.WithPadding(base64.StdPadding).EncodeToString(HMAC.Sum(nil))
	return fmt.Sprintf("hmac %s:%s", appName, digest)
}
