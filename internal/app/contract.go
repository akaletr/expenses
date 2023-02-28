package app

type App interface {
	Init() error
	Run() error
	Stop() error
}
