# Slack app to report incidents in PagerDuty

## Overview

This application integrates Slack with PagerDuty, enabling users to report incidents directly from Slack to PagerDuty using a slash command. The `/incident-report` command allows users to specify a service and a description, which are then used to create an incident in PagerDuty.

## Features

- **Slash Command Integration**: Users can trigger incident reports using the `/incident-report` command directly in Slack. This command expects two inputs enclosed in quotes: the service name and the incident description.
  
- **Dynamic Input Parsing**: Handles service names and descriptions that include spaces by requiring that inputs be enclosed in quotes, e.g., `/incident-report "Service Name" "This is a detailed description of the incident."`

- **Validation of Service Names**: Before creating an incident, the application checks if the provided service name exists in PagerDuty, enhancing reliability and user feedback.

- **Secure Communication**: Uses Slack's signing secrets to verify that incoming requests are from Slack, ensuring secure communication between Slack and your server.

## Configuration

### Environment Variables

The application uses the following environment variables:

- `SLACK_SIGNING_SECRET`: Your Slack app's signing secret, used to verify that incoming requests are coming from Slack.
- `PAGERDUTY_API_TOKEN`: Your PagerDuty API token, used to authenticate API requests to PagerDuty.
- `PAGERDUTY_EMAIL`: The email associated with the PagerDuty account, used when creating incidents.

### Setting Up Slack

1. **Create a Slack App**: If you haven't already, create a Slack app in your workspace.
2. **Enable Slash Commands**: In your Slack app configuration, set up a slash command `/incident-report`.
3. **Request URL**: Configure the request URL to point to where your application is hosted followed by `/incident-report`, e.g., `https://your-server.com/incident-report`.

### Setting Up PagerDuty

1. **Create an API Key**: In PagerDuty, generate an API token with permissions to read services and create incidents.
2. **Identify the Email**: Ensure you have an email that is associated with the PagerDuty account and has necessary permissions to create incidents.

## Running the Server

1. **Install Dependencies**: Make sure Go is installed and then install required packages if not already available.
    ```bash
    go get github.com/PagerDuty/go-pagerduty
    go get github.com/slack-go/slack
    ```
2. **Set Environment Variables**: Set the `SLACK_SIGNING_SECRET`, `PAGERDUTY_API_TOKEN`, and `PAGERDUTY_EMAIL` in your environment.
3. **Run the Application**:
    ```bash
    go run main.go
    ```

## Usage

To report an incident, use the `/incident-report` command in any channel where the Slack app is installed.
