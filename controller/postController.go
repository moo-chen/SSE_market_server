package controller

import (
	"fmt"
	"image"
	"loginTest/api"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

type PostMsg struct {
	UserTelephone string
	Title         string
	Content       string
	Partition     string
	Photos        string
	TagList       string
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
	tagList := requestPostMsg.TagList
	tags := strings.Split(tagList, "|")
	tagString := strings.Join(tags, ",")
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
	if user.UserID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}

	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}

	currentDateTime := time.Now()
	if user.Banend.After(currentDateTime) {
		response.Response(c, http.StatusBadRequest, 400, nil, "你尚处于禁言状态中，不得发帖")
		return
	}

	newPost := model.Post{
		UserID:     int(user.UserID),
		Partition:  partition,
		Title:      title,
		Ptext:      content,
		LikeNum:    0,
		CommentNum: 0,
		BrowseNum:  0,
		Heat:       0,
		PostTime:   time.Now(),
		Photos:     photos,
		Tag:        tagString,
	}
	db.Create(&newPost)
	response.Response(c, http.StatusOK, 200, nil, "发帖成功")
}

type PostResponse struct {
	PostID        uint
	UserName      string
	UserScore     int
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
	Tag           string
}

type BrowseMeg struct {
	UserTelephone string
	Partition     string
	Searchinfo    string
	Searchsort    string //用于分表查询 分为home,history,save三种
	Limit         int
	Offset        int
}

func Browse(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestBrowseMeg BrowseMeg
	c.Bind(&requestBrowseMeg)
	userTelephone := requestBrowseMeg.UserTelephone
	partition := requestBrowseMeg.Partition
	searchinfo := requestBrowseMeg.Searchinfo
	searchsort := requestBrowseMeg.Searchsort
	limit := requestBrowseMeg.Limit
	offset := requestBrowseMeg.Offset
	// 查询用户
	var temUser model.User
	db.Where("phone = ?", userTelephone).First(&temUser)
	if temUser.UserID == 0 {
		return
	}
	// posts是查询的原数据,postResponses是post基础上添加了用户点赞和删除的信息
	var posts []model.Post
	var postResponses []PostResponse

	// 如果是收藏"save"查询,首先查询psave,再查询post
	if searchsort == "save" {
		// saves是用户的收藏列表
		var saves []model.Psave
		db.Order("psaveID DESC").Offset(offset).Limit(limit).Where("userID = ?", temUser.UserID).Find(&saves)
		// 根据收藏列表得到帖子列表
		for _, save := range saves {
			var post model.Post
			db.Where("postID = ?", save.PtargetID).First(&post)
			if post.PostID != 0 {
				posts = append(posts, post)
			}
		}
		// 对每个帖子查询是否点赞
		for _, post := range posts {
			isLiked := false
			var like model.Plike
			db.Where("userID = ? AND ptargetID = ?", temUser.UserID, post.PostID).First(&like)
			if like.PlikeID != 0 {
				isLiked = true
			}
			var user model.User
			if post.UserID == 0 {
				user.Name = "管理员"
				user.Phone = "11111111111"
			} else {
				db.Where("userID = ?", post.UserID).First(&user)
			}
			postResponse := PostResponse{
				PostID:        uint(post.PostID),
				UserName:      user.Name,
				UserScore:     user.Score,
				UserTelephone: user.Phone,
				UserAvatar:    user.AvatarURL,
				Title:         post.Title,
				Content:       post.Ptext,
				Like:          post.LikeNum,
				Comment:       post.CommentNum,
				Browse:        post.BrowseNum,
				Heat:          post.Heat,
				PostTime:      post.PostTime,
				IsSaved:       true,
				IsLiked:       isLiked,
				Photos:        post.Photos,
				Tag:           post.Tag,
			}
			postResponses = append(postResponses, postResponse)
		}
	} else {
		// "home"查询
		if searchsort == "home" {
			if partition == "主页" || len(partition) == 0 {
				if len(searchinfo) == 0 {
					db.Order("postID DESC").Offset(offset).Limit(limit).Find(&posts)
				} else {
					db.Order("postID DESC").Offset(offset).Limit(limit).Where("title LIKE ? OR ptext LIKE ? OR tag LIKE ?", "%"+searchinfo+"%", "%"+searchinfo+"%", "%"+searchinfo+"%").Find(&posts)
				}
			} else {
				db.Order("postID DESC").Offset(offset).Limit(limit).Find(&posts, "`partition` = ?", partition)
			}
		} else if searchsort == "history" {
			// historys是用户的发帖记录
			db.Order("postID DESC").Offset(offset).Limit(limit).Find(&posts, "`userID` = ?", temUser.UserID)
		}
		// 对每个帖子查询是否点赞和收藏
		for _, post := range posts {
			if post.PostID == 0 {
				continue
			}
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
			if post.UserID == 0 {
				user.Name = "管理员"
				user.Phone = "11111111111"
			} else {
				db.Where("userID = ?", post.UserID).First(&user)
			}
			postResponse := PostResponse{
				PostID:        uint(post.PostID),
				UserName:      user.Name,
				UserScore:     user.Score,
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
				Tag:           post.Tag,
			}
			postResponses = append(postResponses, postResponse)
		}
	}
	c.JSON(http.StatusOK, postResponses)
}

func GetPostNum(c *gin.Context) {
	db := common.GetDB()
	// 获取参数
	var requestBrowseMeg BrowseMeg
	c.Bind(&requestBrowseMeg)
	userTelephone := requestBrowseMeg.UserTelephone
	partition := requestBrowseMeg.Partition
	searchinfo := requestBrowseMeg.Searchinfo
	searchsort := requestBrowseMeg.Searchsort
	var count int
	// 下面的searchsort分为home,save,history三类
	if searchsort == "home" {
		if partition == "主页" || len(partition) == 0 {
			if len(searchinfo) == 0 {
				db.Model(&model.Post{}).Count(&count)
			} else {
				db.Model(&model.Post{}).
					Where("title LIKE ? OR ptext LIKE ? OR tag LIKE ?", "%"+searchinfo+"%", "%"+searchinfo+"%", "%"+searchinfo+"%").
					Where("postId != ?", 0).
					Count(&count)
			}
		} else {
			db.Model(&model.Post{}).Where("`partition` = ?", partition).Count(&count)
		}
	} else {
		// 查询用户
		var user model.User
		db.Where("phone = ?", userTelephone).First(&user)
		if user.UserID == 0 {
			return
		}
		if searchsort == "save" {
			db.Model(&model.Psave{}).
				Joins("INNER JOIN posts ON psaves.ptargetID = posts.postID").
				Where("psaves.userID = ? AND posts.postID != ?", user.UserID, 0).
				Count(&count)
		} else if searchsort == "history" {
			db.Model(&model.Post{}).
				Where("userID = ?", user.UserID).
				Where("postId != ?", 0).
				Count(&count)
		}
	}
	// 将结果返回给客户端
	c.JSON(http.StatusOK, gin.H{
		"Postcount": count,
	})
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
	if user.UserID == 0 {
		return
	}
	var post model.Post
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
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
	if user.UserID == 0 {
		return
	}
	var post model.Post
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
	if isLiked {
		// var liketime model.Plike
		// var Time time.Time
		// db.Where("time = ?", Time).First(&liketime)
		var like model.Plike
		db.Where("userID = ? AND ptargetID = ?", user.UserID, post.PostID).First(&like)
		// fmt.Println("likeID ", like.PlikeID)
		if like.PlikeID != 0 {
			currentTime := time.Now()
			fmt.Println("currentTime", currentTime)
			fmt.Println("likeTime", like.Time)
			timedif := currentTime.Sub(like.Time)
			hours := timedif.Hours()
			days := int(hours / 24)
			fmt.Println("days ", days)
			if days > 0 {
				weightlikePower := math.Pow(0.5, float64(days))
				weightLike := float64(3)
				deleteHeat := math.Pow(weightLike, weightlikePower)
				fmt.Println("deleteHeat ", deleteHeat)
				db.Model(&post).Update("heat", post.Heat-deleteHeat)
			} else {
				weightLike := float64(3)
				db.Model(&post).Update("heat", post.Heat-weightLike)
			}
			db.Model(&post).Update("like_num", post.LikeNum-1)
			db.Delete(&like)
		}
	} else {
		newLike := model.Plike{
			UserID:    user.UserID,
			PtargetID: post.PostID,
			Time:      time.Now(),
		}
		if newLike.UserID != 0 && newLike.PtargetID != 0 {
			db.Model(&post).Update("like_num", post.LikeNum+1)
			// 在这里设置 点赞 的权重
			weightLike := float64(3)
			db.Model(&post).Update("heat", post.Heat+weightLike)
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
	if user.UserID == 0 {
		return
	}
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
	db.Model(&post).Update("browse_num", post.BrowseNum+1)
	// 在这里设置 浏览 的权重
	weightBrowse := float64(1)
	db.Model(&post).Update("heat", post.Heat+weightBrowse)
	newBrowse := model.Pbrowse{
		UserID:    user.UserID,
		PtargetID: post.PostID,
		Time:      time.Now(),
	}
	db.Create(&newBrowse)
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
	if post.PostID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "需要删除的帖子不存在")
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != post.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "帖子删除成功"})
}

type Reportmsg struct {
	TargetID      uint
	Targettype    string
	UserTelephone string
	Reason        string
}

func SubmitReport(c *gin.Context) {
	db := common.GetDB()
	var reportmsg Reportmsg
	c.Bind(&reportmsg)
	TargetID := reportmsg.TargetID
	Targettype := reportmsg.Targettype
	userTelephone := reportmsg.UserTelephone
	Reason := reportmsg.Reason
	if len(Reason) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "举报内容不能为空")
		return
	}
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	newSue := model.Sue{
		Targettype: Targettype,
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
	UserScore     int
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
	Tag           string
}

type PostDetailsMsg struct {
	UserTelephone string `json:"userTelephone"`
	PostID        uint   `json:"postID"`
}

func ShowDetails(c *gin.Context) {
	db := common.GetDB()
	var requestPostDetailsMsg PostDetailsMsg
	c.Bind(&requestPostDetailsMsg)
	userTelephone := requestPostDetailsMsg.UserTelephone
	postID := requestPostDetailsMsg.PostID
	var temUser model.User
	db.Where("phone = ?", userTelephone).First(&temUser)
	if temUser.UserID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "无法解析当前用户")
		return
	}
	if postID == 0 {
		response.Response(c, http.StatusBadRequest, 404, nil, "接收到的postID为空")
		return
	}
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
	if post.PostID == 0 {
		return
	}
	var user model.User
	if post.UserID == 0 {
		user.Name = "管理员"
		user.Phone = "11111111111"
	} else {
		db.Where("userID = ?", post.UserID).First(&user)
	}
	postDetailsResponse := PostDetailsResponse{
		PostID:        uint(post.PostID),
		UserName:      user.Name,
		UserScore:     user.Score,
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
		Tag:           post.Tag,
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

	// 生成略缩图
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
		return
	}
	defer src.Close()
	// 读取上传文件的内容，并解码为 image.Image 对象
	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解码失败"})
		return
	}

	// 对图像进行缩略图处理
	resizedImage := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "D:\\SSE_market\\test\\serverTest2\\SSE_market_server\\public\\resized/" + filename
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 更新 Post 的 Photos 字段
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	fileURL := "https://localhost:8080/resized/" + filename

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL, "message": "上传成功"})
}

func UploadZip(c *gin.Context) {
	const maxUploadSize = 10 << 20 // 10 MB

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
	fileURL := "https://localhost:8080/uploads/" + filename

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
		Time:       time.Now(),
		Status:     "wait",
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
	// 对feedbacks按照时间进行排序
	sort.Slice(feedbacks, func(i, j int) bool {
		return feedbacks[j].Time.Before(feedbacks[i].Time)
	})
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

func CalculateHeat(c *gin.Context) {
	// 获取所有帖子的 postID, likenum, commentnum, browsenum
	db := common.GetDB()
	// 从数据库中获取所有的帖子，并将结果存储在posts切片中。
	var posts []model.Post
	db.Find(&posts)
	// 对 postStats 切片按热度进行排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Heat > posts[j].Heat
	})
	var heatResponsesTop10 []HeatResponse
	// 只返回前10个帖子
	for i := 0; i < 10 && i < len(posts); i++ {
		post := posts[i]
		heatResponse := HeatResponse{
			PostID: uint(post.PostID),
			Title:  post.Title,
			Heat:   post.Heat,
		}
		heatResponsesTop10 = append(heatResponsesTop10, heatResponse)
	}
	c.JSON(http.StatusOK, heatResponsesTop10)
}
