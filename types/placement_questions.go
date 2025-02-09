package types

type Answer struct {
	Text string
	Value int 
}

type QuestionInput struct {
	Id string
	Label string
	Answers []Answer
}

type QuestionInputs struct {
	SelectQuestions []QuestionInput
	CheckboxQuestions []QuestionInput
}

type Question struct {
	Id string
	Text string
	Answers []string
	AnswerPoints []string
	InputType string
	PossiblePoints int	
}