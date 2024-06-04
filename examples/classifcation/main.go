package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hemanta212/instructor-go/pkg/instructor"
	anthropic "github.com/liushuangls/go-anthropic/v2"
)

type LabelType string

const (
	LabelTechIssue    LabelType = "tech_issue"
	LabelBilling      LabelType = "billing"
	LabelGeneralQuery LabelType = "general_query"
)

type Label struct {
	Type LabelType `json:"type" jsonschema:"title=Label type,description=Type of the label,enum=tech_issue,enum=billing,enum=general_query"`
}

type Prediction struct {
	Labels []Label `json:"labels" jsonschema:"title=Predicted labels,description=Labels of the prediction"`
}

func classify(data string) *Prediction {
	ctx := context.Background()

	client := instructor.FromAnthropic(
		anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY")),
		instructor.WithMode(instructor.ModeToolCall),
		instructor.WithMaxRetries(3),
	)

	var prediction Prediction
	resp, err := client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Haiku20240307,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(fmt.Sprintf("Classify the following support ticket: %s", data)),
		},
		MaxTokens: 500,
	},
		&prediction,
	)
	_ = resp // sends back original response so no information loss from original API
	if err != nil {
		panic(err)
	}

	return &prediction
}

func main() {

	ticket := "My account is locked and I can't access my billing info."
	prediction := classify(ticket)

	assert(prediction.contains(LabelTechIssue), "Expected ticket to be related to tech issue")
	assert(prediction.contains(LabelBilling), "Expected ticket to be related to billing")
	assert(!prediction.contains(LabelGeneralQuery), "Expected ticket NOT to be a general query")

	fmt.Printf("%+v\n", prediction)
	/*
		&{Labels:[{Type:tech_issue} {Type:billing}]}
	*/
}

/******/

func (p *Prediction) contains(label LabelType) bool {
	for _, l := range p.Labels {
		if l.Type == label {
			return true
		}
	}
	return false
}

func assert(condition bool, message string) {
	if !condition {
		fmt.Println("Assertion failed:", message)
	}
}
