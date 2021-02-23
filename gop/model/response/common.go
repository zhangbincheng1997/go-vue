package response

// PageResult struct
type PageResult struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}
