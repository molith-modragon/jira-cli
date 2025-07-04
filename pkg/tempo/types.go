package tempo

import "time"

// Team represents a Tempo team.
type Team struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	LeadID      string    `json:"leadId"`
	Lead        *User     `json:"lead"`
	Members     []*Member `json:"members"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// User represents a Tempo user.
type User struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Username     string `json:"username"`
	Active       bool   `json:"active"`
}

// Member represents a team member.
type Member struct {
	User        *User     `json:"user"`
	Role        string    `json:"role"`
	Commitment  int       `json:"commitment"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}