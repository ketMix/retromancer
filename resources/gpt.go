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

const url = "https://api.openai.com/v1/chat/completions"

type GPT struct {
	key          string
	SystemPrompt string // The prompt to use for the system
	MaxTokens    int    // 1-2048
	Model        string // "davinci" or "curie"
	Style        string // the style of the translation
}

// Does some set up for GPT
//  - finds the api key from "assets/key.txt"
//	- sets default values
func CreateGPT(fs multipath.FS) (*GPT, error) {
	fileName := "key.txt"
	file, err := fs.ReadFile(fileName)

	// Check if an error was returned
	if err != nil {
		return nil, err
	}

	k := string(file)
	// Check if the key is empty
	if k == "" {
		return nil, fmt.Errorf("key is empty")
	}
	return &GPT{
		key:       k,
		Model:     "gpt-3.5-turbo",
		MaxTokens: 2048,
		Style:     "an old mumbling sage",
		SystemPrompt: `
			You are assisting with writing the story and text for a game.
			
			The game is set in a fantasy world and the player is a wizard
			who can reverse items and spells.

			- You will receive a prompt with a JSON object containing the key value pairs of the original text.
			- For each key you should create a new phrase that is different of the original.
			- A style will be requested that you should use for creating the new phrase.
			- All values should be less than or equal to the length of their original value.
			- Escape characters in the original value should be removed.
			- After creating the phrase you will translate the new phrase into the request locale.
			- There should be no escape characters in the translated phrase.
			- The value for the key should be the translated phrase only.
			- All keys must be present in the output.
			Your response must be YAML unmarshalable.
		`,
	}, nil
}

func (g *GPT) Fetch(requestBody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.key)

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
		Translate to "%s" in YAML format:
		"%s"
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

	body, err := gpt.Fetch(promptJson)
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

func GetGPTLocale(fs multipath.FS, baseLocale *Locale, locale string) (*Locale, error) {
	if baseLocale == nil {
		return nil, fmt.Errorf("baseLocale is nil")
	}
	gpt, err := CreateGPT(fs)
	if err != nil {
		return nil, err
	}
	gptLocale, err := gpt.GetResponse(baseLocale, locale)
	if err != nil {
		return nil, err
	}
	return &gptLocale, nil
}
