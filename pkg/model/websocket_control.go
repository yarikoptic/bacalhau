package model

type WSControl struct {
	// users can authenticate either by JWT tokens or username/password
	JWTToken string `json:"jwt_token"`
	Username string `json:"username"`
	Password string `json:"password"`

	ActionType string `json:"action_type"` // subscribe, unsubscribe
	TaskID     string `json:"task_id"`
}
