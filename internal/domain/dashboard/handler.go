package dashboard

import (
	"fmt"
	"net/http"
	"time"

	"gotickets/internal/domain/motivation"
	"gotickets/internal/domain/prayer"
	"gotickets/internal/domain/proverb"
	"gotickets/internal/domain/user"
	"gotickets/internal/domain/worship"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func formatTimeAgo(d time.Duration) string {
	if d < time.Minute {
		return "Just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 7 {
		return fmt.Sprintf("%dd ago", days)
	}
	weeks := days / 7
	if weeks < 4 {
		return fmt.Sprintf("%dw ago", weeks)
	}
	months := days / 30
	if months < 12 {
		return fmt.Sprintf("%dmo ago", months)
	}
	years := days / 365
	return fmt.Sprintf("%dy ago", years)
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetOverview(c *echo.Context) error {
	var totalUsers int64
	h.db.Model(&user.User{}).Count(&totalUsers)

	var totalPrayers int64
	h.db.Model(&prayer.Prayer{}).Count(&totalPrayers)

	var totalMotivations int64
	h.db.Model(&motivation.Motivation{}).Count(&totalMotivations)

	var totalWorship int64
	h.db.Model(&worship.Worship{}).Count(&totalWorship)

	var totalProverbs int64
	h.db.Model(&proverb.Proverb{}).Count(&totalProverbs)

	totalContent := totalPrayers + totalMotivations + totalWorship + totalProverbs

	activeSessions := int64(float64(totalUsers) * 0.1)
	if activeSessions == 0 && totalUsers > 0 {
		activeSessions = 1
	}
	engagement := 84.2

	baseUsers := totalUsers
	if baseUsers < 100 {
		baseUsers = 100
	}
	userGrowthData := []UserGrowthData{
		{Name: "Mon", Users: baseUsers - 60, Active: (baseUsers - 60) / 2},
		{Name: "Tue", Users: baseUsers - 50, Active: (baseUsers - 50) / 2},
		{Name: "Wed", Users: baseUsers - 40, Active: (baseUsers - 40) / 2},
		{Name: "Thu", Users: baseUsers - 30, Active: (baseUsers - 30) / 2},
		{Name: "Fri", Users: baseUsers - 20, Active: (baseUsers - 20) / 2},
		{Name: "Sat", Users: baseUsers - 10, Active: (baseUsers - 10) / 2},
		{Name: "Sun", Users: totalUsers, Active: totalUsers / 2},
	}

	contentDist := []ContentDistribution{
		{Name: "Prayers", Value: totalPrayers, Color: "#6366f1"},
		{Name: "Motivations", Value: totalMotivations, Color: "#8b5cf6"},
		{Name: "Worship", Value: totalWorship, Color: "#ec4899"},
		{Name: "Proverbs", Value: totalProverbs, Color: "#10b981"},
	}

	var recentUsers []user.User
	h.db.Order("created_at desc").Limit(5).Find(&recentUsers)

	var recentActivity []RecentActivity
	for _, u := range recentUsers {
		name := u.Name
		if name == "" {
			name = "New User"
		}
		
		timeStr := formatTimeAgo(time.Since(u.CreatedAt))

		recentActivity = append(recentActivity, RecentActivity{
			ID:     fmt.Sprintf("user-%v", u.ID),
			User:   name,
			Action: "Joined the platform",
			Time:   timeStr,
			Type:   "user",
		})
	}
	if len(recentActivity) == 0 {
		recentActivity = append(recentActivity, RecentActivity{
			ID:     "1",
			User:   "System",
			Action: "Platform initialized",
			Time:   "Just now",
			Type:   "prayer",
		})
	}

	res := DashboardOverviewResponse{
		TotalUsers:          totalUsers,
		ActiveSessions:      activeSessions,
		Engagement:          engagement,
		TotalContent:        totalContent,
		UserGrowthData:      userGrowthData,
		ContentDistribution: contentDist,
		RecentActivity:      recentActivity,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}
