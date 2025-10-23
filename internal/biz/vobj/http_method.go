package vobj

//go:generate stringer -type=HTTPMethod -linecomment -output=http_method__string.go
type HTTPMethod int8

const (
	HTTPMethodUnknown HTTPMethod = iota // 未知
	HTTPMethodGet                       // GET
	HTTPMethodPost                      // POST
	HTTPMethodPut                       // PUT
	HTTPMethodDelete                    // DELETE
	HTTPMethodPatch                     // PATCH
)
