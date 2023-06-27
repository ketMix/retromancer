package main

import (
	"ebijam23/resources"
)

type Localizer struct {
	manager       *ResourceManager
	locale        string
	currentLocale *resources.Locale
	backupLocale  *resources.Locale
	GPT           *resources.GPT
	gptActive     bool
}

func (l *Localizer) SetGPTStyle(style string) {
	l.GPT.Style = style
}

func (l *Localizer) CheckGPTKey() bool {
	return l.GPT.CheckKey()
}

func (l *Localizer) GPTIsActive() bool {
	return l.gptActive
}

func (l *Localizer) InitGPT() {
	l.GPT = resources.InitGPT(l.manager.files)
}

func (l *Localizer) Locale() string {
	return l.locale
}

func (l *Localizer) SetLocale(loc string, gpt bool) {
	l.locale = loc
	l.backupLocale = l.manager.GetAs("locales", "en", (*resources.Locale)(nil)).(*resources.Locale)

	if !gpt {
		l.currentLocale = l.manager.GetAs("locales", loc, (*resources.Locale)(nil)).(*resources.Locale)
		return
	}
	currentLocale, err := l.GPT.GetLocale(l.backupLocale, loc)
	if err != nil {
		l.currentLocale = l.manager.GetAs("locales", loc, (*resources.Locale)(nil)).(*resources.Locale)
		l.gptActive = false
	} else {
		l.currentLocale = currentLocale
		l.gptActive = true
	}
}

func (l *Localizer) Get(key string) string {
	s := l.currentLocale.Get(key)
	if s == key {
		return l.backupLocale.Get(key)
	}
	return s
}
