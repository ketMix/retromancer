package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kettek/go-multipath/v2"
	"gopkg.in/yaml.v2"
)

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `yaml:"content"`
		} `yaml:"message"`
	} `yaml:"choices"`
}

type MessageBody struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type GPTRequestBody struct {
	Model    string        `json:"model"`
	Messages []MessageBody `json:"messages"`
}

type GPT struct {
	Key          string
	SystemPrompt string // The prompt to use for the system
	MaxTokens    int    // 1-2048
	Model        string // "davinci" or "curie"
	Style        string // the style of the translation
}

// Does some set up for GPT
//  - finds the api key from "assets/key.txt"
//	- sets default values
func InitGPT(fs multipath.FS) *GPT {
	fileName := "key.txt"
	file, err := fs.ReadFile(fileName)

	k := ""
	// Check if an error was returned
	if err == nil {
		k = string(file)
	}

	return &GPT{
		Key:       k,
		Model:     "gpt-3.5-turbo",
		Style:     "an old mumbling sage",
		MaxTokens: 4096,
		SystemPrompt: `
			You are assisting with writing the story and text for a game.
			
			The game is set in a fantasy world and the player is a wizard
			who can reverse items and spells.

			- You will receive a prompt with a JSON object containing the key value pairs of the original text.
			- For each key you should create a new phrase that is different from the original.
			- A style will be requested that you should use for creating the new phrase.
			- All values should have length less than or equal to the length of their original value.
			- After creating the phrase you will translate the new phrase into the requested language.
			- There should be no escape characters in the translated phrase.
			- The value for the key should be the translated phrase only.
			- All keys must be present in the output.
			- If the value is one word, the phrase should be one word.

			Your response must be YAML unmarshalable.
		`,
	}
}

func (g *GPT) Fetch(method, url string, requestBody *[]byte) ([]byte, error) {
	if method == "POST" && requestBody == nil {
		return nil, fmt.Errorf("request body cannot be nil for POST requests")
	}
	var req *http.Request
	var err error
	if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(*requestBody))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Creates the prompt
func (gpt *GPT) createPrompt(inputLocale *Locale, locale string) string {
	str := "{"
	for k, v := range *inputLocale {
		str += fmt.Sprintf(`"%s": "%s",`, k, v)
	}
	str += "}"
	return fmt.Sprintf(`
		It should be in the similar spelling, case sensitivity, grammar, and phrasing of a %s style.
		Translate to the %s language in YAML format:
		%s
	`, gpt.Style, locale, str)
}

func (gpt *GPT) GetResponse(inputLocale *Locale, locale string) (Locale, error) {

	prompt := GPTRequestBody{
		Model: gpt.Model,
		Messages: []MessageBody{
			{
				Role:    "system",
				Content: gpt.SystemPrompt,
			},
			{
				Role:    "user",
				Content: gpt.createPrompt(inputLocale, locale),
			},
		},
	}

	promptJson, err := json.Marshal(prompt)
	if err != nil {
		return nil, err
	}

	body, err := gpt.Fetch("POST", "https://api.openai.com/v1/chat/completions", &promptJson)
	if err != nil {
		return nil, err
	}

	var s *GPTResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}
	if len(s.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}
	content := s.Choices[0].Message.Content

	var respLocale *Locale
	if err := yaml.Unmarshal([]byte(content), &respLocale); err != nil {
		return nil, err
	}
	return *respLocale, nil
}

func (g *GPT) GetLocale(baseLocale *Locale, locale string) (*Locale, error) {
	gptLocale, err := g.GetResponse(baseLocale, locale)
	if err != nil {
		return nil, err
	}
	return &gptLocale, nil
}

func (g *GPT) CheckKey() bool {
	if g.Key == "" {
		return false
	}

	_, err := g.Fetch("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return false
	}
	return true
}
