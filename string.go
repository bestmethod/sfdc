package sfdc

import "fmt"

func (sfdc *Reporting) String() string {
	return fmt.Sprintf("sandbox=%t baseUrl=%s serverHost=%s sessionId=%s", sfdc.Sandbox, sfdc.baseUrl, sfdc.serverHost, sfdc.sessionId)
}
