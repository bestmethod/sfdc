package sfdc

import (
	"encoding/json"
	"fmt"
)

func (sfdc *Reporting) String() string {
	return fmt.Sprintf("sandbox=%t baseUrl=%s serverHost=%s sessionId=%s", sfdc.Sandbox, sfdc.baseUrl, sfdc.serverHost, sfdc.sessionId)
}

func (out *GetReportOutput) String() string {
	str, _ := json.MarshalIndent(out, "", "    ")
	return string(str)
}
