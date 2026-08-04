package response

type PermissionWithAssignedRole struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Assigned bool   `json:"assigned"`
}
