package middleware

// UserContext define os dados do usuário autenticado via JWT
type UserContext struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

// AgentContext define os dados do Agente/Bot autenticado via API Key
type AgentContext struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	KeyID     string `json:"key_id"`
}
