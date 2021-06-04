package plugins

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
)

var (
	Plugin Plugins
)

func init() {
	loadPlugins("./plugins.json")
}

type Plugins struct {
	AVG             bool `json:"avg"`
	Zoner           bool `json:"zoner"`
	Fprot           bool `json:"fprot"`
	WindowsDefender bool `json:"windows_defender"`
	Escan           bool `json:"escan"`
	Mcafee          bool `json:"mcafee"`
	Clamav          bool `json:"clamav"`
	Avira           bool `json:"avira"`
	Kaspersky       bool `json:"kaspersky"`
	DrWeb           bool `json:"drweb"`
	Comodo          bool `json:"comodo"`
	Bitdefender     bool `json:"bitdefender"`
	Avast           bool `json:"avast"`
}

func loadPlugins(path string) {
	tmp := new(Plugins)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("error read plugins : " + err.Error())
	}

	if err := json.Unmarshal(file, &tmp); err != nil {
		panic("error unmarshal plugins : " + err.Error())
	}

	Plugin = *tmp
}

func EnableCount() int {

	return len(Enables())
}

func Enables() []string {
	v := reflect.ValueOf(Plugin)
	if v.String() == "<invalid Value>" {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typeOfT := v.Type()
	resp := []string{}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if f.Kind() == reflect.Bool {
			value := f.Bool()
			if value {
				resp = append(resp, typeOfT.Field(i).Name)
			}

		}
	}
	return resp
}
