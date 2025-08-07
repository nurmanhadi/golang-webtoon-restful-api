package dto

type DashboardResponse struct {
	TotalComic        int `json:"total_comic"`
	TotalChapter      int `json:"total_chapter"`
	TotalUser         int `json:"total_user"`
	TotalViews        int `json:"total_views"`
	TotalViewsDaily   int `json:"total_views_daily"`
	TotalViewsWeekly  int `json:"total_views_weekly"`
	TotalViewsMonthly int `json:"total_views_monthly"`
}
