package entity

type Member struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func (m *Member) Valid() bool {
	return len(m.Email) > 0 && len(m.Pass) > 0
}
