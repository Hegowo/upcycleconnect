package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"upcycleconnect/backend/models"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForumHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	prev := '-'
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prev = r
		} else if prev != '-' {
			b.WriteRune('-')
			prev = '-'
		}
	}
	return strings.Trim(b.String(), "-")
}

func (h *ForumHandler) currentUser(c *gin.Context) *models.User {
	u, _ := c.Get("user")
	user, _ := u.(*models.User)
	return user
}

func (h *ForumHandler) ListCategories(c *gin.Context) {
	var cats []models.ForumCategory
	if err := h.DB.Order("sort_order, id").Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	for i := range cats {
		var tc, rc int64
		h.DB.Model(&models.ForumThread{}).Where("category_id = ? AND deleted_at IS NULL", cats[i].ID).Count(&tc)
		h.DB.Table("forum_replies").
			Joins("JOIN forum_threads ON forum_replies.thread_id = forum_threads.id").
			Where("forum_threads.category_id = ? AND forum_threads.deleted_at IS NULL AND forum_replies.deleted_at IS NULL", cats[i].ID).
			Count(&rc)
		cats[i].ThreadCount = int(tc)
		cats[i].ReplyCount = int(rc)
	}
	c.JSON(http.StatusOK, gin.H{"data": cats})
}

func (h *ForumHandler) ShowCategory(c *gin.Context) {
	slug := c.Param("slug")
	var cat models.ForumCategory
	if err := h.DB.Where("slug = ?", slug).First(&cat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.ForumThread{}).Where("category_id = ?", cat.ID).Count(&total)

	var threads []models.ForumThread
	h.DB.Where("category_id = ?", cat.ID).Preload("User").
		Order("pinned DESC, updated_at DESC").Limit(limit).Offset(offset).Find(&threads)

	for i := range threads {
		var rc int64
		h.DB.Model(&models.ForumReply{}).Where("thread_id = ? AND deleted_at IS NULL", threads[i].ID).Count(&rc)
		threads[i].ReplyCount = int(rc)
		if threads[i].User != nil {
			threads[i].User.Password = ""
		}
	}

	c.JSON(http.StatusOK, gin.H{"category": cat, "threads": threads, "total": total, "page": page})
}

func (h *ForumHandler) ShowThread(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var thread models.ForumThread
	if err := h.DB.Preload("User").Preload("Category").First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.User != nil {
		thread.User.Password = ""
	}

	var replies []models.ForumReply
	h.DB.Where("thread_id = ?", thread.ID).Preload("User").Order("created_at ASC").Find(&replies)
	for i := range replies {
		if replies[i].User != nil {
			replies[i].User.Password = ""
		}
	}
	thread.ReplyCount = len(replies)

	var bans []models.ForumBan
	h.DB.Preload("User").Where("thread_id = ?", thread.ID).Find(&bans)
	for i := range bans {
		if bans[i].User != nil {
			bans[i].User.Password = ""
		}
	}

	c.JSON(http.StatusOK, gin.H{"thread": thread, "replies": replies, "bans": bans})
}

func (h *ForumHandler) CreateThread(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		CategoryID uint   `json:"category_id" binding:"required"`
		Title      string `json:"title" binding:"required,min=5,max=255"`
		Content    string `json:"content" binding:"required,min=10"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cat models.ForumCategory
	if err := h.DB.First(&cat, body.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_not_found"})
		return
	}

	thread := models.ForumThread{
		CategoryID: body.CategoryID,
		UserID:     user.ID,
		Title:      strings.TrimSpace(body.Title),
		Content:    body.Content,
	}
	if err := h.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"thread": thread})
}

func (h *ForumHandler) UpdateThread(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Title != "" {
		thread.Title = strings.TrimSpace(body.Title)
	}
	if body.Content != "" {
		thread.Content = body.Content
	}
	h.DB.Save(&thread)
	c.JSON(http.StatusOK, gin.H{"thread": thread})
}

func (h *ForumHandler) CloseThread(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	thread.Locked = !thread.Locked
	h.DB.Save(&thread)
	c.JSON(http.StatusOK, gin.H{"locked": thread.Locked})
}

func (h *ForumHandler) DeleteThread(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	h.DB.Delete(&thread)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) BanUser(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	threadID, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only_owner"})
		return
	}

	userID, _ := strconv.Atoi(c.Param("userId"))
	if uint(userID) == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_ban_yourself"})
		return
	}

	ban := models.ForumBan{
		ThreadID: uint(threadID),
		UserID:   uint(userID),
		BannedBy: user.ID,
	}
	h.DB.Where(models.ForumBan{ThreadID: uint(threadID), UserID: uint(userID)}).
		FirstOrCreate(&ban)

	var banWithUser models.ForumBan
	h.DB.Preload("User").First(&banWithUser, ban.ID)
	if banWithUser.User != nil {
		banWithUser.User.Password = ""
	}

	c.JSON(http.StatusCreated, gin.H{"ban": banWithUser})
}

func (h *ForumHandler) UnbanUser(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	threadID, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only_owner"})
		return
	}

	userID, _ := strconv.Atoi(c.Param("userId"))
	h.DB.Where("thread_id = ? AND user_id = ?", threadID, userID).Delete(&models.ForumBan{})
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) CreateReply(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread_not_found"})
		return
	}
	if thread.Locked {
		c.JSON(http.StatusForbidden, gin.H{"error": "thread_locked"})
		return
	}

	var ban models.ForumBan
	if err := h.DB.Where("thread_id = ? AND user_id = ?", thread.ID, user.ID).First(&ban).Error; err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "banned"})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply := models.ForumReply{
		ThreadID: thread.ID,
		UserID:   user.ID,
		Content:  body.Content,
	}
	if err := h.DB.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	h.DB.Model(&thread).Update("updated_at", reply.CreatedAt)

	h.DB.Preload("User").First(&reply, reply.ID)
	if reply.User != nil {
		reply.User.Password = ""
	}
	c.JSON(http.StatusCreated, gin.H{"reply": reply})
}

func (h *ForumHandler) UpdateReply(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var reply models.ForumReply
	if err := h.DB.First(&reply, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if reply.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply.Content = body.Content
	h.DB.Save(&reply)
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

func (h *ForumHandler) DeleteReply(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var reply models.ForumReply
	if err := h.DB.First(&reply, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	var thread models.ForumThread
	h.DB.First(&thread, reply.ThreadID)

	if reply.UserID != user.ID && thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	h.DB.Delete(&reply)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) CreateReport(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Type     string `json:"type" binding:"required"`
		TargetID uint   `json:"target_id" binding:"required"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Type != "thread" && body.Type != "reply" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_type"})
		return
	}

	if body.Type == "thread" {
		var thread models.ForumThread
		if h.DB.First(&thread, body.TargetID).Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "target_not_found"})
			return
		}
	} else {
		var reply models.ForumReply
		if h.DB.First(&reply, body.TargetID).Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "target_not_found"})
			return
		}
	}

	report := models.ForumReport{
		Type:       body.Type,
		TargetID:   body.TargetID,
		ReporterID: user.ID,
		Reason:     strings.TrimSpace(body.Reason),
		Status:     "pending",
	}
	h.DB.Create(&report)
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (h *ForumHandler) UploadMedia(c *gin.Context) {
	user := h.currentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no_file"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_type"})
		return
	}

	dir := "/uploads/forum"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage_error"})
		return
	}

	filename := fmt.Sprintf("%d_%d%s", user.ID, time.Now().UnixNano(), ext)
	path := filepath.Join(dir, filename)
	out, err := os.Create(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage_error"})
		return
	}
	defer out.Close()
	io.Copy(out, file)

	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("/uploads/forum/%s", filename)})
}

func (h *ForumHandler) AdminCreateCategory(c *gin.Context) {
	var body struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := models.ForumCategory{
		Name:        strings.TrimSpace(body.Name),
		Slug:        slugify(body.Name),
		Description: strings.TrimSpace(body.Description),
		Icon:        body.Icon,
		Color:       body.Color,
		SortOrder:   body.SortOrder,
	}
	if err := h.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"category": cat})
}

func (h *ForumHandler) AdminUpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cat models.ForumCategory
	if err := h.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Name != "" {
		cat.Name = strings.TrimSpace(body.Name)
		cat.Slug = slugify(body.Name)
	}
	cat.Description = body.Description
	cat.Color = body.Color
	cat.SortOrder = body.SortOrder
	h.DB.Save(&cat)
	c.JSON(http.StatusOK, gin.H{"category": cat})
}

func (h *ForumHandler) AdminDeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.DB.Delete(&models.ForumCategory{}, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) AdminDeleteThread(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.DB.Delete(&models.ForumThread{}, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) AdminDeleteReply(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.DB.Delete(&models.ForumReply{}, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ForumHandler) AdminPinThread(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	thread.Pinned = !thread.Pinned
	h.DB.Save(&thread)
	c.JSON(http.StatusOK, gin.H{"pinned": thread.Pinned})
}

func (h *ForumHandler) AdminLockThread(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var thread models.ForumThread
	if err := h.DB.First(&thread, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	thread.Locked = !thread.Locked
	h.DB.Save(&thread)
	c.JSON(http.StatusOK, gin.H{"locked": thread.Locked})
}

func (h *ForumHandler) AdminListThreads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 30
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.ForumThread{}).Unscoped().Count(&total)

	var threads []models.ForumThread
	h.DB.Unscoped().Preload("User").Preload("Category").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)

	for i := range threads {
		if threads[i].User != nil {
			threads[i].User.Password = ""
		}
		var rc int64
		h.DB.Model(&models.ForumReply{}).Where("thread_id = ?", threads[i].ID).Count(&rc)
		threads[i].ReplyCount = int(rc)
	}
	c.JSON(http.StatusOK, gin.H{"data": threads, "total": total})
}

func (h *ForumHandler) AdminListReports(c *gin.Context) {
	status := c.DefaultQuery("status", "pending")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 30
	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.ForumReport{}).Where("status = ?", status).Count(&total)

	var reports []models.ForumReport
	h.DB.Where("status = ?", status).Preload("Reporter").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&reports)

	for i := range reports {
		if reports[i].Reporter != nil {
			reports[i].Reporter.Password = ""
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": reports, "total": total})
}

func (h *ForumHandler) AdminResolveReport(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var report models.ForumReport
	if err := h.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	report.Status = "resolved"
	h.DB.Save(&report)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
