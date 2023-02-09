package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"net/http"

	"github.com/MikeCodeSun/go-mux-auth/model"
	"github.com/MikeCodeSun/go-mux-auth/util"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
  var msg string
	// get user input from req.body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprintln("read all err")
		fmt.Fprintf(w, msg)
		return
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprintln("un marshal err")
		fmt.Fprintf(w, msg)
		return
	}
	// trim input
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	// validate user input
	err = validate.Struct(user)
	if err != nil {
		fmt.Println(err.Error())	
		errMsg := util.HandlerErrorMsg(err)
		if err := json.NewEncoder(w).Encode(errMsg); err != nil {
			fmt.Println(err.Error())
			return
	  }
		return
	}
	
	// check user name & email is already reigster
	row := model.Db.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE name=? OR email=?)", user.Name, user.Email)

	var exist bool
	if err = row.Scan(&exist); err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprint("user exist")
		fmt.Fprintf(w, msg)
		return
	} else if !exist  {
		
 	// hash password
		bPassword,err:= bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			fmt.Println(err.Error())
			msg = fmt.Sprintln("hash password err")
			fmt.Fprintf(w, msg)
			return
		}
		user.Password= string(bPassword)
		// register user 
		result, err := model.Db.Exec("INSERT INTO users (name, email, password) VALUES(?,?,?)", user.Name, user.Email, user.Password)
		if err != nil {
			fmt.Println(err.Error())
			msg = fmt.Sprintln("insert database err")
			fmt.Fprintf(w, msg)
			return
		}
		fmt.Println(result)
		fmt.Fprintf(w,"register ok!")
		return
	} else {
		fmt.Println("user exist")
		fmt.Fprintln(w, "user already exist!")
		return
	}
	
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// var inputErrors = make([]util.ApiError, 2)
	var inputErrors []util.ApiError
	// get log in user input from r.body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println(err.Error())
		msg := fmt.Sprint("decode err")
		fmt.Fprintln(w, msg)
		return
	}

	// get error msg from r.body input field 
	if len(user.Name)==0 || strings.TrimSpace(user.Name ) == "" {
		nameError := util.ApiError{Param:"name", Message: "name empty"}
		// inputErrors[0] = nameError
		inputErrors = append(inputErrors, nameError)
	}
	if len(user.Password)== 0 {
		passwordError := util.ApiError{Param:"password",Message:  "password empty"}
		// inputErrors[1] = passwordError
		inputErrors = append(inputErrors, passwordError)
	}
	if (len(inputErrors) >0 ){
		json.NewEncoder(w).Encode(inputErrors)
		return
	}
	// check user is register?
	var existUser  model.User
	if err := model.Db.QueryRow("SELECT * FROM users WHERE name=?", user.Name).Scan(&existUser.ID, &existUser.Name, &existUser.Email, &existUser.Password,  &existUser.Created_at); err != nil {
		fmt.Println(err.Error())
		msg := fmt.Sprintf("User name:%s not exist", user.Name)
		fmt.Fprintln(w, msg)
		return
	} 
		// user exist
		// check user password is match
		if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password)); err != nil {
			fmt.Println(err.Error())
			fmt.Fprintln(w, "User password not right!")
			return
		}
		// password right ,then generate jwt token
		token, err := util.GenerateJwt(existUser.Name, existUser.Email, existUser.ID,)
		if err != nil {
			fmt.Fprintln(w, "Generate token err...")
			return
		}
		// deal token way 1 set cookie
		cookie := http.Cookie{
			Name: "cookie",
			Value: token,
			Path: "/",
			Domain: "localhost",
			Expires:time.Now().Local().Add(time.Hour * 24),
			Secure: false,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
		fmt.Fprintln(w, "login succussfully!")

		// deal token or way 2 return token
		// fmt.Fprintln(w, token)
}

// protected route
func ProtectedPage(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value("claim").(*util.CustomClaims)
	msg := fmt.Sprintf("Hello %v Welcome! Your email is %v!", claim.Name, claim.Email )
	fmt.Fprintln(w, msg)
}

// logot

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name: "cookie",
		Value: "",
		Path: "/",
		Domain: "localhost",
		Expires:time.Now().Local().Add(time.Hour * 0),
		Secure: false,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, "log out succussfully!")
}
