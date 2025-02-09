package newsfeed

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
	Id    int    `json:"id"`
}
