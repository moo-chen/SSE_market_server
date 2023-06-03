package controller

import (
	"github.com/gin-gonic/gin"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"time"
)

type NoticeResponse struct {
	NoticeID     int       `json:"noticeID"`
	ReceiverName string    `json:"receiverName"`
	SenderName   string    `json:"senderName"`
	SenderAvatar string    `json:"senderAvatar"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	Read         bool      `json:"read"`
	PostID       int       `json:"postID"`
	Target       int       `json:"target"`
	PcommentID   int       `json:"pcommentID"`
	Time         time.Time `json:"time"`
}

func GetNotice(c *gin.Context) {
	db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	var user model.User
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "游客无法访问通知")
		return
	} else {
		user = value.(model.User)
	}
	//执行操作
	var notices []model.Notice
	db.Find(&notices, "receiver =?", user.UserID)
	var noticeResponse []NoticeResponse
	for _, notice := range notices {
		var temuser model.User
		db.Where("userID =?", notice.Sender).First(&temuser)
		var tempcomment model.Pcomment
		if notice.Type == "ccomment" {
			var temccoment model.Ccomment
			db.Where("ccommentID=?", notice.Target).First(&temccoment)
			db.Where("pcommentID =?", temccoment.CtargetID).First(&tempcomment)
			noticeResponse = append(noticeResponse, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				PcommentID:   temccoment.CtargetID,
				PostID:       tempcomment.PtargetID,
				Time:         temccoment.Time,
			})
		} else if notice.Type == "pcomment" {
			db.Where("pcommentID =?", notice.Target).First(&tempcomment)
			noticeResponse = append(noticeResponse, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				PostID:       tempcomment.PtargetID,
				Time:         tempcomment.Time,
			})
		} else {
			noticeResponse = append(noticeResponse, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
			})

		}

	}
	c.JSON(http.StatusOK, &noticeResponse)
}
func ReadNotice(c *gin.Context) {
	noticeID := c.Param("noticeID")
	db := common.GetDB()
	// Update notice data in database to set 'read' field to true
	err := db.Model(&model.Notice{}).Where("noticeID = ?", noticeID).Update("read", true).Error
	if err != nil {
		c.JSON(500, gin.H{"更新已读出现error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
