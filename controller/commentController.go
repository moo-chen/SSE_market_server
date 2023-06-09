package controller

import (
	"github.com/gin-gonic/gin"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"time"
)

type CommentResponse struct {
	PcommentID   int
	Author       string
	AuthorAvatar string
	CommentTime  time.Time
	Content      string
	LikeNum      int
	SubComments  []Subcomment
	IsLiked      bool
}
type Commentsmsg struct {
	UserTelephone string `json:"userTelephone"`
	PostID        int    `json:"postID"`
}
type Subcomment struct {
	CcommentID     int       `json:"ccommentID"`
	Author         string    `json:"author"`
	AuthorAvatar   string    `json:"authorAvatar"`
	CommentTime    time.Time `json:"commentTime"`
	Content        string    `json:"content"`
	LikeNum        int       `json:"likeNum"`
	IsLiked        bool      `json:"isLiked"`
	UserTargetName string    `json:"userTargetName"`
}

// GetComments 给前端返回对应帖子的评论以及每条帖子评论的评论
func GetComments(c *gin.Context) {
	db := common.GetDB()
	var msg Commentsmsg
	c.Bind(&msg)
	usertelephone := msg.UserTelephone
	postid := msg.PostID
	if usertelephone == "" || postid == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "服务器无法成功解析请求")
		return
	}
	var temUser model.User
	db.Where("phone = ?", usertelephone).First(&temUser)

	var comments []CommentResponse
	var pcomments []model.Pcomment
	db.Find(&pcomments, "ptargetID = ?", postid)
	for _, pcomment := range pcomments {
		isLike := false
		var like model.Pclike
		db.Where("pctargetID = ? AND userID = ?", pcomment.PcommentID, temUser.UserID).First(&like)
		if like.PclikeID != 0 {
			isLike = true
		}
		var commentuser model.User
		db.Where("userID = ?", pcomment.UserID).First(&commentuser)
		comment := CommentResponse{
			PcommentID:   pcomment.PcommentID,
			Author:       commentuser.Name,
			AuthorAvatar: commentuser.Profile,
			CommentTime:  pcomment.Time,
			Content:      pcomment.Pctext,
			LikeNum:      pcomment.LikeNum,
			SubComments:  GetSubComments(pcomment,temUser.UserID),
			IsLiked:      isLike,
		}
		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, comments)
}

// GetSubComments 返回pcomment帖子的评论对应的子评论列表
func GetSubComments(pcomment model.Pcomment, userID int) []Subcomment {
	db := common.GetDB()
	var ccomments []model.Ccomment
	db.Find(&ccomments, "ctargetID =?", pcomment.PcommentID)
	var subcomments []Subcomment
	for _, ccomment := range ccomments {
		isLike := false
		var like model.Cclike
		db.Where("cctargetID =? AND userID =?", ccomment.CcommentID, userID).First(&like)
		if like.CclikeID != 0 {
			isLike = true
		}
		var commentuser model.User
		db.Where("userID =?", ccomment.UserID).First(&commentuser)
		comment := Subcomment{
			CcommentID:     ccomment.CcommentID,
			Author:         commentuser.Name,
			AuthorAvatar:   commentuser.Profile,
			CommentTime:    ccomment.Time,
			Content:        ccomment.Cctext,
			LikeNum:        ccomment.LikeNum,
			IsLiked:        isLike,
			UserTargetName: ccomment.UserTargetName,
		}
		subcomments = append(subcomments, comment)
	}
	if len(subcomments) == 0 {
		return []Subcomment{}
	}
	return subcomments
}

// 用于接收来自前端发表帖子的评论的结构体
type PcommentMsg struct {
	UserTelephone string
	PostID        int
	Content       string
}

// PostPcomment 进行帖子评论
func PostPcomment(c *gin.Context) {
	db := common.GetDB()
	var msg PcommentMsg
	c.Bind(&msg)
	//if api.GetSuggestion(msg.Content) == "Block" {
	//	response.Response(c, http.StatusBadRequest, 400, nil, "评论内容含有不良信息,请重新编辑")
	//	return
	//}
	var user model.User
	db.Where("phone = ?", msg.UserTelephone).First(&user)
	pcomment := model.Pcomment{
		UserID:    user.UserID,
		PtargetID: msg.PostID,
		Pctext:    msg.Content,
		Time:      time.Now(),
		LikeNum:   0,
	}
	db.Create(&pcomment)
	var post model.Post
	db.Where("postID = ?", msg.PostID).First(&post)
	db.Model(&post).Update("comment_num", post.CommentNum+1)
	comment := CommentResponse{
		PcommentID:   pcomment.PcommentID,
		Author:       user.Name,
		AuthorAvatar: user.Profile,
		CommentTime:  pcomment.Time,
		Content:      pcomment.Pctext,
		LikeNum:      pcomment.LikeNum,
		SubComments:  GetSubComments(pcomment,user.UserID),
		IsLiked:      false,
	}
	c.JSON(http.StatusOK, comment)
}

// CcommentMsg 用于接收来自前端发表评论的评论的结构体
type CcommentMsg struct {
	UserTelephone  string `json:"userTelephone"`
	PcommentID     int    `json:"pcommentID"`
	PostID        int		`json:"postID"`
	Content        string `json:"content"`
	UserTargetName string `json:"userTargetName"`
}

// PostCcomment 发表评论的评论
func PostCcomment(c *gin.Context) {
	db := common.GetDB()
	var msg CcommentMsg
	err := c.Bind(&msg)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "Bind()"+err.Error())
		return
	}
	content := msg.Content
	if len(content) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容不能为空")
		return
	}
	if len(msg.UserTelephone) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论人不能为空")
	}
	//if api.GetSuggestion(content) == "Block" {
	//	response.Response(c, http.StatusBadRequest, 400, nil, "评论内容含有不良信息,请重新编辑")
	//	return
	//}
	var user model.User
	db.Where("phone =?", msg.UserTelephone).First(&user)
	newCcomment := model.Ccomment{
		UserID:         user.UserID,
		CtargetID:      msg.PcommentID,
		Cctext:         msg.Content,
		Time:           time.Now(),
		LikeNum:        0,
		UserTargetName: msg.UserTargetName,
	}
	var post model.Post
	db.Where("postID = ?", msg.PostID).First(&post)
	db.Model(&post).Update("comment_num", post.CommentNum+1)
	db.Create(&newCcomment)
	response.Response(c, http.StatusOK, 200, nil, "评论成功！")
}

type PclikeMsg struct {
	UserTelephone string `json:"userTelephone"`
	PcommentID    uint   `json:"pcommentID"`
	IsLiked       bool   `json:"isLiked"`
}

func UpdatePcommentLike(c *gin.Context) {
	db := common.GetDB()
	var requestLikeMsg PclikeMsg
	c.Bind(&requestLikeMsg)
	userTelephone := requestLikeMsg.UserTelephone
	pcommentID := requestLikeMsg.PcommentID
	isLiked := requestLikeMsg.IsLiked
	if len(userTelephone) == 0 || pcommentID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数有误")
		return
	}
	// Find the user by telephone
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	var pcomment model.Pcomment
	db.Where("pcommentID = ?", pcommentID).First(&pcomment)
	if isLiked {
		db.Model(&pcomment).Update("like_num", pcomment.LikeNum-1)
		var like model.Pclike
		db.Where("userID = ? AND pctargetID = ?", user.UserID, pcomment.PcommentID).First(&like)
		if like.PclikeID != 0 {
			db.Delete(&like)
		}
	} else {
		newLike := model.Pclike{
			UserID:     user.UserID,
			PctargetID: pcomment.PcommentID,
		}
		if newLike.UserID != 0 && newLike.PctargetID != 0 {
			db.Model(&pcomment).Update("like_num", pcomment.LikeNum+1)
			db.Create(&newLike)
		}
	}
}

type CclikeMsg struct {
	UserTelephone string `json:"userTelephone"`
	CcommentID    uint   `json:"ccommentID"`
	IsLiked       bool   `json:"isLiked"`
}

func UpdateCcommentLike(c *gin.Context) {
	db := common.GetDB()
	var requestLikeMsg CclikeMsg
	c.Bind(&requestLikeMsg)
	userTelephone := requestLikeMsg.UserTelephone
	ccommentID := requestLikeMsg.CcommentID
	isLiked := requestLikeMsg.IsLiked
	if len(userTelephone) == 0 || ccommentID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数有误")
		return
	}
	// Find the user by ID
	var user model.User
	db.Where("phone =?", userTelephone).First(&user)
	var ccomment model.Ccomment
	db.Where("ccommentID =?", ccommentID).First(&ccomment)
	if isLiked {
		db.Model(&ccomment).Update("like_num", ccomment.LikeNum-1)
		var like model.Cclike
		db.Where("userID =? AND cctargetID =?", user.UserID, ccomment.CcommentID).First(&like)
		if like.CclikeID != 0 {
			db.Delete(&like)
		}
	} else {
		newLike := model.Cclike{
			UserID:     user.UserID,
			CctargetID: ccomment.CcommentID,
		}
		if newLike.UserID != 0 && newLike.CctargetID != 0 {
			db.Model(&ccomment).Update("like_num", ccomment.LikeNum+1)
			db.Create(&newLike)
		}
	}
}
