package project

type SerialGroup string

type Sequentiality struct {
	MaxInFlight int
	SerialGroup SerialGroup
}
