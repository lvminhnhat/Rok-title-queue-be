package model

import (
	"context"
	"strings"
	"sync"
	"time"
)

type Title struct {
	sync.Mutex                              // Thêm mutex
	Config                Config            `bson:"config" json:"config"`                               // Cấu hình của các title
	Queue                 []TitleAssignment `bson:"title_queue,omitempty" json:"title_queue,omitempty"` // Hàng đợi của các title
	TimeDone              TimeDone          `bson:"title_cooldown" json:"title_cooldown"`               // Thời gian để có thể thực thi title mới
	HomeKingdomMap        string            `bson:"home_kingdom_map" json:"home_kingdom_map"`
	LostKingdomMap        string            `bson:"lost_kingdom_map" json:"lost_kingdom_map"`
	UserDiscordChannelID  string            `bson:"user_discord_channel_id" json:"user_discord_channel_id"`
	AdminDiscordChannelID string            `bson:"admin_discord_channel_id" json:"admin_discord_channel_id"`
}

func (t *Title) NewTitle() *Title {
	var title Title
	title.Config = title.Config.NewConfig()
	return &title
}
func (t *Title) DeleteTitle(titleA TitleAssignment) {
	for i, title := range t.Queue {
		if title.PlayerID == titleA.PlayerID {
			// Kiểm tra nếu mảng chỉ có 1 phần tử
			if len(t.Queue) == 1 {
				// Nếu chỉ có 1 phần tử, gán mảng rỗng
				t.Queue = []TitleAssignment{}
			} else {
				// Nếu có nhiều phần tử, xóa phần tử tại vị trí i
				t.Queue = append(t.Queue[:i], t.Queue[i+1:]...)
			}
			break
		}
	}
}

func (t *Title) Finish(titleA TitleAssignment) {
	t.SetTimeDone(titleA.Title, int64(time.Now().UTC().Unix())+t.Config.GetTitleDuration(titleA.Title))
	t.DeleteTitle(titleA)
}
func (t *Title) Done(titleA TitleAssignment) {
	t.SetTimeDone(titleA.Title, int64(time.Now().UTC().Unix()))
	// xóa bỏ title đã hoàn thành
	t.DeleteTitle(titleA)
}
func (t *Title) SetTimeDone(title string, time int64) {
	switch title {
	case "duke":
		t.TimeDone.Duke = time
	case "architect":
		t.TimeDone.Architect = time
	case "scientist":
		t.TimeDone.Scientist = time
	case "justice":
		t.TimeDone.Justice = time
	}
}
func (t *Title) GetTimeDone(title string) int64 {
	switch title {
	case "duke":
		return t.TimeDone.Duke
	case "architect":
		return t.TimeDone.Architect
	case "scientist":
		return t.TimeDone.Scientist
	case "justice":
		return t.TimeDone.Justice
	default:
		return 0
	}
}

func (t *Title) GetMap(mapName string) string {
	switch mapName {
	case "home":
		return t.HomeKingdomMap
	case "lost":
		return t.LostKingdomMap
	default:
		return ""
	}
}
func (t *Title) GetTitleAssignment() (TitleAssignment, bool) {
	t.Lock()
	defer t.Unlock()

	// Dọn dẹp titles hết hạn trước
	t.CleanExpiredTitles()

	var titleA TitleAssignment
	currentTime := int64(time.Now().UTC().Unix())

	// Tìm title có thời gian chờ lâu nhất và đủ điều kiện
	for _, title := range t.Queue {
		timeDone := t.GetTimeDone(title.Title)

		if currentTime > timeDone {
			// Title đầu tiên thỏa điều kiện
			if titleA.PlayerID == "" {
				titleA = title
			} else {
				// So sánh thời gian thêm vào để lấy title chờ lâu nhất
				if title.TimeAdd < titleA.TimeAdd {
					titleA = title
				}
			}
		}
	}

	if titleA.PlayerID == "" {
		return TitleAssignment{}, false
	}
	return titleA, true
}
func (t *Title) AddTitle(title TitleAssignment) bool {
	t.Lock()
	defer t.Unlock()

	// Dọn dẹp titles hết hạn trước
	t.CleanExpiredTitles()

	// Thêm kiểm tra
	if title.PlayerID == "" || title.Title == "" {
		return false
	}
	title.SetTimeAdd()
	// chặn thêm title nếu title đã tồn tại trong hàng đợi
	if t.Config.GetTitleDuration(title.Title) == -1 {
		return false
	}
	for _, titleA := range t.Queue {
		if titleA.PlayerID == title.PlayerID {
			return false
		}
	}
	t.Queue = append(t.Queue, title)
	return true
}

func (t *Title) CleanExpiredTitles() {
	currentTime := int64(time.Now().UTC().Unix())
	newQueue := []TitleAssignment{}

	for _, title := range t.Queue {
		if currentTime <= title.TimeAdd+t.Config.GetTitleDuration(title.Title) {
			newQueue = append(newQueue, title)
		}
	}

	t.Queue = newQueue
}

func (t *Title) StartCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				t.Lock()
				t.CleanExpiredTitles()
				t.Unlock()
			}
		}
	}()
}

// Cấu trúc Config cho các title
type Config struct {
	Duke      int64 `bson:"duke, omitempty"`      // Thời gian giữ title Công tước
	Architect int64 `bson:"architect, omitempty"` // Thời gian giữ title Kiến trúc sư
	Scientist int64 `bson:"scientist, omitempty"` // Thời gian giữ title Nhà khoa học
	Justice   int64 `bson:"justice, omitempty"`   // Thời gian giữ title Công lý
}

func (c *Config) GetTitleDuration(title string) int64 {
	// chuyển đổi title về chữ thường
	title = strings.ToLower(title)
	switch title {
	case "duke":
		return c.Duke
	case "architect":
		return c.Architect
	case "scientist":
		return c.Scientist
	case "justice":
		return c.Justice
	default:
		return 0
	}
}

func (c *Config) SetTitleDuration(title string, duration int64) {
	// chuyển đổi title về chữ thường
	title = strings.ToLower(title)
	switch title {
	case "duke":
		c.Duke = duration
	case "architect":
		c.Architect = duration
	case "scientist":
		c.Scientist = duration
	case "justice":
		c.Justice = duration
	}
}
func (c *Config) NewConfig() Config {
	return Config{
		Duke:      60 * 2,
		Architect: 60 * 2,
		Scientist: 60 * 2,
		Justice:   60 * 5,
	}
}

// Cấu trúc cho hàng đợi Title
type TitleAssignment struct {
	PlayerID string `json:"player_id" bson:"player_id"` // ID của người chơi
	Title    string `json:"title" bson:"title"`         // Tên của title (Duke, Architect, Scientist, Justice)
	Local    Local  `json:"local" bson:"local"`         // Vị trí của người chơi
	TimeAdd  int64  `json:"time_add" bson:"time_add"`   // Thời gian thêm title vào hàng đợi
}

func (t *TitleAssignment) SetTimeAdd() {
	t.TimeAdd = int64(time.Now().UTC().Unix())
}

type Local struct {
	// map chỉ có thể có 2 dữ liệu là "home" hoặc "lost"
	Map string `json:"map"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type TimeDone struct {
	Duke      int64 `bson:"duke" json:"duke"`           // Thời gian hoàn thành title Công tước
	Architect int64 `bson:"architect" json:"architect"` // Thời gian hoàn thành title Kiến trúc sư
	Scientist int64 `bson:"scientist" json:"scientist"` // Thời gian hoàn thành title Nhà khoa học
	Justice   int64 `bson:"justice" json:"justice"`     // Thời gian hoàn thành title Công lý
}
