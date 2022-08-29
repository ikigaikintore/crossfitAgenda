package event_handler

var Queue *gobbit

const NewWodCommand = "wod.command.new"

func init() {
	Queue = New()
}

func Topics() []string {
	return []string{
		NewWodCommand,
	}
}
