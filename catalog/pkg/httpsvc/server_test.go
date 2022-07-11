package httpsvc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	type testCase struct {
		expErr bool
	}
	tcs := map[string]testCase{
		"success": {},
		"invalid port": {
			expErr: true,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given && When:
			s, err := NewServer(nil)

			// Then:
			if tc.expErr {
				require.Nil(t, s)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, s)
			require.Equal(t, ":3000", s.srv.Addr)
			require.Equal(t, 5*time.Second, s.srv.ReadTimeout)
			require.Equal(t, 2*time.Second, s.srv.ReadHeaderTimeout)
			require.Equal(t, 30*time.Second, s.srv.WriteTimeout)
			require.Equal(t, 120*time.Second, s.srv.IdleTimeout)
		})
	}
}
