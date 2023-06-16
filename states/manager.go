package states

type ResourceManager interface {
	Get(category string, name string) interface{}
	GetAs(category string, name string, target interface{}) interface{}
}
