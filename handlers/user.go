package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"unicode"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"golang.org/x/crypto/bcrypt"
)

func HandleNewUser(w http.ResponseWriter, r *http.Request) {

	user := types.User{}
	user.Email = sanitize(r.FormValue("email"))
	user.Password = sanitize(r.FormValue("password"))
	user.Username = sanitize(r.FormValue("username"))

	signupForm := types.SignupForm{}
	signupForm.Email = user.Email
	signupForm.Username= user.Username
	signupForm.Password = user.Password

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		signupForm.EmailError = "please enter a valid email address"
	}

	_, err = db.GetUserByEmail(r, user.Email)
	if err == nil {
		signupForm.EmailError= "this email address is already taken"
	}

	_, err = db.GetUserByUsername(r, user.Username)
	if err == nil {
		signupForm.UsernameError = "this username is already taken"
	}

	if len(user.Username) == 0 {
		signupForm.UsernameError = "please enter a username"
	}

	signupForm.PasswordError = validatePassword(user.Password)

	if signupForm.EmailError != "" || signupForm.UsernameError != "" || signupForm.PasswordError != "" {
		renderComponent(w, "signup", "signup", signupForm)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = string(hash)

	dbUser, err := db.NewUser(r, user)
	if err != nil {
		HandleClientError(w, err)
		return
	}

	err = NewSession(r, w, dbUser)
	if err != nil {
		HandleClientError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}

func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	loginForm := types.LoginForm{}
	user := types.User{}
	user.Email = sanitize(r.FormValue("email"))
	user.Password = sanitize(r.FormValue("password"))

	dbUser, err := db.GetUserByEmail(r, user.Email)
	if err != nil {
		log.Printf("failed getting user from DB: %v", err)
		loginForm.Email = user.Email
		loginForm.EmailError = "email not found"
		renderComponent(w, "login", "login", loginForm)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Printf("failed checking password: %v", err)
		loginForm.Email = user.Email
		loginForm.PasswordError = "password is incorrect"
		renderComponent(w, "login", "login", loginForm)
		return
	}

	err = NewSession(r, w, dbUser)
	if err != nil {
		log.Printf("issue creating session: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}

func HandleEmailResetLink(w http.ResponseWriter, r *http.Request) {
	email := sanitize(r.FormValue("email"))

	_, err := db.GetUserByEmail(r, email)
	if err != nil {
		fmt.Fprint(w, "<span class=\"error\">Account not found</span>")
		return
	}

	// sendResetLink(email)

	resetToken, err := db.NewResetToken(r, email)	
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	fmt.Println(resetToken)

	fmt.Fprintf(w, "Email successfuly sent to: %s", email)
	// TODO: this should be the link sent to the email
	// resetUrl := fmt.Sprintf("/login/reset/password?token=%s",  resetToken.Token)
}

func HandleSignOut(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

func HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	resetToken, err := r.Cookie("Reset-Token")
	if err != nil {
		log.Printf("Cookie 'Reset-Token' not found %v", err)
		HandleClientError(w, err)
		return
	}

	// 1 - validate token
	// 2 - update password

	fmt.Fprintf(w, "Token Found: %s", resetToken.Value)
}

func sendResetLink(email string) {
	// Sign up for email service with namecheap (soulclimbers.org)
	c, err := smtp.Dial("smtp.gmail.com:25")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Mail("some_email_address"); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(email); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func validatePassword(s string) string {
	required := []string{"minimum eight characters", "a letter", "a number", "a special character"}
	if len(s) >= 8 {
		required[0] = ""
	}
	for _, char := range s {
		switch {
		case unicode.IsLetter(char):
			required[1] = ""
		case unicode.IsNumber(char):
			required[2] = ""
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			required[3] = ""
		}
	}

	errorMessage := ""
	for i := 0; i < len(required); i++ {
		if len(errorMessage) == 0 && len(required[i]) > 0 {
			errorMessage = required[i]
		} else if len(required[i]) > 0 {
			errorMessage += ", " + required[i]
		}
	}

	return errorMessage 
}
