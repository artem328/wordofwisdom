package zenquotes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider struct {
	cli *http.Client
}

func NewProvider() *Provider {
	return &Provider{cli: &http.Client{Timeout: 5 * time.Second}}
}

func (p *Provider) GetQuote(ctx context.Context) (string, error) {
	const url = "https://zenquotes.io/api/random"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("could not prepare request: %w", err)
	}

	resp, err := p.cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not fetch random quote: %w", err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("unsucessfull response code %d", resp.StatusCode)
	}

	type quote struct {
		Quote  string `json:"q"`
		Author string `json:"a"`
	}

	var quotes []quote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return "", fmt.Errorf("decode HTTP response: %w", err)
	}

	if len(quotes) == 0 || quotes[0].Quote == "" {
		return "", fmt.Errorf("no quotes were provided :(")
	}

	q := quotes[0]

	if q.Author == "" {
		return q.Quote, nil
	}

	return q.Quote + "\n\t — " + q.Author, nil
}
