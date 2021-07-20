package sfdc

type Reporting struct {
	Sandbox    bool
	sessionId  string
	serverHost string
	baseUrl    string
}

type GetReportInput struct {
}

type GetReportOutput struct {
}
