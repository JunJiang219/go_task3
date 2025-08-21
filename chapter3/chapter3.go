package chapter3

import (
	"gorm.io/gorm"
)

type User struct {
	ID    uint
	Name  string
	Posts []Post
}

type Post struct {
	ID       uint
	Title    string
	Content  string
	Comments []Comment
	UserID   uint
}

type Comment struct {
	ID      uint
	Content string
	PostID  uint
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})

	user := User{
		Name: "张三",
		Posts: []Post{
			{
				Title:    "今日头条",
				Content:  "A股牛市吹响号角...",
				Comments: []Comment{{Content: "尽扯犊子"}, {Content: "假牛快跑"}},
			},
			{
				Title:    "娱乐头条",
				Content:  "zls与东家反目",
				Comments: []Comment{{Content: "吃瓜ing"}, {Content: "吃瓜ing 2"}, {Content: "吃瓜ing 3"}},
			},
		},
	}
	db.Create(&user)

	// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息
	// user := User{ID: 1}
	// db.Preload("Posts.Comments").First(&user)
	// fmt.Printf("%+v\n", user)

	// 使用Gorm查询评论数量最多的文章信息
	// var posts []Post
	// subQuery := db.Model(&Comment{}).Select("post_id, count(*) as comment_count").Group("post_id").Order("comment_count desc").Limit(1)
	// db.Debug().Joins("JOIN (?) AS top_posts ON posts.id = top_posts.post_id", subQuery).Preload("Comments").Find(&posts)
	// fmt.Printf("%+v\n", posts)
}
