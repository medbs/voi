package core

type Message interface {
	Parse() error
	Process(*VoIPServer) error
}
