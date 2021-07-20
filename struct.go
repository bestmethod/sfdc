package sfdc

// Reporting - basic object to init, before calling .Login() on it
type Reporting struct {
	Sandbox    bool   // true == sfdc test
	User       string // username to login with
	Password   string // password to login with
	SecToken   string // security token to use to login
	sessionId  string // this is the token/session key that we use to communicate with sfdc post login
	serverHost string // the serverHost parameter returns the hostname to use for further sfdc communication post login
	baseUrl    string // the baseUrl of the sfdc host for future communications post login
}

// GetReportInput - fill with details required to GetReport()
type GetReportInput struct {
}

// GetReportOutput - the return data from GetReport()
type GetReportOutput struct {
}
