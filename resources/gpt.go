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

var gptKey string

type GPT struct {
	Key          string
	IsInitKey    bool   // If true, the key was from initialization
	SystemPrompt string // The prompt to use for the system
	MaxTokens    int    // 1-8192
	Model        string // "gpt-3.5-turbo-16k", "gpt-3.5-turbo" or "gpt-4"
	Style        string // the style of the translation
}

// Does some set up for GPT
//   - uses the gptKey variable that is set during release builds
//   - or finds the api key from "assets/key.txt"
//   - sets default values
func InitGPT(fs multipath.FS) *GPT {
	var k string
	k = ""
	if gptKey != "" {
		k = gptKey
	} else {
		fileName := "key.txt"
		file, err := fs.ReadFile(fileName)

		// Check if an error was returned
		if err == nil {
			k = string(file)
		}
	}

	return &GPT{
		Key:       k,
		Model:     "gpt-3.5-turbo-16k",
		Style:     "an old mumbling sage",
		MaxTokens: 8192,
		SystemPrompt: `
			You are an assistant who has been asked to translate the values of a YAML object into a different style.
			You must emulate the style or personality of the user-requested style.

			You must abide by the following rules:
			- Do not modify the keys
			- You will receive a query with a YAML object containing the key value pairs of the original text.
			- For each key you should create a new value that must be different from the original.
			- A style will be requested that you should use for creating the new value.
			- You will adhere your responses to the style and apply the given writing style or personality to it.
			- There should be no escape characters in the new values.
			- All key value pairs in the input must be present in the output.
			- If the value is one word, the phrase should be one word.
			- Your response must be YAML unmarshalable.
			- Only respond with the new YAML object.
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
	if resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error: %s", string(body))
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Creates the prompt
func (gpt *GPT) createPrompt(inputLocale *Locale, locale string) string {
	prompt := fmt.Sprintf(`You must strictly apply the style/personality of "%s" to the values.`, gpt.Style)
	if locale != "en" {
		prompt += fmt.Sprintf(`You must translate the values into the "%s" language using the appropriate character set.`, locale)
	}
	prompt += "The YAML object to be modified. All of the keys in this object must be present in the output:\n"
	for k, v := range *inputLocale {
		prompt += fmt.Sprintf(`"%s": "%s"\n`, k, v)
	}
	return prompt
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

func (g *GPT) GetKey() string {
	if g.IsInitKey {
		return ""
	}
	return g.Key
}

func (g *GPT) SetKey(key string) {
	g.Key = key
	g.IsInitKey = false
}
