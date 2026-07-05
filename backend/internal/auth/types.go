package auth

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var DemoUser = User{
	Name:  "Niklas Studio",
	Email: "niklas@example.com",
}