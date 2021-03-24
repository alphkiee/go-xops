package response

// HostRep ...
type HostRep struct {
	HostName  string `json:"host_name"`
	IP        string `json:"ip"`
	OsVersion string `json:"os_version"`
	AuthType  string `json:"auth_type"`
	Creator   string `json:"creator"`
}

//CmdRep ...
type CmdRep struct {
	IP     string      `json:"ip"`
	Cmd    string      `json:"cmd"`
	Stdout interface{} `json:"stdout"`
}
