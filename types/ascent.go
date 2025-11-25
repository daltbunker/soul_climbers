package types

type AscentForm struct {
	ClimbId int
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
	AscentType string 
	AscentDate string 
	Over200Pounds bool
	Attempts string 
	Comment string 
	CreatedBy string
}
