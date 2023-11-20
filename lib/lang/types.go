package lang

type ReqXBodyType string

const (
	JSON      ReqXBodyType = "json"
	FORM      ReqXBodyType = "form"
	MULTIPART ReqXBodyType = "multipart"
	GRPC      ReqXBodyType = "grpc"
)

type ReqXBody interface {
	GetType() ReqXBodyType
}

type ReqXRequestUrl struct {
	Hostname string            `json:"hostname"`
	Path     string            `json:"path"`
	Protocol string            `json:"protocol"`
	Query    map[string]string `json:"query"`
}

type ReqXRequest struct {
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Url     ReqXRequestUrl    `json:"url"`
	Type    ReqXBodyType      `json:"type"`
	Body    ReqXBody          `json:"body"`
}

type ReqXAssertion struct {
	Operator string `json:"operator"`
	Field    string `json:"field"`
	Operand  string `json:"operand"`
	JsonPath string `json:"jsonPath"`
}

type ReqX struct {
	Request    ReqXRequest     `json:"request"`
	Assertions []ReqXAssertion `json:"assertions"`
}
