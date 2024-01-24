package models

import "github.com/golang-jwt/jwt"

type Visitor struct {
	Name         string
	LikesDueling bool
}

type CustomClaims struct {
	Username string `json:"username" form:"username" param:"username"`
	jwt.StandardClaims
}

type User struct {
	Username string `form:"username" param:"username"`
	Password string `form:"password"`
}

// type User struct {
// 	Name  string `json:"name" form:"name" query:"name"`
// 	Email string `json:"email" form:"email" query:"email"`
//   }

//   type UserDTO struct {
// 	Name    string
// 	Email   string
// 	IsAdmin bool
//   }
