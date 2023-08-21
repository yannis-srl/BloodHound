package model

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateTokenValue_charset(t *testing.T) {
	for _, tc := range []struct {
		char uint8
		exp  string
	}{
		{0, "0"},
		{9, "9"},
		{10, "a"},
		{35, "z"},
		{36, "A"},
		{61, "Z"},
		// 62 and 63 are ignored since outside range

		// wraps around % 64
		{64, "0"},
		{65, "1"},
		{125, "Z"},
		{128, "0"},
		{253, "Z"},
		// 254 and 255 are ignored since they are outside the range % 64
	} {
		t.Run(fmt.Sprintf("char %d", tc.char), func(t *testing.T) {
			val, err := generateTokenValue(bytes.NewReader(bytes.Repeat([]byte{tc.char}, 64)))
			require.Nil(t, err)
			require.Equal(t, strings.Repeat(tc.exp, 64), val)
		})
	}
}
func TestGenerateTokenValue_sequences(t *testing.T) {
	for _, tc := range []struct {
		rdata []byte
		exp   string
	}{
		{
			[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		},
		{
			[]byte{193, 136, 20, 214, 147, 220, 62, 75, 125, 239, 13, 190, 208, 236, 229, 118, 147, 18, 43, 179, 230, 106, 197, 1, 3, 65, 136, 226, 20, 132, 129, 158, 177, 124, 227, 92, 96, 176, 145, 182, 76, 254, 77, 138, 219, 157, 168, 248, 47, 2, 184, 46, 111, 179, 193, 189, 132, 226, 224, 254, 144, 82, 169, 161, 241, 37, 169, 5, 197, 198, 221, 151, 1, 52, 56, 156, 35, 15, 27, 66, 46, 247, 125, 88, 102, 224, 132, 88, 154, 83, 47, 22, 125, 141, 87, 208, 36, 136, 66, 173, 19, 134, 125, 231, 53, 14, 237, 171, 121, 201, 215, 51, 238, 124, 161, 67, 76, 133, 137, 237, 28, 82, 172, 65, 212, 169, 96, 70},
			"18kmjsbZLdgIBSjiHPCG51318yk41uNYzswMhScdartEUL2UKLP1Z4ywgiFxNBF5",
		},
	} {
		t.Run(fmt.Sprintf("0x%x", tc.rdata), func(t *testing.T) {
			val, err := generateTokenValue(bytes.NewReader(tc.rdata))
			require.Nil(t, err)
			require.Equal(t, tc.exp, val)
		})
	}
}

func TestNewTokenString(t *testing.T) {
	tok, err := NewTokenString("t_tok")
	require.Nil(t, err)
	require.Equal(t, "T_TOK", tok.Prefix)
	require.NotEmpty(t, tok.value)
	require.Len(t, tok.value, 64)
	require.Equal(t, crc32.ChecksumIEEE([]byte(tok.Prefix+tok.value)), tok.cksum)
}

func TestNewTokenString_badprefix(t *testing.T) {
	_, err := NewTokenString("")
	require.Error(t, err)
}

func TestTokenString_String(t *testing.T) {
	for _, tc := range []struct {
		n   string
		tok TokenString
		exp string
	}{
		{
			"zero",
			TokenString{},
			"",
		},
		{
			"missing value",
			TokenString{Prefix: "BLAH", cksum: 12345},
			"",
		},
		{
			"pattern value and cksum",
			TokenString{
				Prefix: "TOK1",
				value:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				cksum:  987654, // encodes to "48VU"
			},
			"TOK1_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0048VU",
		},
		{
			"random value and cksum",
			TokenString{
				Prefix: "ASDF",
				value:  "18kmjsbZLdgIBSjiHPCG51318yk41uNYzswMhScdartEUL2UKLP1Z4ywgiFxNBF5",
				cksum:  2055449580, // encodes to "2f6skA"
			},
			"ASDF_18kmjsbZLdgIBSjiHPCG51318yk41uNYzswMhScdartEUL2UKLP1Z4ywgiFxNBF52f6skA",
		},
		{
			"max cksum value",
			TokenString{
				Prefix: "CKSUM",
				value:  "0000000000000000000000000000000000000000000000000000000000000000",
				cksum:  math.MaxUint32, // 4294967295 encodes to "4GFfc3"
			},
			"CKSUM_00000000000000000000000000000000000000000000000000000000000000004GFfc3",
		},
	} {
		t.Run(tc.n, func(t *testing.T) {
			require.Equal(t, tc.exp, tc.tok.String())
		})
	}
}
