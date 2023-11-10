package provider

import (
	"strings"

	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func Parse(content string) ([]option.Outbound, error) {
	var outbounds []option.Outbound
	var err error
	if strings.Contains(content, "\"outbounds\"") {
		var options option.Options
		err = options.UnmarshalJSON([]byte(content))
		if err != nil {
			return nil, E.Cause(err, "decode config at ")
		}
		outbounds = options.Outbounds
		return outbounds, nil
	} else if strings.Contains(content, "proxies") {
		outbounds, err = newClashParser(content)
		if err != nil {
			return nil, err
		}
		return outbounds, nil
	}
	outbounds, err = newNativeURIParser(content)
	if err != nil {
		return nil, err
	}
	return outbounds, err
}
