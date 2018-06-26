package formatters

import (
	. "github.com/LindsayBradford/crm/logging/shared"
	"encoding/json"

)

// NameValuePairFormatter formats a LogAttributes array into a string of comma-separated name-value pairs.
// TODO: Supply example encoding.
type JsonFormatter struct {}

func (this *JsonFormatter) Initialise() {}

func (this *JsonFormatter) Format(attributes LogAttributes) string {
	atttributesAsJson, _:= json.Marshal(attributes)
	return string(atttributesAsJson)
}
