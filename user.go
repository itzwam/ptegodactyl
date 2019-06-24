package ptegodactyl

import (
	"strconv"
	"time"
)

// User is all infos on a user
type User struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		UUID       string      `json:"uuid"`
		Username   string      `json:"username"`
		Email      string      `json:"email"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Language   string      `json:"language"`
		RootAdmin  bool        `json:"root_admin"`
		TwoFa      bool        `json:"2fa"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"attributes"`
	client *AppClient
}

// ListUsers list all users
func (c *AppClient) ListUsers() ([]User, error) {
	users := []User{}
	err := c.list("/application/users", &users)
	if err != nil {
		return nil, err
	}
	for k := range users {
		users[k].client = c
	}
	return users, nil
}

// GetUser return a user object based on userID
func (c *AppClient) GetUser(id int) (User, error) {
	u := User{client: c}
	err := c.get("/application/users/"+strconv.Itoa(id), &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

type createUserPayload struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// CreateUser create a user and return a user object
func (c *AppClient) CreateUser(username, firstName, lastName, email, password string) (User, error) {
	u := User{client: c}
	payload := createUserPayload{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	err := c.post("/application/users/", &payload, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// EditUserPayload is a payload to send if you want to edit a user
type EditUserPayload struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Edit a user and return a user object
func (u *User) Edit(payload EditUserPayload) (User, error) {
	user := User{client: u.client}
	err := Replace(u.Attributes, &payload)
	if err != nil {
		return user, err
	}
	err = u.client.patch("/application/users/"+strconv.Itoa(u.Attributes.ID), &payload, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Delete a user
func (u *User) Delete() error {
	err := u.client.delete("/application/users/"+strconv.Itoa(u.Attributes.ID), nil)
	return err
}
