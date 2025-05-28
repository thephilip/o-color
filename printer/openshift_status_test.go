package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hidetatz/kubecolor/color"
	"github.com/hidetatz/kubecolor/testutil"
)

func TestOpenShiftStatusPrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		darkBackground bool // Though not used by OpenShiftStatusPrinter directly, good for consistency
		expectedOutput string
	}{
		{
			name:           "project context",
			input:          "In project my-project on server https://api.example.com:6443",
			darkBackground: true,
			// The entire line should be Cyan, as the URL rule won't override a part of an already Cyan line.
			expectedOutput: color.Apply("In project my-project on server https://api.example.com:6443", color.Cyan) + "\n",
		},
		{
			name:           "service name",
			input:          "svc/my-service - 1 pod",
			darkBackground: true,
			expectedOutput: color.Apply("svc/my-service", color.Green) + " - 1 pod\n",
		},
		{
			name:           "deployment config name",
			input:          "dc/my-app deploys istag/my-app:latest",
			darkBackground: true,
			expectedOutput: color.Apply("dc/my-app", color.Blue) + " deploys istag/my-app:latest\n",
		},
		{
			name:           "url",
			input:          "Access via https://my-app.example.com",
			darkBackground: true,
			// Note: The URL in "In project" also gets Magenta due to current regex order.
			// This specific test is for a URL not part of "In project".
			expectedOutput: "Access via " + color.Apply("https://my-app.example.com", color.Magenta) + "\n",
		},
		{
			name:           "keyword running",
			input:          "Deployment #1 (latest): running",
			darkBackground: true,
			expectedOutput: "Deployment #1 (latest): " + color.Apply("running", color.Green) + "\n",
		},
		{
			name:           "keyword deployed",
			input:          "Deployment #2: deployed 5 minutes ago",
			darkBackground: true,
			expectedOutput: "Deployment #2: " + color.Apply("deployed", color.Green) + " 5 minutes ago\n",
		},
		{
			name:           "keyword failed",
			input:          "Pod my-pod-1-build: build failed",
			darkBackground: true,
			expectedOutput: "Pod my-pod-1-build: build " + color.Apply("failed", color.Red) + "\n",
		},
		{
			name: "mixed output",
			input: `In project default on server https://localhost:8443

svc/kubernetes - 172.30.0.1:443 -> 8443
  dc/router deploys router:latest
    deployment #1 running for 2 hours
    deployment #0 failed 3 hours ago

View details with 'oc describe <resource>/<name>' or list everything with 'oc get all'.`,
			darkBackground: true,
			expectedOutput: color.Apply("In project default on server https://localhost:8443", color.Cyan) + "\n" + // Adjusted this line
				"\n" +
				color.Apply("svc/kubernetes", color.Green) + " - 172.30.0.1:443 -> 8443\n" +
				"  " + color.Apply("dc/router", color.Blue) + " deploys router:latest\n" +
				"    deployment #1 " + color.Apply("running", color.Green) + " for 2 hours\n" +
				"    deployment #0 " + color.Apply("failed", color.Red) + " 3 hours ago\n" +
				"\n" +
				"View details with 'oc describe <resource>/<name>' or list everything with 'oc get all'.\n",
		},
		{
			name:           "project context with url colorization",
			input:          "In project my-project on server https://api.example.com:6443",
			darkBackground: true,
			// The entire line should be Cyan, as the URL rule won't override a part of an already Cyan line.
			expectedOutput: color.Apply("In project my-project on server https://api.example.com:6443", color.Cyan) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printer := OpenShiftStatusPrinter{DarkBackground: tt.darkBackground}
			printer.Print(strings.NewReader(tt.input), &buf)

			testutil.MustEqual(t, tt.expectedOutput, buf.String())
		})
	}
}
