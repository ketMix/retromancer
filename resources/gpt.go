package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
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
	key          string
	SystemPrompt string // The prompt to use for the system
	MaxTokens    int    // 1-2048
	Model        string // "davinci" or "curie"
	Style        string // the style of the translation
}

func findKey() string {
	return ""
}

// Does some set up for GPT
//  - finds the api key from "assets/key.txt"
//	- sets default values
func CreateGPT() (*GPT, error) {
	// The file to open
	fileName := "assets/key.txt"

	// Open the file
	file, err := os.Open(fileName)

	// Check if an error was returned
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	k := string(bytes)
	return &GPT{
		key:       k,
		Model:     "gpt-3.5-turbo",
		MaxTokens: 1000,
		Style:     "an eloquent sage",
		SystemPrompt: `
			You are assisting with writing the story and text for a game.
			
			The game is set in a fantasy world and the player is a wizard
			who can reverse items and spells.

			- You will receive a prompt with a JSON object containing the key value pairs of the original text.
			- For each key you should create a new phrase that is different of the original.
			- This new phrase should be less than or equal to the length of the original value.
			- A style be requested that you should use for creating the new phrase.
			- Escape characters in the original value should be removed.
			- After creating the phrase you will translate the new phrase into the request locale.
			- There should be no escape characters in the translated phrase.
			
			Your response must be JSON serializable.
		`,
	}, nil
}

// Creates the prompt
func (gpt *GPT) createPrompt(inputLocale Locale, locale string) string {
	str := "{"
	for k, v := range inputLocale {
		str += fmt.Sprintf(`"%s": "%s",`, k, v)
	}
	str += "}"
	return fmt.Sprintf(`
		It should be in the similar spelling, case sensitivity, and grammar
		as a person who is "%s" would write it.
		Translate the following JSON values to "%s" localization:
		"%s"
	`, gpt.Style, locale, str)
}

func (gpt *GPT) GetResponse(inputLocale Locale, locale string) (Locale, error) {
	url := "https://api.openai.com/v1/chat/completions"
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(promptJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gpt.key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	var s *GPTResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}
	if len(s.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}
	content := s.Choices[0].Message.Content

	var respLocale *Locale
	if err := json.Unmarshal([]byte(content), &respLocale); err != nil {
		return nil, err
	}
	return *respLocale, nil
}

func GetGPTLocale(baseLocale Locale, locale string) (*Locale, error) {
	gpt, err := CreateGPT()
	if err != nil {
		return nil, err
	}
	gptLocale, err := gpt.GetResponse(baseLocale, locale)
	if err != nil {
		return nil, err
	}
	return &gptLocale, nil
}
