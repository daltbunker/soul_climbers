package types

type AscentForm struct {
	NewAscent bool
	ClimbId int
	ClimbName string 
	ClimbType string
	Date string
	DateError string 
	Grade string
	GradeError string 
	RatingOptions []FormOption
	RatingError string
	AttemptOptions []FormOption
	AttemptError string
	WeightOptions []FormOption
	WeightError string
	Comment string
	CommentError string
}

type Ascent struct {
	ClimbId int
	Grade string
	Rating string 
	AscentDate string 
	Over200Pounds bool
	Attempts string 
	Comment string 
	CreatedBy string
}

type UserAscent struct {
	Ascent Ascent
	Area Area
	Climb Climb
}