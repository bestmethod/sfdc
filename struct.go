package sfdc

type Reporting struct {
	Sandbox    bool
	User       string
	Password   string
	SecToken   string
	sessionId  string
	serverHost string
	baseUrl    string
}

type GetReportInput struct {
}

type GetReportOutput struct {
}
