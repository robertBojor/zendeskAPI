package zendesk

import "time"

type CreateTaskOptions struct {
	Content      string `json:"content"`
	DueDate      string `json:"due_date"`
	OwnerID      int64  `json:"owner_id"`
	ResourceType string `json:"resource_type"`
	ResourceID   int64  `json:"resource_id"`
	Completed    bool   `json:"completed"`
	RemindAt     string `json:"remind_at"`
}
type createTaskRequest struct {
	Data CreateTaskOptions `json:"data"`
	Meta struct{
		Type string `json:"type"`
	}
}
type CreateTaskResponse struct {
	Data struct {
		ID           int64      `json:"id"`
		CreatorID    int64      `json:"creator_id"`
		OwnerID      int64      `json:"owner_id"`
		ResourceType string     `json:"resource_type"`
		ResourceID   int64      `json:"resource_id"`
		Completed    bool       `json:"completed"`
		CompletedAt  *time.Time `json:"completed_at"`
		DueDate      *time.Time `json:"due_date"`
		Overdue      bool       `json:"overdue"`
		RemindAt     *time.Time `json:"remind_at"`
		Content      string     `json:"content"`
		CreatedAt    *time.Time `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	} `json:"data"`
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}

func (z *API) CreateTask(options *CreateTaskOptions) {
	if options == nil {
		return
	}
	body := createTaskRequest{
		Data: *options,
	}
	body.Meta.Type = "task"
	z.createRequest("POST", "/v2/tasks", body).execute()
}