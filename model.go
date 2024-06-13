package mlog

type ClientLogModel struct {
	SourceIP       string      `json:"sourceIP,omitempty"`
	RequestHeader  interface{} `json:"requestHeader,omitempty"`
	Request        string      `json:"request,omitempty"`
	ResponseHeader interface{} `json:"responseHeader,omitempty"`
	Response       string      `json:"response,omitempty"`
}

type LegacyLogModel struct {
	StepName             string      `json:"stepname,omitempty"`
	Endpoint             string      `json:"endpoint,omitempty"`
	Method               string      `json:"method,omitempty"`
	LegacyRequestHeader  interface{} `json:"legacyRequestHeader,omitempty"`
	LegacyRequest        string      `json:"legacyRequest,omitempty"`
	LegacyResponseHeader interface{} `json:"legacyResponseHeader,omitempty"`
	LegacyResponse       string      `json:"legacyResponse,omitempty"`
}

type MessageQueueLogModel struct {
	MessageQueueHeader    interface{} `json:"messageQueueHeader,omitempty"`
	MessageQueueTopic     string      `json:"messageQueueTopic,omitempty"`
	MessageQueueValue     interface{} `json:"messageQueueValue,omitempty"`
	MessageQueuePartition interface{} `json:"messageQueuePartition,omitempty"`
	MessageQueueOffset    interface{} `json:"messageQueueOffset,omitempty"`
}
