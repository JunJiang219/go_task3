package chapter3

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID      uint
	Name    string
	Posts   []Post
	PostNum uint
}

type Post struct {
	ID           uint
	Title        string
	Content      string
	Comments     []Comment
	UserID       uint
	CommentNum   uint
	CommentState string
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&User{ID: p.UserID}).Update("post_num", gorm.Expr("post_num + 1")).Error
}

type Comment struct {
	ID      uint
	Content string
	PostID  uint
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Comment-AfterDelete: %+v\n", c)
	var count uint
	err = tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Scan(&count).Error
	if err != nil {
		return err
	}

	status := "有评论"
	if count == 0 {
		status = "无评论"
	}

	err = tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(map[string]interface{}{
		"comment_num":   count,
		"comment_state": status,
	}).Error

	return
}

// GetPostsWithMaxComments 获取评论数量最多的文章（允许多篇并列）
func GetPostsWithMaxComments(db *gorm.DB) ([]Post, error) {
	// 步骤1：构建子查询，统计每篇文章的评论数量（包括0评论）
	subQuery := db.Model(&Post{}).
		Select("posts.id, COUNT(comments.id) AS comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id")

	// 步骤2：获取最大评论数
	var maxCount int
	if err := db.Table("(?) AS t", subQuery).Select("MAX(comment_count)").Scan(&maxCount).Error; err != nil {
		return nil, err
	}

	// 步骤3：获取具有最大评论数的文章ID
	var postIDs []uint
	if err := db.Table("(?) AS t", subQuery).
		Select("id").
		Where("comment_count = ?", maxCount).
		Scan(&postIDs).Error; err != nil {
		return nil, err
	}

	// 如果没有符合条件的文章，返回空切片
	if len(postIDs) == 0 {
		return []Post{}, nil
	}

	// 步骤4：查询完整的文章信息
	var posts []Post
	if err := db.Where("id IN ?", postIDs).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})

	// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
	// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
	user1 := User{
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
	db.Create(&user1)

	// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息
	user2 := User{ID: 1}
	db.Preload("Posts.Comments").First(&user2)
	fmt.Printf("%+v\n", user2)

	// 使用Gorm查询评论数量最多的文章信息
	posts, err := GetPostsWithMaxComments(db)
	if err != nil {
		panic("查询失败: " + err.Error())
	}

	fmt.Println("评论数量最多的文章:")
	for _, post := range posts {
		fmt.Printf("ID: %d, 标题: %s\n", post.ID, post.Title)
	}

	// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段
	post := Post{
		Title:   "历史时事",
		Content: "八一建军节",
		UserID:  1,
	}
	db.Create(&post)

	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
	db.Where("post_id = ?", 1).Delete(&Comment{PostID: 1})
}
