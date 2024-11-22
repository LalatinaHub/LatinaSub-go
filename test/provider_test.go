package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/LalatinaHub/LatinaSub-go/provider"
	"github.com/stretchr/testify/require"
)

var (
	vmessLink = "vmess://eyJhZGQiOiJ6b29tLnVzIiwiYWlkIjowLCJob3N0IjoiZGUtbmV3MDEuZGFsdXF1YW4udG9wIiwiaWQiOiI2ZTc5YTQwMy02OWVhLTQ1ZjQtYjk4OS1mMTAxNmEyOGM3NGEiLCJuZXQiOiJ3cyIsInBhdGgiOiIvIiwicG9ydCI6ODA4MCwicHMiOiIxOTYg8J+HuPCfh6wgT1JBQ0xFIENPUlBPUkFUSU9OIFdTIENETiBOVExTIiwidGxzIjoiIiwic2VjdXJpdHkiOiJhdXRvIiwic2tpcC1jZXJ0LXZlcmlmeSI6dHJ1ZSwic25pIjoiZGUtbmV3MDEuZGFsdXF1YW4udG9wIn0="
	vlessLink = "vless://b59a150a-7096-49fd-8081-368c518e83d0@zoom.us:443?allowInsecure=true&security=tls&serviceName=v2rayfree1_v2rayfree1_v2rayfree1_v2rayfree1_v2rayfree1_v2rayfree1_v2rayfree1_v2rayfree1&sni=&type=grpc#777%20%F0%9F%87%B9%F0%9F%87%B7%20STARK%20INDUSTRIES%20SOLUTIONS%20LTD%20GRPC%20CDN%20TLS"
)

func TestProvider(t *testing.T) {
	t.Run("vmess_parser", func(t *testing.T) {
		require.NoError(t, testParser(vmessLink))
	})

	t.Run("vless_parser", func(t *testing.T) {
		require.NoError(t, testParser(vlessLink))
	})
}

func testParser(content string) error {
	result, err := provider.Parse(content)
	if err != nil {
		return err
	}

	if j, err := json.MarshalIndent(result, "", "  "); err == nil {
		fmt.Println(string(j[:]))
	} else {
		return err
	}

	return nil
}
