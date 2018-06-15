package models

//Comment type contains post comment info
type Comment struct {
	Model

	UserName  string `form:"user_name" binding:"required"`
	Content   string `form:"content"`
	Published bool   `form:"published"`
	UserID    uint64
	User      User   `binding:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	PostID    uint64 `binding:"required" form:"post_id"`
}
