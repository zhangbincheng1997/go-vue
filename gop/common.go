package main

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Status struct
type Status struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

var statusOptions = []Status{
	{0, "失败"},
	{1, "成功"},
	{2, "待定"},
}

// Status const
const (
	FAILURE = iota
	SUCCESS
	WAITING
)
