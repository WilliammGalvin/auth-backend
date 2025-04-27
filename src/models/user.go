package models

type User struct {
	Id            string  `json:"id"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	DisplayName   string  `json:"display_name"`
	ProfileImgx64 *string `json:"profile_img_x64"`
}

type NewUser struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}
