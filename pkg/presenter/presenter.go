package presenter

type (
	Response struct {
		ClientID  string `json:"client_id,omitempty"`
		UserAgent string `json:"user_agent"`
		Version   int    `json:"version"`
	}

	ResponseV2 struct {
		Response
		Message string `json:"message"`
	}
)

func NewResponse(clientID, userAgent string) Response {
	return Response{
		ClientID:  clientID,
		UserAgent: userAgent,
		Version:   1,
	}
}

func NewResponseV2(response Response) ResponseV2 {
	return ResponseV2{
		Response: Response{
			ClientID:  response.ClientID,
			UserAgent: response.UserAgent,
			Version:   2,
		},
		Message: "looks like the feature flag is on",
	}
}
