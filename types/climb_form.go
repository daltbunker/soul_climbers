package types

type FormOption struct {
	Name string
	Selected bool
}

type ClimbDraft struct {
	Name string
	RouteType string
	Country string
	AreaId int32
	Area string
	SubAreas string
}

type ClimbForm struct {
	Part int
	Name string
	NameError string
	RouteType string
	RouteTypeOptions []FormOption
	RouteTypeError string
	Country string
	CountryOptions []FormOption
	CountryError string
	Area string
	AreaError string
	SubArea string
	SubAreas []string
	SubAreasError string
}

type Result struct {
	Name string
	Id int
}

type SearchResult struct {
	Results []Result
	InputId string
}