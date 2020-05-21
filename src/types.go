// Package src is a collection of functions for juno-benchmark
package src

// An InitMsg represents an example module initialization function from Juno
type InitMsg struct {
	RequestID string `json:"requestId"`
	Type      int    `json:"type"`
	ModuleID  string `json:"moduleId"`
	Version   string `json:"version"`
}

// A Job represents an example module call function from Juno
type Job struct {
	RequestID string    `json:"requestId"`
	Type      int       `json:"type"`
	Function  string    `json:"function"`
	Arguments arguments `json:"arguments"`
}

type arguments struct {
	Req string `json:"req"`
}
