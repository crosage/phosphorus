package structs

type User struct {
	Uid      int    `json:"uid,omitempty"`
	Username string `json:"username"`
	Passhash string `json:"passhash,omitempty"`
	Type     string `json:"usertype,omitempty"`
}
