package model

import (
	"time"
)

// Post [...]
type Post struct {
	Postid     int       `gorm:"primaryKey;column:postID;type:int;not null" json:"-"`
	Userid     int       `gorm:"index:postuser;column:userID;type:int;default:null" json:"userId"`
	User       User      `gorm:"joinForeignKey:userID;foreignKey:userID;references:Userid" json:"userList"`
	Partition  string    `gorm:"column:partition;type:varchar(10);default:null" json:"partition"`
	Title      string    `gorm:"column:title;type:varchar(20);default:null" json:"title"`
	Ptext      string    `gorm:"column:ptext;type:varchar(5000);default:null" json:"ptext"`
	CommentNum int       `gorm:"column:comment_num;type:int;default:null" json:"commentNum"`
	LikeNum    int       `gorm:"column:like_num;type:int;default:null" json:"likeNum"`
	PostTime   time.Time `gorm:"column:post_time;type:datetime;default:null" json:"postTime"`
	Heat       float64   `gorm:"column:heat;type:double;default:null" json:"heat"`
	Photos     string    `gorm:"column:photos;type:varchar(1000);default:null" json:"photos"`
}
