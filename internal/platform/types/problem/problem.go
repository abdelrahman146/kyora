package problem

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Status     int            `json:"status"`
	Title      string         `json:"title"`
	Detail     string         `json:"detail"`
	Type       string         `json:"type,omitempty"`
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"extensions,omitempty"`
	err        error          `json:"-"`
}

const aboutBlank = "about:blank"

func (p *Problem) With(key string, value any) *Problem {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}
	p.Extensions[key] = value
	return p
}

func (p *Problem) WithError(err error) *Problem {
	p.err = err
	return p
}

func (p *Problem) Error() string {
	return p.Detail
}

func (p *Problem) Unwrap() error {
	return p.err
}

func (p *Problem) Is(target error) bool {
	if pd, ok := target.(*Problem); ok {
		return p.Status == pd.Status && p.Title == pd.Title && p.Detail == pd.Detail && p.Type == pd.Type && p.Instance == pd.Instance
	}
	return false
}

func (p *Problem) As(target any) bool {
	if pd, ok := target.(**Problem); ok {
		*pd = p
		return true
	}
	return false
}

func (p *Problem) ServeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

func (p *Problem) ServeText(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(p.Status)
	_, _ = w.Write([]byte(p.Detail))
}

func (p *Problem) JSON() ([]byte, error) {
	return json.Marshal(p)
}

func New(status int, title, detail, typeURL, instance string) *Problem {
	p := &Problem{
		Status:   status,
		Title:    title,
		Detail:   detail,
		Type:     typeURL,
		Instance: instance,
	}
	if p.Type == "" {
		p.Type = aboutBlank
	}
	return p
}

func InternalError() *Problem {
	return &Problem{
		Status: 500, Title: "Internal Server Error", Detail: "An internal server error occurred", Type: aboutBlank,
	}
}

func NotFound(detail string) *Problem {
	return &Problem{
		Status: 404, Title: "Not Found", Detail: detail, Type: aboutBlank,
	}
}

func BadRequest(detail string) *Problem {
	return &Problem{
		Status: 400, Title: "Bad Request", Detail: detail, Type: aboutBlank,
	}
}

func Unauthorized(detail string) *Problem {
	return &Problem{
		Status: 401, Title: "Unauthorized", Detail: detail, Type: aboutBlank,
	}
}

func Forbidden(detail string) *Problem {
	return &Problem{
		Status: 403, Title: "Forbidden", Detail: detail, Type: aboutBlank,
	}
}

func UnprocessableEntity(detail string) *Problem {
	return &Problem{
		Status: 422, Title: "Unprocessable Entity", Detail: detail, Type: aboutBlank,
	}
}

func Conflict(detail string) *Problem {
	return &Problem{
		Status: 409, Title: "Conflict", Detail: detail, Type: aboutBlank,
	}
}
