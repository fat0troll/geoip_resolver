package requesterinterface

// RequesterInterface is an interface which represents functions of Requester package
type RequesterInterface interface {
	Init()

	ProcessRequest(IPAddress string) map[string]string
}
