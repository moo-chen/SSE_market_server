package controller

import (
	"fmt"
	"loginTest/api"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CommentResponse struct {
	PcommentID      int
	Author          string
	AuthorTelephone string
	AuthorAvatar    string
	CommentTime     time.Time
	Content         string
	LikeNum         int
	SubComments     []Subcomment
	IsLiked         bool
}
type Commentsmsg struct {
	UserTelephone string `json:"userTelephone"`
	PostID        int    `json:"postID"`
}
type Subcomment struct {
	CcommentID      int       `json:"ccommentID"`
	Author          string    `json:"author"`
	AuthorTelephone string    `json:"authorTelephone"`
	AuthorAvatar    string    `json:"authorAvatar"`
	CommentTime     time.Time `json:"commentTime"`
	Content         string    `json:"content"`
	LikeNum         int       `json:"likeNum"`
	IsLiked         bool      `json:"isLiked"`
	UserTargetName  string    `json:"userTargetName"`
	ShowMenu        bool      `json:"showMenu"`
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
			PcommentID:      pcomment.PcommentID,
			Author:          commentuser.Name,
			AuthorTelephone: commentuser.Phone,
			AuthorAvatar:    commentuser.AvatarURL,
			CommentTime:     pcomment.Time,
			Content:         pcomment.Pctext,
			LikeNum:         pcomment.LikeNum,
			SubComments:     GetSubComments(pcomment, temUser.UserID),
			IsLiked:         isLike,
		}
		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, comments)
}

type IDmesg struct {
	PcommentID uint
}

func DeletePcomment(c *gin.Context) {
	db := common.GetDB()
	var ID IDmesg
	c.Bind(&ID)
	PcommentID := ID.PcommentID
	var pcomment model.Pcomment
	db.Where("pcommentID = ?", PcommentID).First(&pcomment)
	// 剪掉相应的热度
	currentTime := time.Now()
	timedif := currentTime.Sub(pcomment.Time)
	hours := timedif.Hours()
	days := int(hours / 24)
	fmt.Println("days: ", days)
	weightComment := float64(6)
	var post model.Post
	db.Where("postID = ?", pcomment.PtargetID).First(&post)
	// 帖子评论数减相应数字
	var ccomment model.Ccomment
	var count int64
	db.Model(&ccomment).Where("ctargetID = ?", pcomment.PcommentID).Count(&count)
	db.Model(&post).UpdateColumn("comment_num", gorm.Expr("comment_num - ?", count+1))
	if days > 0 {
		weightCommentPower := math.Pow(0.5, float64(days))
		deleteHeat := math.Pow(weightComment, weightCommentPower)
		db.Model(&post).Update("heat", post.Heat - (deleteHeat + float64(count)))
	} else {
		deleteCcommentHeat := float64(count * int64(weightComment))
		db.Model(&post).Update("heat", post.Heat - (weightComment + deleteCcommentHeat))
	}
	//
	db.Delete(&pcomment)
}

type IDmesag struct {
	CcommentID uint
}

func DeleteCcomment(c *gin.Context) {
	db := common.GetDB()
	var ID IDmesag
	c.Bind(&ID)
	CcommentID := ID.CcommentID
	var ccomment model.Ccomment
	db.Where("ccommentID = ?", CcommentID).First(&ccomment)
	// 剪掉相应的热度
	currentTime := time.Now()
	timedif := currentTime.Sub(ccomment.Time)
	hours := timedif.Hours()
	days := int(hours / 24)
	fmt.Println("days: ", days)
	weightComment := float64(6)
	var targetcommentid model.Pcomment
	db.Where("pcommentID= ?", ccomment.CtargetID).First(&targetcommentid)
	var post model.Post
	db.Where("postID = ?", targetcommentid.PtargetID).First(&post)
	// 帖子评论数减一
	db.Model(&post).UpdateColumn("comment_num", gorm.Expr("comment_num - ?", 1))
	if days > 0 {
		weightCommentPower := math.Pow(0.5, float64(days))
		deleteHeat := math.Pow(weightComment, weightCommentPower)
		db.Model(&post).Update("heat", post.Heat-deleteHeat)
	} else {
		db.Model(&post).Update("heat", post.Heat-weightComment)
	}
	//
	db.Delete(&ccomment)
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
			CcommentID:      ccomment.CcommentID,
			Author:          commentuser.Name,
			AuthorTelephone: commentuser.Phone,
			AuthorAvatar:    commentuser.AvatarURL,
			CommentTime:     ccomment.Time,
			Content:         ccomment.Cctext,
			LikeNum:         ccomment.LikeNum,
			IsLiked:         isLike,
			UserTargetName:  ccomment.UserTargetName,
			ShowMenu:        false,
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
	if len(msg.Content) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容不能为空")
		return
	}
	if len(msg.UserTelephone) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论人不能为空")
		return
	}
	if api.GetSuggestion(msg.Content) == "Block" {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容含有不良信息,请重新编辑")
		return
	}
	var user model.User
	var tempost model.Post
	db.Where("phone = ?", msg.UserTelephone).First(&user)
	db.Where("postID =?", msg.PostID).First(&tempost)
	pcomment := model.Pcomment{
		UserID:    user.UserID,
		PtargetID: msg.PostID,
		Pctext:    msg.Content,
		Time:      time.Now(),
		LikeNum:   0,
	}
	// 创建一条帖子评论
	db.Create(&pcomment)
	// 如果用户自己评论自己的帖子，则不用通知
	if tempost.UserID != user.UserID {
		notice := model.Notice{
			Receiver: tempost.UserID,
			User:     model.User{},
			Sender:   user.UserID,
			Type:     "pcomment",
			Ntext:    msg.Content,
			Time:     time.Now(),
			Read:     false,
			Target:   pcomment.PcommentID,
		}
		// 创建一条通知
		db.Create(&notice)
	}

	var post model.Post
	db.Where("postID = ?", msg.PostID).First(&post)
	db.Model(&post).Update("comment_num", post.CommentNum+1)
	// 在这里设置 评论 的权重
	weightComment := float64(6)
	db.Model(&post).Update("heat", post.Heat+weightComment)
	comment := CommentResponse{
		PcommentID:   pcomment.PcommentID,
		Author:       user.Name,
		AuthorAvatar: user.AvatarURL,
		CommentTime:  pcomment.Time,
		Content:      pcomment.Pctext,
		LikeNum:      pcomment.LikeNum,
		SubComments:  GetSubComments(pcomment, user.UserID),
		IsLiked:      false,
	}
	c.JSON(http.StatusOK, comment)
}

// CcommentMsg 用于接收来自前端发表评论的评论的结构体
type CcommentMsg struct {
	UserTelephone  string `json:"userTelephone"`
	PcommentID     int    `json:"pcommentID"`
	PostID         int    `json:"postID"`
	Content        string `json:"content"`
	UserTargetName string `json:"userTargetName"`
	CcommentID     int    `json:"ccommentID"`
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
		return
	}
	if api.GetSuggestion(content) == "Block" {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容含有不良信息,请重新编辑")
		return
	}
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
	// 数据库创建一条新的评论的评论
	db.Create(&newCcomment)
	var tempcomment model.Pcomment
	db.Where("pcommentID =?", msg.PcommentID).First(&tempcomment)
	// 如果是评论的评论
	// 如果是用户在自己发的一级评论下发回复，那么不需要通知
	if tempcomment.UserID != user.UserID {
		notice1 := model.Notice{
			Receiver: tempcomment.UserID,
			User:     model.User{},
			Sender:   user.UserID,
			Type:     "ccomment",
			Ntext:    msg.Content,
			Time:     time.Now(),
			Read:     false,
			Target:   newCcomment.CcommentID,
		}
		// 数据库创建一条通知
		db.Create(&notice1)
	}
	// 如果是二级评论的回复
	if msg.UserTargetName != "" {
		var temccomment model.Ccomment
		db.Where("ccommentID =?", msg.CcommentID).First(&temccomment)
		// 如果是自己回复自己就不用发通知,还有一种情况，就是上面的一级回复已经发了通知，这里就不需要重发了
		if temccomment.UserID != user.UserID && tempcomment.UserID != temccomment.UserID {
			notice2 := model.Notice{
				Receiver: temccomment.UserID,
				User:     model.User{},
				Sender:   user.UserID,
				Type:     "ccomment",
				Ntext:    msg.Content,
				Time:     time.Now(),
				Read:     false,
				Target:   newCcomment.CcommentID,
			}
			// 数据库创建一条通知
			db.Create(&notice2)
		}
	}
	// 如果是评论的回复

	var post model.Post
	db.Where("postID = ?", msg.PostID).First(&post)
	db.Model(&post).Update("comment_num", post.CommentNum+1)
	// 在这里设置 评论 的权重
	weightComment := float64(6)
	db.Model(&post).Update("heat", post.Heat+weightComment)
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
