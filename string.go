package sfdc

import (
	"encoding/json"
	"fmt"
)

// sfdc.String - print details of connection
func (sfdc *Reporting) String() string {
	return fmt.Sprintf("sandbox=%t baseUrl=%s serverHost=%s sessionId=%s", sfdc.Sandbox, sfdc.baseUrl, sfdc.serverHost, sfdc.sessionId)
}

// GetReportOutput.String - print json representation of the report output struct
func (out *GetReportOutput) String() string {
	str, _ := json.Marshal(out)
	return string(str)
}
