// Package src is a collection of functions for juno-benchmark
package src

// An InitializeModuleRequest represents an example module initialization function from Juno
type InitializeModuleRequest struct {
	RequestID string `json:"requestId"`
	Type      int    `json:"type"`
	ModuleID  string `json:"moduleId"`
	Version   string `json:"version"`
}

// A DeclareFunctionRequest represents an example declare function for a module from Juno
type DeclareFunctionRequest struct {
	RequestID string `json:"requestId"`
	Type      int    `json:"type"`
	Function  string `json:"function"`
}

// A FunctionCallRequest represents an example module call function from Juno
type FunctionCallRequest struct {
	RequestID string    `json:"requestId"`
	Type      int       `json:"type"`
	Function  string    `json:"function"`
	Arguments arguments `json:"arguments"`
}

type arguments struct {
	Req string `json:"req"`
}
