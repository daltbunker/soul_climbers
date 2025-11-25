package utils

type Role string

const (
	User Role = "user" 
	Admin Role = "admin" 
)

func GetAscentRatings() []string {
	return []string{"*", "**", "***", "****"}
}

func GetAscentWeights() []string {
	return []string{"over 200 pounds", "under 200 pounds"}
}

func GetAscentAttempts() []string {
	return []string{"more than 2 tries", "soft second go", "flash"}
}
