package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/slack-go/slack"
)

func main() {
	http.HandleFunc("/incident-report", func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
		if err != nil {
			fmt.Printf("Error creating verifier: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := verifier.Write(body); err != nil {
			fmt.Printf("Error writing body to verifier: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = verifier.Ensure(); err != nil {
			fmt.Printf("Error verifying signature: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		slashCommand, err := slack.SlashCommandParse(r)
		if err != nil {
			fmt.Printf("Error parsing slash command: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serviceName, description, err := parseCommand(slashCommand.Text)
		if err != nil {
			fmt.Fprintf(w, "Error parsing input: %v", err)
			return
		}

		fromEmail := os.Getenv("PAGERDUTY_EMAIL")
		if fromEmail == "" {
			fmt.Fprintf(w, "No 'from' email configured. Set PAGERDUTY_EMAIL environment variable.")
			return
		}

		// Check if the service exists in PagerDuty
		serviceID, err := getServiceID(serviceName)
		if err != nil {
			fmt.Fprintf(w, "Error checking service: %v", err)
			return
		}
		if serviceID == "" {
			fmt.Fprintf(w, "No such service found in PagerDuty")
			return
		}

		// Create the incident in PagerDuty
		if err := createPagerDutyIncident(serviceID, description, fromEmail); err != nil {
			fmt.Fprintf(w, "Error creating incident: %v", err)
			return
		}

		fmt.Fprintf(w, "Incident reported successfully")
	})

	fmt.Println("[INFO] Server listening on :3000")
	http.ListenAndServe(":3000", nil)
}

func parseCommand(input string) (serviceName, description string, err error) {
	re := regexp.MustCompile(`^"(.+?)"\s+"(.+)"$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil || len(matches) < 3 {
		return "", "", fmt.Errorf("input must be in format: /incident-report \"Service Name\" \"Description\"")
	}
	return matches[1], matches[2], nil
}

func getServiceID(serviceName string) (string, error) {
	client := pagerduty.NewClient(os.Getenv("PAGERDUTY_API_TOKEN"))
	options := pagerduty.ListServiceOptions{
		Query: serviceName,
	}

	services, err := client.ListServices(options)
	if err != nil {
		return "", err
	}
	for _, service := range services.Services {
		if service.Name == serviceName {
			return service.ID, nil
		}
	}
	return "", nil
}

func createPagerDutyIncident(serviceID, description, fromEmail string) error {
	client := pagerduty.NewClient(os.Getenv("PAGERDUTY_API_TOKEN"))

	incidentOptions := pagerduty.CreateIncidentOptions{
		Type:  "incident",
		Title: description,
		Service: &pagerduty.APIReference{
			ID:   serviceID,
			Type: "service_reference",
		},
	}

	ctx := context.Background()
	incident, err := client.CreateIncidentWithContext(ctx, fromEmail, &incidentOptions)
	if err != nil {
		return err
	}

	fmt.Println("Created incident: ", incident.ID)
	return nil
}
