package controller

import (
	"loginTest/api"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"time"
	"unicode/utf8"
	"github.com/gin-gonic/gin"
)

type PostMsg struct {
	UserTelephone string
	Title string
	Content string
	Partition string
}

func Post(c *gin.Context) {
	db := common.GetDB()
	var requestPostMsg PostMsg
	c.Bind(&requestPostMsg)
	// 获取参数
	userTelephone := requestPostMsg.UserTelephone
	title := requestPostMsg.Title
	content := requestPostMsg.Content
	partition :=requestPostMsg.Partition
	// 验证数据
	if !isTelephoneExist(db, userTelephone) {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}
	if len(title) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题不能为空")
		return
	}
	
	if utf8.RuneCountInString(title) > 15 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题最多为15个字")
		return
	}
	
	if len(content) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "内容不能为空")
		return
	}

	if utf8.RuneCountInString(title) > 5000 {
		response.Response(c, http.StatusBadRequest, 400, nil, "内容最多为5000个字")
		return
	}
	
	if len(partition) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "分区不能为空")
		return
	}

	if api.GetSuggestion(title) == "Block" || api.GetSuggestion(content) == "Block" {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题或内容含有不良信息,请重新编辑")
		return
	}
	
	var user model.User
    db.Where("telephone = ?", userTelephone).First(&user)

	newPost := model.Post{
		UserID: int(user.ID),
		Partition: partition,
		Title: title,
		Content: content,
		Like: 0,
		Comment: 0,
		Heat: 0,
		PostTime: time.Now(),
	}
	db.Create(&newPost)
	response.Response(c, http.StatusOK, 200, nil, "发帖成功")
}

type PostResponse struct {
	PostID uint
    UserName string
    Title string
    Content string
    Like int
    Comment int
    PostTime time.Time
	IsLiked bool
}
type BrowseMeg struct {
	UserTelephone string
	Partition string
}

func Browse(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestBrowseMeg BrowseMeg
	c.Bind(&requestBrowseMeg)
	userTelephone := requestBrowseMeg.UserTelephone
	partition := requestBrowseMeg.Partition
	var temUser model.User
    db.Where("telephone = ?", userTelephone).First(&temUser)
	var posts []model.Post
	if partition == "主页" || len(partition) == 0 {
		db.Find(&posts)
	} else {
		db.Find(&posts,"`partition` = ?",partition)
	}
	var postResponses []PostResponse
	for _, post := range posts {
		isLiked := false
		var like model.Like
    	db.Where("user_id = ? AND post_id = ?", temUser.ID, post.ID).First(&like)
    	if like.ID != 0 {
        	isLiked = true
    	}
		var user model.User
        db.Where("id = ?", post.UserID).First(&user)
		postResponse := PostResponse{
			PostID: post.ID,
            UserName: user.Name,
            Title: post.Title,
            Content: post.Content,
            Like: post.Like,
            Comment: post.Comment,
            PostTime: post.PostTime,
			IsLiked: isLiked,
        }
		postResponses = append(postResponses, postResponse)
	}
	c.JSON(http.StatusOK, postResponses)
}

type LikeMsg struct {
	UserTelephone string
	PostID uint
	IsLiked bool
}

func UpdateLike(c *gin.Context) {
	db := common.GetDB()
	var requestLikeMsg LikeMsg
	c.Bind(&requestLikeMsg)
	userTelephone := requestLikeMsg.UserTelephone
	postID := requestLikeMsg.PostID
	isLiked := requestLikeMsg.IsLiked
	// Find the user by telephone
	var user model.User
	db.Where("telephone = ?", userTelephone).First(&user)
	var post model.Post
	db.Where("id = ?", postID).First(&post)
	if isLiked {
		db.Model(&post).Update("like", post.Like - 1)
		var like model.Like
    	db.Where("user_id = ? AND post_id = ?", user.ID, post.ID).First(&like)
    	if like.ID != 0 {
        	db.Delete(&like)
    	}
	}else {
		newLike := model.Like {
			UserID: uint(user.ID),
			PostID: uint(post.ID),
		}
		if newLike.UserID != 0 && newLike.PostID != 0 {
			db.Model(&post).Update("like", post.Like + 1)
			db.Create(&newLike)
		}
	}
}