package params

type CommonReq struct {
	Head  string `json:"head"`
	Value string `json:"value"`
}

type LoginValue struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Label    string `json:"label"`
	Facility string `json:"facility"`
}

type LoginResp struct {
	Verify bool     `json:"verify"`
	Uinfo  UserInfo `json:"u_info"`
	Token  string   `json:"token"`
}

type UserInfo struct {
	Account  string `json:"account"`
	UserName string `json:"user_name"`
	PhoneNo  string `json:"phone_no"`
}

type RobotListResp struct {
	LoginResp
	RobotList []Robot `json:"robot_list"`
}

type Robot struct {
	SoftID    string `json:"soft_id"`
	RobotName string `json:"robot_name"`
	JobState  string `json:"job_state"`
	Online    int    `json:"online"`
	IpAddress string `json:"ip_address"`
	Port      string `json:"port"`
}

type StopTaskReq struct {
	Head  string `json:"head"`
	Value string `json:"value"`
	Sid   int    `json:"sid"`
}

type StopTaskRsp struct {
	Code int `json:"code"`
}
