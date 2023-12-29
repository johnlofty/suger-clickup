package models

import (
	"fmt"
	"regexp"
	"strconv"
)

const TextFormat = "(username)[user_id] says: {content}"

type ClickupTaskCommentsResponse struct {
	Comments []ClickupTaskComment
}
type ClickupTaskComment struct {
	ID          string             `json:"id,omitempty"`
	CommentText string             `json:"comment_text,omitempty"`
	Date        string             `json:"date,omitempty"`
	User        ClickupCommentUser `json:"user,omitempty"`
}

type ClickupCommentUser struct {
	ID             int64  `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
}

type ClickupCommentParseResult struct {
	Username string
	UserID   int32
	Content  string
}

func (r *ClickupCommentParseResult) Format() string {
	return fmt.Sprintf("(%s)[%d] says: %s", r.Username, r.UserID, r.Content)
}

func ExtractComment(input string) (ClickupCommentParseResult, bool) {
	re := regexp.MustCompile(`\((.*?)\)\[(\d+)\] says: (.*)`)
	match := re.FindStringSubmatch(input)

	var result ClickupCommentParseResult
	if len(match) >= 4 {
		username := match[1]
		userIDStr := match[2]
		content := match[3]
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		result.Username = username
		result.UserID = int32(userID)
		result.Content = content
		return result, true

	}
	return result, false
}

type ClickupCreateTaskCommentRequest struct {
	CommentText string `json:"comment_text,omitempty"`
}

type ClickupCreateTaskCommentResponse struct {
	ID int64
}
