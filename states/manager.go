package states

type GPT interface {
	GetKey() string
	SetKey() string
	CheckKey() bool
	GetStyle() string
	SetStyle() string
}

type Localizer interface {
	SetGPTStyle(style string)
	GetGPTStyle() string
	SetGPTKey(key string)
	GetGPTKey() string
	CheckGPTKey() bool
	GPTIsActive() bool
	InitGPT()
	Locale() string
	SetLocale(loc string, gpt bool)
	Get(key string) string
}

type Resource interface {
	Get(category string, name string) interface{}
	GetAs(category string, name string, target interface{}) interface{}
	GetNamesWithPrefix(category string, prefix string) []string
}
