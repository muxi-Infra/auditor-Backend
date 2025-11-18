package request

type Contents struct {
	Topic       Topics  `json:"topic"`
	LastComment Comment `json:"last_comment"`
	NextComment Comment `json:"next_comment"`
}

type Topics struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}

type Comment struct {
	Content  string   `json:"content"`
	Pictures []string `json:"pictures"`
}
