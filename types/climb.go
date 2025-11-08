package types

type Area struct {
	AreaId int32
	Name string
	Country string
	SubAreas string
	SubArea string
}

type Climb struct {
	AreaId int32
	Name string
	Type string
}