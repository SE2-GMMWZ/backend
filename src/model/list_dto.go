package model

type UserListResponse struct {
	Users       []User `json:"users"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
}

type BookingListResponse struct {
	Bookings    []Booking `json:"bookings"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
}

type CommentListResponse struct {
	Comments    []Comment `json:"comments"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
}

type GuideListResponse struct {
	Guides      []Guide `json:"guides"`
	CurrentPage int     `json:"current_page"`
	TotalPages  int     `json:"total_pages"`
}

type DockingSpotListResponse struct {
	DockingSpots []DockingSpot `json:"docking_spots"`
	CurrentPage  int           `json:"current_page"`
	TotalPages   int           `json:"total_pages"`
}
