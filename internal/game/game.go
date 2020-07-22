package game

type GameMessage struct {
	PlayerID string
	Message  interface{}
}

//make game interface
type Game interface {
	Create()
	Join()
}
