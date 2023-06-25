package main

import (
	"ebijam23/resources"
	"fmt"
)

type Localizer struct {
	manager       *ResourceManager
	locale        string
	currentLocale *resources.Locale
	backupLocale  *resources.Locale
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
	fmt.Println("Fetching from GPT")
	currentLocale, err := resources.GetGPTLocale(l.manager.files, l.backupLocale, loc)
	if err != nil {
		fmt.Println("Failed to get GPT locale:", err)
		l.currentLocale = l.manager.GetAs("locales", loc, (*resources.Locale)(nil)).(*resources.Locale)
	} else {
		l.currentLocale = currentLocale
	}
}

func (l *Localizer) Get(key string) string {
	s := l.currentLocale.Get(key)
	if s == key {
		return l.backupLocale.Get(key)
	}
	return s
}
