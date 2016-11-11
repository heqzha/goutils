package encoding

import(
	"encoding/json"
)

func ObjToJSON(o interface{})(string, error){
	js, err := json.Marshal(o)
	if err != nil{
		return "", err
	}
	return string(js), nil
}

func JSONToObj(js string, o interface{})(error){
	err := json.Unmarshal([]byte(js), o)
	if err != nil{
		return err
	}
	return nil
}
