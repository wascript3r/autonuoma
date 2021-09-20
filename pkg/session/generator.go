package session

type Generator interface {
	GenerateID() (string, error)
}
