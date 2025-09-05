package models

type Command interface {
	Execute(args []string) (string, error)
	GetDescription() string
	GetUsage() string
}