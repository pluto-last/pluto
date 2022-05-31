package utils

type WSPack struct {
	TaskID string
	Type   string
	Code   int
	Data   interface{}
	Msg    string
}

func WSPackOK(Type string, data interface{}) (*WSPack, error) {
	return &WSPack{Type: Type, Data: data}, nil
}
