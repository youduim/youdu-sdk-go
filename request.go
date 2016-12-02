package ydapp

type ApiRequest struct {
	*Json
}

func NewRequest() *ApiRequest {
	return &ApiRequest{
		Json: NewEmptyJson(),
	}
}
