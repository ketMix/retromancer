package main

import "ebijam23/resources"

type Localizer struct {
	manager       *ResourceManager
	locale        string
	currentLocale *resources.Locale
	backupLocale  *resources.Locale
}

func (l *Localizer) SetLocale(loc string) {
	l.locale = loc
	l.backupLocale = l.manager.GetAs("locales", "en", (*resources.Locale)(nil)).(*resources.Locale)
	l.currentLocale = l.manager.GetAs("locales", loc, (*resources.Locale)(nil)).(*resources.Locale)
}

func (l *Localizer) Get(key string) string {
	s := l.currentLocale.Get(key)
	if s == key {
		return l.backupLocale.Get(key)
	}
	return s
}
