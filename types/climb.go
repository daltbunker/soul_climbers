package types

type Area struct {
	AreaId int32
	Name string
	Country string
	SubAreas string
	SubArea string
}

type Climb struct {
	ClimbId int32
	AreaId int32
	Name string
	Type string
	SubAreas string
}

type ClimbPage struct {
	Climb Climb
	Area Area
	Ascents []Ascent
}

type AreaPage struct {
	Area Area
	SubArea string
	Climbs []Climb
}