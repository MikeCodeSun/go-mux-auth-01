package util

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Name string `json:"name"`
	Email string `json:"email"`
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(name, email string, id int, ) (string,error){
	claims := CustomClaims{
		Name: name,
		Email: email,
		ID: id,
		RegisteredClaims:jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24*time.Hour)),
		},
	}
	secret := os.Getenv("SECRET")
	token, err:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		
		return "", err
	}
	return token, nil
}

// how to http mux auth middleware
func AuthJwt(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET")
		cookie, err:=r.Cookie("cookie")
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("No cookie Found, NOT allow to!"))
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintln(w, "jwt parse err")
		return
	}
	if !token.Valid {
		fmt.Println(err.Error())
		fmt.Fprintln(w, "Jwt Token not valid")
		return
	}
	
	claim := token.Claims.(*CustomClaims)
	ctx := context.WithValue(r.Context(), "claim", claim)
	
	
	next.ServeHTTP(w, r.WithContext(ctx))
	})
}