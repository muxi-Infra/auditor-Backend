package request

type LoginReq struct {
	Code string `json:"code"`
}
type RegisterReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
