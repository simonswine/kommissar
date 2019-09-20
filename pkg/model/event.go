package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	Object struct {
		Kind       string `json:"kind"`
		APIVersion string `json:"apiVersion"`
		Metadata   struct {
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
		} `json:"metadata"`
	}
	RawObject *json.RawMessage `json:"-"`
	Timestamp time.Time        `json:"-"`
}
