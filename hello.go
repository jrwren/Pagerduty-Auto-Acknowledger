package main

import (
	"context"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"os"
)

var authToken = os.Getenv("PD_AUTH_TOKEN")
var userId = os.Getenv("PD_USER_ID")
var email = os.Getenv("PD_EMAIL")

func main() {
	if authToken == "" || userId == "" || email == "" {
		fmt.Println("Missing env variables!!!")
		os.Exit(1)
	}
	client := pagerduty.NewClient(authToken)
	eps, err := client.ListIncidentsWithContext(context.TODO(), pagerduty.ListIncidentsOptions{Limit: 500, Statuses: []string{"triggered"}, UserIDs: []string{userId}})
	if err != nil {
		panic(err)
	}
	if len(eps.Incidents) == 0 {
		fmt.Println("No alerts detected.")
	} else {
		fmt.Sprintf("%d incidents detected", len(eps.Incidents))
	}
	for _, p := range eps.Incidents {
		// ack all of the incidents
		id := p.ID
		client.ManageIncidentsWithContext(context.TODO(), email, []pagerduty.ManageIncidentsOptions{{ID: id, Status: "acknowledged"}})
	}
}
