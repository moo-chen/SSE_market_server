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
	Title         string
	Content       string
	Partition     string
}

func Post(c *gin.Context) {
	db := common.GetDB()
	var requestPostMsg PostMsg
	c.Bind(&requestPostMsg)
	// 获取参数
	userTelephone := requestPostMsg.UserTelephone
	title := requestPostMsg.Title
	content := requestPostMsg.Content
	partition := requestPostMsg.Partition
	// 验证数据
	if len(userTelephone) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "返回的手机号为空")
		return
	}
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
	db.Where("phone = ?", userTelephone).First(&user)

	newPost := model.Post{
		UserID:    int(user.UserID),
		Partition: partition,
		Title:     title,
		Ptext:     content,
		Heat:      0,
		PostTime:  time.Now(),
	}
	db.Create(&newPost)
	response.Response(c, http.StatusOK, 200, nil, "发帖成功")
}

type PostResponse struct {
	PostID        uint
	UserName      string
	UserTelephone string
	Title         string
	Content       string
	Like          int
	Comment       int
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
}

type BrowseMeg struct {
	UserTelephone string
	Partition     string
	Searchinfo    string
}

func Browse(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestBrowseMeg BrowseMeg
	c.Bind(&requestBrowseMeg)
	userTelephone := requestBrowseMeg.UserTelephone
	partition := requestBrowseMeg.Partition
	searchinfo := requestBrowseMeg.Searchinfo
	var temUser model.User
	db.Where("phone = ?", userTelephone).First(&temUser)
	var posts []model.Post
	if partition == "主页" || len(partition) == 0 {
		if len(searchinfo) == 0 {
			db.Find(&posts)
		} else {
			db.Where("title LIKE ? OR ptext LIKE ?", "%"+searchinfo+"%", "%"+searchinfo+"%").Find(&posts)
		}
	} else {
		db.Find(&posts, "`partition` = ?", partition)
	}
	var postResponses []PostResponse
	for _, post := range posts {
		isSaved := false
		var save model.Psave
		db.Where("userID = ? AND ptargetID = ?", temUser.UserID, post.PostID).First(&save)
		if save.PsaveID != 0 {
			isSaved = true
		}
		isLiked := false
		var like model.Plike
		db.Where("userID = ? AND ptargetID = ?", temUser.UserID, post.PostID).First(&like)
		if like.PlikeID != 0 {
			isLiked = true
		}
		var user model.User
		db.Where("userID = ?", post.UserID).First(&user)
		postResponse := PostResponse{
			PostID:        uint(post.PostID),
			UserName:      user.Name,
			UserTelephone: user.Phone,
			Title:         post.Title,
			Content:       post.Ptext,
			Like:          post.LikeNum,
			Comment:       post.CommentNum,
			PostTime:      post.PostTime,
			IsSaved:       isSaved,
			IsLiked:       isLiked,
		}
		postResponses = append(postResponses, postResponse)
	}
	c.JSON(http.StatusOK, postResponses)
}

type SaveMsg struct {
	UserTelephone string
	PostID        uint
	IsSaved       bool
}

func UpdateSave(c *gin.Context) {
	db := common.GetDB()
	var requestSaveMsg SaveMsg
	c.Bind(&requestSaveMsg)
	userTelephone := requestSaveMsg.UserTelephone
	postID := requestSaveMsg.PostID
	isSaved := requestSaveMsg.IsSaved
	// Find the user by telephone
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	if isSaved {
		var save model.Psave
		db.Where("userID = ? AND ptargetID = ?", user.UserID, post.PostID).First(&save)
		if save.PsaveID != 0 {
			db.Delete(&save)
		}
	} else {
		newSave := model.Psave{
			UserID:    user.UserID,
			PtargetID: post.PostID,
		}
		if newSave.UserID != 0 && newSave.PtargetID != 0 {
			db.Create(&newSave)
		}
	}
}

type LikeMsg struct {
	UserTelephone string
	PostID        uint
	IsLiked       bool
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
	db.Where("phone = ?", userTelephone).First(&user)
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	if isLiked {
		db.Model(&post).Update("like_num", post.LikeNum-1)
		var like model.Plike
		db.Where("userID = ? AND ptargetID = ?", user.UserID, post.PostID).First(&like)
		if like.PlikeID != 0 {
			db.Delete(&like)
		}
	} else {
		newLike := model.Plike{
			UserID:    user.UserID,
			PtargetID: post.PostID,
		}
		if newLike.UserID != 0 && newLike.PtargetID != 0 {
			db.Model(&post).Update("like_num", post.LikeNum+1)
			db.Create(&newLike)
		}
	}
}

type IDmsg struct {
	PostID uint
}

func DeletePost(c *gin.Context) {
	db := common.GetDB()
	var ID IDmsg
	c.Bind(&ID)
	PostID := ID.PostID
	var post model.Post
	db.Where("postID = ?", PostID).First(&post)
	db.Delete(&post)
}

type Reportmsg struct {
	TargetID      uint
	UserTelephone string
	Reason        string
}

func SubmitReport(c *gin.Context) {
	db := common.GetDB()
	var reportmsg Reportmsg
	c.Bind(&reportmsg)
	TargetID := reportmsg.TargetID
	userTelephone := reportmsg.UserTelephone
	Reason := reportmsg.Reason
	if len(Reason) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "举报内容不能为空")
		return
	}
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	newSue := model.Sue{
		Targettype: "post",
		TargetID:   int(TargetID),
		UserID:     int(user.UserID),
		User:       user,
		Reason:     Reason,
		Time:       time.Now(),
		Status:     "wait",
		Finish:     false,
	}
	db.Create(&newSue)
	response.Response(c, http.StatusOK, 200, nil, "举报发送成功")
}

type PostDetailsResponse struct {
	PostID        uint
	UserName      string
	UserTelephone string
	Title         string
	Content       string
	Like          int
	Comment       int
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
}

type PostDetailsMsg struct {
	UserTelephone string
	PostID        uint
}

func ShowDetails(c *gin.Context) {
	db := common.GetDB()
	var requestPostDetailsMsg PostDetailsMsg
	c.Bind(&requestPostDetailsMsg)
	userTelephone := requestPostDetailsMsg.UserTelephone
	postID := requestPostDetailsMsg.PostID
	var temUser model.User
	db.Where("phone = ?", userTelephone).First(&temUser)
	isLiked := false
	var like model.Plike
	db.Where("userID = ? AND ptargetID = ?", temUser.UserID, postID).First(&like)
	if like.PlikeID != 0 {
		isLiked = true
	}
	isSaved := false
	var save model.Psave
	db.Where("userID = ? AND ptargetID = ?", temUser.UserID, postID).First(&save)
	if save.PsaveID != 0 {
		isSaved = true
	}
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	var user model.User
	db.Where("userID = ?", post.UserID).First(&user)
	postDetailsResponse := PostDetailsResponse{
		PostID:        uint(post.PostID),
		UserName:      user.Name,
		UserTelephone: user.Phone,
		Title:         post.Title,
		Content:       post.Ptext,
		Like:          post.LikeNum,
		Comment:       post.CommentNum,
		PostTime:      post.PostTime,
		IsSaved:       isSaved,
		IsLiked:       isLiked,
	}
	c.JSON(http.StatusOK, postDetailsResponse)
}
