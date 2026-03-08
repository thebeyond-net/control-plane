package input

type Message struct {
	ChatID        int
	ID            int
	InteractionID string
	Text          string
	Args          []string
}
