package hyprland

import (
	"fmt"
	"strconv"
)

func GetOption(r Reader, name string) *Option {
	if name == "" {
		return nil
	}
	out, ok := hyprctlJSON(r, "getoption", name)
	if !ok {
		return nil
	}
	var option Option
	if !decodeJSON(out, &option) {
		return nil
	}
	return &option
}

func (r Reader) GetOption(name string) *Option {
	return GetOption(r, name)
}

func OptionString(o *Option) string {
	if o == nil {
		return ""
	}
	if o.Custom != "" {
		return o.Custom
	}
	if o.Str != "" {
		return o.Str
	}
	if o.Float != 0 {
		return strconv.FormatFloat(o.Float, 'f', -1, 64)
	}
	if o.Int != 0 {
		return strconv.Itoa(o.Int)
	}
	if o.Set {
		return fmt.Sprintf("%d", o.Int)
	}
	return ""
}
