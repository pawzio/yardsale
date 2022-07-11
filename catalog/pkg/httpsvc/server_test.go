package httpsvc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	type testCase struct {
		givenOpts            []ServerOption
		expErr               bool
		expReadTimeout       time.Duration
		expReadHeaderTimeout time.Duration
		expWriteTimeout      time.Duration
		expIdleTimeout       time.Duration
		expShutdownGrace     time.Duration
		expPort              string
	}
	tcs := map[string]testCase{
		"success": {
			expReadTimeout:       5 * time.Second,
			expReadHeaderTimeout: 2 * time.Second,
			expWriteTimeout:      30 * time.Second,
			expIdleTimeout:       120 * time.Second,
			expShutdownGrace:     10 * time.Second,
			expPort:              "3000",
		},
		"success with all options": {
			givenOpts:            []ServerOption{WithServerPort("9000")},
			expReadTimeout:       5 * time.Second,
			expReadHeaderTimeout: 2 * time.Second,
			expWriteTimeout:      30 * time.Second,
			expIdleTimeout:       120 * time.Second,
			expShutdownGrace:     10 * time.Second,
			expPort:              "9000",
		},
		"invalid port": {
			givenOpts: []ServerOption{WithServerPort("abcd")},
			expErr:    true,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given && When:
			s, err := NewServer(nil, nil, tc.givenOpts...)

			// Then:
			if tc.expErr {
				require.Nil(t, s)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, s)
			require.Equal(t, fmt.Sprintf(":%s", tc.expPort), s.srv.Addr)
			require.Equal(t, tc.expReadTimeout, s.srv.ReadTimeout)
			require.Equal(t, tc.expReadHeaderTimeout, s.srv.ReadHeaderTimeout)
			require.Equal(t, tc.expWriteTimeout, s.srv.WriteTimeout)
			require.Equal(t, tc.expIdleTimeout, s.srv.IdleTimeout)
			require.Equal(t, tc.expShutdownGrace, s.shutdownGrace)
		})
	}
}
