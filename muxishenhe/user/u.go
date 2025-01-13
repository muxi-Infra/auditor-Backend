package user

type UserStruct struct {
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Hometown     string `json:"hometown"`
	Group        string `json:"group"`
	Timejoin     string `json:"timejoin"`
	Timeleft     string `json:"timeleft"`
	Username     string `json:"username"`
	RoleId       int    `json:"role_id"`
	Left         bool   `json:"left"`
	Info         string `json:"info"`
	AvatarUrl    string `json:"avatar_url"`
	PersonalBlog string `json:"personal_blog"`
	Github       string `json:"github"`
	Flickr       string `json:"flickr"`
	Weibo        string `json:"weibo"`
	Zhihu        string `json:"zhihu"`
}
