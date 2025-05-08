package handler

import (
	"context"
	"fmt"
	"kanbanify-api/model"
	"kanbanify-api/utils"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/sirupsen/logrus"
)

func classifyIssue(issue model.Issue) (model.Variant, error) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("CLAUDE_API_KEY is not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	prompt := fmt.Sprintf(`
Du er en klassifiseringsmotor for et Kanban-system. Du får en tittel og en beskrivelse av en oppgave.

Du skal kun returnere én av følgende kategorier, basert på innholdet:
- "bug"
- "chore"
- "task"

Regler:
- Returner bare én av de tre ordene nøyaktig, og ingenting annet.
- Ikke forklar svaret.
- Ikke bruk punktum eller ekstra tekst.

Oppgave:
Tittel: %s
Beskrivelse: %s
`, issue.Title, issue.Description)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 50,
		Messages: []anthropic.MessageParam{{
			Role:    anthropic.MessageParamRoleUser,
			Content: utils.MakeTextContent(prompt),
		}},
		Model: anthropic.ModelClaude3_5SonnetLatest,
	})

	if err != nil {
		return "", fmt.Errorf("Failed to classify issue: %w", err)
	}

	if message == nil {
		return "", fmt.Errorf("No message returned from Claude")
	}

	res := message.Content[0].Text

	logrus.Info("Claude response: ", res)

	switch res {
	case "bug":
		return model.VariantBug, nil
	case "chore":
		return model.VariantChore, nil
	case "task":
		return model.VariantTask, nil
	default:
		return "", fmt.Errorf("Invalid response from Claude: %s", res)
	}

}
