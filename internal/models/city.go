package models

type OriginCity struct {
	CityCode string
}

type DestinationCity struct {
	CityCode string
}

type Terminal struct {
	TerminalName string
}
type CityTerminal struct {
	CityName  string
	Terminals []Terminal
}
