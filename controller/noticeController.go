package controller

import (
	"github.com/gin-gonic/gin"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"strconv"
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
type NoticeGet struct {
	NoticeList     []NoticeResponse `json:"noticeList"`
	LastID         int              `json:"lastID"`
	TotalNum       int              `json:"totalNum"`
	UnreadTotalNum int              `json:"unreadTotalNum"`
}
type NoticeNumResponse struct {
	TotalNum       int `json:"totalNum"`
	UnreadTotalNum int `json:"unreadTotalNum"`
	ReadTotalNum   int `json:"readTotalNum"`
}

func GetNotice(c *gin.Context) {
	db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	pageSizeStr := c.Query("pageSize")
	requireIDStr := c.Query("requireID")
	readStr := c.Query("read")
	requireID, _ := strconv.Atoi(requireIDStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	read, _ := strconv.Atoi(readStr)
	var user model.User
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "游客无法访问通知")
		return
	} else {
		user = value.(model.User)
	}
	//执行操作，分页返回通知
	var notices []model.Notice
	var total int
	var unreadTotal int
	// requireID==0说明是首次查询,返回pagesize条通知,计算totalNum
	if requireID == 0 {
		db.Model(&model.Notice{}).Where("receiver =?", user.UserID).Count(&total)
		db.Model(&model.Notice{}).Where("receiver =? AND `read` =?", user.UserID, read).Count(&unreadTotal)
		db.Where("receiver =? AND `read` =?", user.UserID, read).Order("noticeID DESC").Limit(pageSize).Find(&notices)
	} else { //否则查询比requireID小的通知
		db.Where("receiver =? AND `read` =? AND noticeID < ?", user.UserID, read, requireID).Order("noticeID DESC").Limit(pageSize).Find(&notices)
	}
	if len(notices) == 0 {
		response.Response(c, http.StatusOK, 201, nil, "没有更多通知")
		return
	} else if len(notices) < 5 {

	}
	var noticeGet NoticeGet
	noticeGet.TotalNum = total
	noticeGet.LastID = notices[len(notices)-1].NoticeID
	for _, notice := range notices {
		var temuser model.User
		db.Where("userID =?", notice.Sender).First(&temuser)
		var tempcomment model.Pcomment
		if notice.Type == "ccomment" {
			var temccoment model.Ccomment
			db.Where("ccommentID=?", notice.Target).First(&temccoment)
			db.Where("pcommentID =?", temccoment.CtargetID).First(&tempcomment)
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
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
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
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
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
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
	c.JSON(http.StatusOK, &noticeGet)
}

func GetNoticeNum(c *gin.Context) {
	db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	var user model.User
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "游客无法获得通知数量")
		return
	} else {
		user = value.(model.User)
	}
	var noticeNum NoticeNumResponse
	var readNum int
	var unreadNum int
	db.Model(&model.Notice{}).Where("receiver =? AND read =0", user.UserID).Count(&unreadNum)
	db.Model(&model.Notice{}).Where("receiver =? AND read =1", user.UserID).Count(&readNum)
	noticeNum.TotalNum = readNum + unreadNum
	noticeNum.ReadTotalNum = readNum
	noticeNum.UnreadTotalNum = unreadNum
	c.JSON(http.StatusOK, &noticeNum)
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
