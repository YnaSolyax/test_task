package storage

type Todo struct {
	Id          int
	Title       string
	Description OptinonalString
	Status      OptinonalInt
}

type OptinonalInt struct {
	Value int
	IsSet bool
}

type OptinonalString struct {
	Value string
	IsSet bool
}
