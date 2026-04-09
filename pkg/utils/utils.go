package utils

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"math"
	mrand "math/rand"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Habeebamoo/Clivo/server/internal/config"
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Slugify(title string) string {
	slug := strings.ToLower(title)

	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")
	return slug
}

func GenerateArticleSlug(authorUsername string, articleTitle string) string {
	mrand.Seed(time.Now().UnixNano())
	randomNum := mrand.Intn(1000)
	slug := Slugify(articleTitle)
	return fmt.Sprintf("%s/%s-%d", authorUsername, slug, randomNum)
}

func FormatText(text string) string {
	caser := cases.Title(language.English)

	firstChar := text[:1]
	restChars := text[1:]

	return caser.String(firstChar) + restChars
}

func SetCookies(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "auth_token",
		Value: token,
		Path: "/",
		Domain: "",
		Expires: time.Now().Add(24*time.Hour),
		MaxAge: 3600,
		Secure: true, //true for production
		HttpOnly: true,
	})
}

func RemoveCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "auth_token",
		Value: "",
		Path: "/",
		Domain: "",
		MaxAge: -1,
		Secure: true, // true for production
		HttpOnly: true,
	})
}

func GenerateRandomId() string {
	b := make([]byte, 16)
	_, err := crand.Read(b)
	if err != nil {
		log.Fatal("failed to read bytes")
	}

	return hex.EncodeToString(b)
}

func GetUserProfile(username string) string {
	clientOrigin, err := config.Get("CLIENT_URL")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/%s", clientOrigin, username)
}

func GetUsernameFromEmail(email string) string {
	local := strings.SplitN(email, "@", 2)[0]
	return "@" + local
}

func GenerateUniqueUsername(base string, exists func(string) bool) string {
	username := base
	i := 1

	for exists(username) {
		username = fmt.Sprintf("%s%d", username, i)
		i++
	}

	return username
}

// Entry point
func CleanEditorJS(raw json.RawMessage) (json.RawMessage, error) {
	var data interface{}

	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	cleaned := cleanRecursive(data)

	out, err := json.Marshal(cleaned)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Recursive cleaner
func cleanRecursive(v interface{}) interface{} {
	switch val := v.(type) {

	case map[string]interface{}:
		for k, v2 := range val {
			val[k] = cleanRecursive(v2)
		}
		return val

	case []interface{}:
		for i, v2 := range val {
			val[i] = cleanRecursive(v2)
		}
		return val

	case string:
		return cleanString(val)

	default:
		return val
	}
}

// Actual string cleaning logic
func cleanString(s string) string {
	// Decode HTML entities (&nbsp; -> space)
	s = html.UnescapeString(s)

	// Replace non-breaking space (unicode 0xA0) with normal space
	s = strings.ReplaceAll(s, "\u00A0", " ")

	// Trim trailing spaces
	s = strings.TrimRight(s, " ")

	return s
}

func MergeUniqueArticle(a1, a2 []models.Article) []models.Article {
	seen := make(map[string]bool)
	results := []models.Article{}

	for _, article := range append(a1, a2...) {
		if !seen[article.ArticleId] {
			seen[article.ArticleId] = true
			results = append(results, article)
		}
	}

	return results
}

func GetArticleReadTime(jsonContent string) int {
	var content models.EditorJSContent
	err := json.Unmarshal([]byte(jsonContent), &content)
	if err != nil {
		return 1
	}

	wordCount := 0
	for _, block := range content.Blocks {
		if block.Type == "paragraph" || block.Type == "header" {
			words := strings.Fields(block.Data.Text)
			wordCount += len(words)
		}
	}

	minutes := float64(wordCount) / 200.0
	if minutes < 1 {
		return 1
	}
	
	return int(math.Ceil(minutes))
}


func GetTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%d mins ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 12*30*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365)))
	}
}

func UploadImage(file multipart.File) (string, error) {
	key, _ := config.Get("IK_PRIVATE_KEY")
	if key == "" {
		return "", fmt.Errorf("env variables missing")
	}

	client := imagekit.NewClient(
		option.WithPrivateKey(key),
	)

	resp, err := client.Files.Upload(context.TODO(), imagekit.FileUploadParams{
		File: file,
		FileName: fmt.Sprintf("%s.jpg", GenerateRandomId()),
	})

	if err != nil {
		return "", fmt.Errorf("upload error")
	}

	return resp.URL, nil
}

func StripHtml(s string) string {
	re := regexp.MustCompile("<.*?>")
	return re.ReplaceAllString(s, "")
}

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max] + "..."
}

func ExtractExcerpt(content json.RawMessage) string {
	var data struct {
		Blocks []struct {
			Type string `json:"type"`
			Data struct {
				Text string `json:"text"`
			} `json:"data"`
		} `json:"blocks"`
	}

	_ = json.Unmarshal(content, &data)

	for _, b := range data.Blocks {
		if b.Type == "paragraph" && b.Data.Text != "" {
			text := StripHtml(b.Data.Text)
			return Truncate(text, 180)
		}
	}

	return "Read this article on Clivo"
}

func IsBot(ua string) bool {
	bots := []string{
		"facebookexternalhit",
		"Twitterbot",
		"LinkedInBot",
		"WhatsApp",
		"Slackbot",
		"Discordbot",
		"Googlebot",
	}

	ua = strings.ToLower(ua)

	for _, b := range bots {
		if strings.Contains(ua, strings.ToLower(b)) {
			return true
		}
	}

	return false
}