package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"loginTest/api"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"time"
	"sort"
	"unicode/utf8"
)

type PostMsg struct {
	UserTelephone string
	Title         string
	Content       string
	Partition     string
	Photos        string
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
	photos := requestPostMsg.Photos
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
		BrowseNum: 0,
		Heat:      0,
		PostTime:  time.Now(),
		Photos:    photos,
	}
	db.Create(&newPost)
	response.Response(c, http.StatusOK, 200, nil, "发帖成功")
}

type PostResponse struct {
	PostID        uint
	UserName      string
	UserTelephone string
	UserAvatar    string
	Title         string
	Content       string
	Like          int
	Comment       int
	Browse        int
	Heat          float64
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
	Photos        string
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
			UserAvatar:    user.AvatarURL,
			Title:         post.Title,
			Content:       post.Ptext,
			Like:          post.LikeNum,
			Comment:       post.CommentNum,
			Browse:        post.BrowseNum,
			Heat:          post.Heat,
			PostTime:      post.PostTime,
			IsSaved:       isSaved,
			IsLiked:       isLiked,
			Photos:        post.Photos,
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

type BrowseMsg struct {
	UserTelephone string
	PostID        uint
	// BrowseNum     int
}

func UpdateBrowseNum(c *gin.Context) {
	db := common.GetDB()
	var requestBrowseMsg BrowseMsg
	c.Bind(&requestBrowseMsg)
	userTelephone := requestBrowseMsg.UserTelephone
	postID := requestBrowseMsg.PostID
	// browseNum := requestBrowseMsg.BrowseNum 不用获取直接+1
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	db.Model(&post).Update("browse_num", post.BrowseNum+1)
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
	UserAvatar    string
	Title         string
	Content       string
	Like          int
	Comment       int
	Browse        int
	Heat          float64
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
	Photos        string
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
		UserAvatar:    user.AvatarURL,
		Title:         post.Title,
		Content:       post.Ptext,
		Like:          post.LikeNum,
		Comment:       post.CommentNum,
		PostTime:      post.PostTime,
		IsSaved:       isSaved,
		IsLiked:       isLiked,
		Browse:        post.BrowseNum,
		Heat:          post.Heat,
		Photos:        post.Photos,
	}
	c.JSON(http.StatusOK, postDetailsResponse)
}

func UploadPhotos(c *gin.Context) {
	//UserID := c.PostForm("UserID")
	// 获取前端传过来 图片
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	// 文件保存路径和文件名可以根据实际情况修改
	// 文件名我们采用了当前时间戳和原始文件名的组合，以避免文件名冲突
	// 时间戳采用 nanoseconds 级别，可以几乎确保每个文件名都是唯一的
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename
	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 更新 Post 的 Photos 字段
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	fileURL := "http://localhost:8080/uploads/" + filename

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL, "message": "上传成功"})
}

func UploadZip(c *gin.Context) {
	const maxUploadSize = 100 << 20 // 100 MB

	// 获取前端传过来的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件太大，不能超过10MB"})
		return
	}

	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}
	defer fileBytes.Close()

	buffer := make([]byte, 512)
	_, err = fileBytes.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}

	if http.DetectContentType(buffer) != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是zip格式"})
		return
	}

	// 文件保存路径和文件名可以根据实际情况修改
	// 文件名我们采用了当前时间戳和原始文件名的组合，以避免文件名冲突
	// 时间戳采用 nanoseconds 级别，可以几乎确保每个文件名都是唯一的
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename

	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 更新 Post 的 Photos 字段
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	fileURL := "http://localhost:8080/uploads/" + filename

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"zipURL": fileURL, "message": "上传成功"})
}

func SubmitFeedback(c *gin.Context) {
	db := common.GetDB()

	// Create a struct to hold the incoming JSON body
	var feedbackInput struct {
		Ftext      string `json:"ftext"`
		Attachment string `json:"attachment"`
	}

	// Bind the incoming JSON to the struct
	if err := c.BindJSON(&feedbackInput); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Invalid request body")
		return
	}

	// Create a new feedback entry
	feedback := model.Feedback{
		Ftext:      feedbackInput.Ftext,
		Attachment: feedbackInput.Attachment,
	}

	db.Create(&feedback)

	if db.NewRecord(feedback) {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Failed to submit feedback")
		return
	}

	// Convert to JSON and respond
	response.Success(c, gin.H{
		"feedbackID": feedback.FeedbackID,
		"ftext":      feedback.Ftext,
		"attachment": feedback.Attachment,
	}, "Feedback submitted successfully")
}

func GetAllFeedback(c *gin.Context) {
	db := common.GetDB()

	var feedbacks []model.Feedback
	db.Find(&feedbacks)

	if len(feedbacks) == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "No feedback found")
		return
	}

	response.Success(c, gin.H{"feedbacks": feedbacks}, "Feedback retrieved successfully")
}

//func GetFeedback(c *gin.Context) {
//	db := common.GetDB()
//
//	feedbackID, err := strconv.Atoi(c.PostForm("feedbackID"))
//	if err != nil {
//		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid feedback ID")
//		return
//	}
//
//	var feedback model.Feedback
//	if err := db.First(&feedback, feedbackID).Error; err != nil {
//		if gorm.IsRecordNotFoundError(err) {
//			response.Response(c, http.StatusNotFound, 404, nil, "Feedback not found")
//		} else {
//			response.Response(c, http.StatusInternalServerError, 500, nil, "Database error")
//		}
//		return
//	}
//
//	response.Success(c, gin.H{"feedback": feedback}, "Feedback retrieved successfully")
//}

type HeatResponse struct {
	PostID uint
	Title  string
	Heat   float64
}

type PostStats struct {
	PostID     int
	LikeNum    int
	CommentNum int
	BrowseNum  int
	HeatValue  float64
	Title      string
}

func CalculateHeat(c *gin.Context) {
	// 获取所有帖子的 postID, likenum, commentnum, browsenum
	db := common.GetDB()
	// 从数据库中获取所有的帖子，并将结果存储在posts切片中。
	var posts []model.Post
	db.Find(&posts)
	// 创建了一个postStats切片，用于存储所有帖子的统计信息。
	var postStats []PostStats
	// 通过遍历posts切片，我们逐个取出每个帖子的信息，然后创建一个新的PostStats结构体，
	// 并将帖子的postID、likenum、commentnum和browsenum赋值给对应的字段。
	for _, post := range posts {
		var stats PostStats
		stats.PostID = post.PostID
		stats.LikeNum = post.LikeNum
		stats.CommentNum = post.CommentNum
		stats.BrowseNum = post.BrowseNum
		stats.Title = post.Title
		postStats = append(postStats, stats)
	}
	// postStats 列表中包含了所有帖子的 postID, likenum, commentnum, browsenum
	// 定义权重
	weightLike := 3
	weightComment := 5
	weightBrowse := 2
	// 计算每个帖子的热度
	for i := range postStats {
		heatValue := float64(
			postStats[i].LikeNum*weightLike +
				postStats[i].CommentNum*weightComment +
				postStats[i].BrowseNum*weightBrowse)
		postStats[i].HeatValue = heatValue
		// 更新post表的heat属性
		db.Model(&model.Post{}).Where("postID = ?", postStats[i].PostID).Update("heat", heatValue)
	}
	// 对 postStats 切片按热度进行排序
	sort.Slice(postStats, func(i, j int) bool {
		return postStats[i].HeatValue > postStats[j].HeatValue
	})
	var heatResponsesTop10 []HeatResponse
	// 只返回前10个帖子
	for i := 0; i < 10 && i < len(postStats); i++ {
		post := postStats[i]
		heatResponse := HeatResponse{
			PostID: uint(post.PostID),
			Title:  post.Title,
			Heat:   post.HeatValue,
		}
		heatResponsesTop10 = append(heatResponsesTop10, heatResponse)
	}
	c.JSON(http.StatusOK, heatResponsesTop10)
}