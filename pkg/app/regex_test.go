package app

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRegexCache(t *testing.T) {
	tests := []struct {
		mapping Mapping
		wantLen int
		wantErr bool
	}{
		{
			mapping: Mapping{
				Request: RequestMapping{
					Path: PathMapping{Pattern: []string{`/[A-z0-9]+/`}},
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			mapping: Mapping{
				Request: RequestMapping{
					Path:    PathMapping{Pattern: []string{`[A-z0-9]+`}},
					Headers: map[string]HeaderMapping{"accept": {Pattern: []string{".*"}}, "x-id": {Pattern: []string{`\d*`}}, "x-debug": {Pattern: []string{".*"}}},
					Body:    BodyMapping{Pattern: []string{`\d{3}\.\d{3}\.\d{3}-\d{2}`}},
				},
			},
			wantLen: 4,
			wantErr: false,
		},
		{
			mapping: Mapping{
				Request: RequestMapping{
					Path: PathMapping{Pattern: []string{`([A-z0-9]+`}},
				},
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			mapping: Mapping{
				Request: RequestMapping{
					Headers: map[string]HeaderMapping{"accept": {Pattern: []string{"((.*json)"}}},
				},
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			mapping: Mapping{
				Request: RequestMapping{
					Body: BodyMapping{Pattern: []string{`\d{)}*`}},
				},
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		rc := NewRegexCache()
		err := rc.AddFromMapping(tt.mapping)
		if !assert.IsEqual(err != nil, tt.wantErr) {
			t.Log("error parsing regex pattern: ", err)
			t.Fail()
		}

		assert.Equal(t, len(rc.cache), tt.wantLen)

		if !tt.wantErr {

			for _, p := range tt.mapping.Request.Path.Pattern {
				_, ok := rc.cache[p]
				assert.Equal(t, true, ok)
			}
			for _, p := range tt.mapping.Request.Body.Pattern {
				_, ok := rc.cache[p]
				assert.Equal(t, true, ok)
			}
			for _, v := range tt.mapping.Request.Headers {
				for _, p := range v.Pattern {
					_, ok := rc.cache[p]
					assert.Equal(t, true, ok)
				}
			}
		}
	}

}
