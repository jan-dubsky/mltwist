package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKey_ValidateProps(t *testing.T) {
	tests := []struct {
		name   string
		k      Key
		s      keyScope
		p      keyPermission
		hasErr bool
	}{{
		name: "valid_ip_key",
		k:    "#r:w:ip",
		s:    'r',
		p:    'w',
	}, {
		name:   "missing_name",
		k:      "#r:w:",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_key_scope",
		k:      "#c:w:ip",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_key_permissions",
		k:      "#r:f:ip",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_key_scope",
		k:      "#c:w:ip",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_wanted_scope",
		k:      "#r:w:ip",
		s:      'i',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_wanted_permissions",
		k:      "#r:w:ip",
		s:      'r',
		p:      'g',
		hasErr: true,
	}, {
		name:   "missing_first_colon",
		k:      "#row:",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "missing_second_colon",
		k:      "#r:wow",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "not_allowed_scope",
		k:      "#r:w:ip",
		s:      'm',
		p:      'w',
		hasErr: true,
	}, {
		name:   "not_allowed_permissions",
		k:      "#r:w:ip",
		s:      'r',
		p:      'r',
		hasErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k.validateProps(tt.s, tt.p)
			if r := require.New(t); tt.hasErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestKey_Validate(t *testing.T) {
	tests := []struct {
		name   string
		k      Key
		s      keyScope
		p      keyPermission
		hasErr bool
	}{{
		name: "non_reserved",
		k:    "foo",
		s:    'r',
		p:    'w',
	}, {
		name: "instruction_pointer",
		k:    "#r:w:ip",
		s:    'r',
		p:    'w',
	}, {
		name:   "empty",
		k:      "",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "empty",
		k:      "",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "unknown_reserved",
		k:      "#r:w:foobar",
		s:      'r',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_scope",
		k:      "#r:w:ip",
		s:      'm',
		p:      'w',
		hasErr: true,
	}, {
		name:   "invalid_permissions",
		k:      "#r:w:ip",
		s:      'r',
		p:      'r',
		hasErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.k.validate(tt.s, tt.p)
			if r := require.New(t); tt.hasErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
