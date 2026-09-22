package dashboard

type DashboardOverviewResponse struct {
	TotalUsers          int64                 `json:"totalUsers"`
	ActiveSessions      int64                 `json:"activeSessions"`
	Engagement          float64               `json:"engagement"`
	TotalContent        int64                 `json:"totalContent"`
	UserGrowthData      []UserGrowthData      `json:"userGrowthData"`
	ContentDistribution []ContentDistribution `json:"contentDistribution"`
	RecentActivity      []RecentActivity      `json:"recentActivity"`
}

type UserGrowthData struct {
	Name   string `json:"name"`
	Users  int64  `json:"users"`
	Active int64  `json:"active"`
}

type ContentDistribution struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
	Color string `json:"color"`
}

type RecentActivity struct {
	ID     string `json:"id"`
	User   string `json:"user"`
	Action string `json:"action"`
	Time   string `json:"time"`
	Type   string `json:"type"`
}
