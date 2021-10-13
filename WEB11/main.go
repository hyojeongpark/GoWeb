package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 30

}

func main() {
	user := User{Name: "Park", Email: "123@gmail.com", Age: 25}
	user2 := User{Name: "Kim", Email: "222@gmail.com", Age: 40}
	users := []User{user, user2}
	// tmpl, err := template.New("Tmpl1").Parse("Name: {{.Name}}\nEmail:{{.Email}}\nAge:{{.Age}}\n")
	tmpl, err := template.New("Tmpl1").ParseFiles("templates/tmpl1.tmpl", "templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl1.tmpl", user)
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
