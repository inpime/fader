package addons

type MainConfiguration interface {
	Configuration

	Include() []string // include files
}

type Configuration interface {
	Merge(interface{}) error
}
