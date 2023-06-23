package resources

type Locale map[string]string

func (l *Locale) Get(key string) string {
	if v, ok := (*l)[key]; ok {
		return v
	}
	return key
}
