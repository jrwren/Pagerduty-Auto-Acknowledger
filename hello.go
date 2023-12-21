package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

var authToken = os.Getenv("PD_AUTH_TOKEN")
var userId = os.Getenv("PD_USER_ID")
var email = os.Getenv("PD_EMAIL")
var client *pagerduty.Client
var noack []string
var timer *time.Timer
var freq = 5 * time.Second

func checkAndAcknowledgeAlert() {
	defer func() {
		timer = time.AfterFunc(freq, checkAndAcknowledgeAlert)
	}()
	eps, err := client.ListIncidentsWithContext(context.TODO(),
		pagerduty.ListIncidentsOptions{
			Limit:    500,
			Statuses: []string{"triggered"},
			UserIDs:  []string{userId}},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
incs:
	for _, p := range eps.Incidents {
		// ack all of the incidents
		id := p.ID
		if len(p.Acknowledgements) > 0 {
			continue
		}
		if strings.ToLower(p.Urgency) == "low" {
			continue
		}
		// Don't ack pages.
		for i := range noack {
			if strings.Contains(p.Title, noack[i]) {
				continue incs
			}
		}
		client.ManageIncidentsWithContext(context.TODO(), email, []pagerduty.ManageIncidentsOptions{{ID: id, Status: "acknowledged"}})
		fmt.Printf("%s acknowledged: Title: %s  https://ciscospark.pagerduty.com/incidents/%s",
			time.Now().Format(time.RFC3339), p.Title, id)
		if !strings.HasSuffix(p.Title, "\n") {
			fmt.Println()
		}
		if p.Body.Details != "" {
			fmt.Printf("\t%s\n", p.Body.Details)
		}
	}
}

func main() {
	if authToken == "" || userId == "" || email == "" {
		fmt.Println("Missing env variables!!!")
		os.Exit(1)
	}

	err := readnoack()
	if err != nil {
		fmt.Println("error reading noack file: ", err)
	}

	client = pagerduty.NewClient(authToken)

	timer = time.AfterFunc(freq, checkAndAcknowledgeAlert)
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	timer.Stop()
}

func readnoack() error {
	r, err := os.Open("noack")
	if err != nil {
		return err
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		noack = append(noack, t)
	}
	fmt.Printf("noack set to: %#v\n", noack)
	return nil
}
