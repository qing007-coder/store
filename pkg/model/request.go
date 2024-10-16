package model

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Account          string `json:"account" binding:"min=3,max=9"`
	Password         string `json:"password" binding:"min=8,max=20"`
	Email            string `json:"email" binding:"email"`
	VerificationCode string `json:"verification_code" binding:"min=6,max=6"`
}

type SendVerificationCodeReq struct {
	Email string `json:"email"`
}

type ModifyPasswordReq struct {
	OldPassword string
	NewPassword string
	RePassword  string
}

type UploadArticleReq struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

type SearchReq struct {
	Content string `json:"content"`
}

type GetArticleDetailsReq struct {
	ID string `json:"id"`
}

type Queue struct {
	ArticleId string
}

type GetJudgeArticleReq struct {
	ArticleId string
}

type JudgeArticleReq struct {
	ArticleID string `json:"article_id"`
	Status    string `json:"status"`
}
