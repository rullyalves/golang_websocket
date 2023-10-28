package input

type CreateComplaintCategoryParamsInput struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentId,omitempty"`
}
