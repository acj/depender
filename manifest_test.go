package main

import (
	"strings"
	"testing"
)

func TestManifest_Generate(t *testing.T) {
	type fields struct {
		PackagePaths         string
		ExcludedPathSegments []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantContains    []byte
		wantErr bool
	}{
		{
			name:    "lists dependencies for all packages under current directory",
			fields:  fields{
				PackagePaths:         "./...",
				ExcludedPathSegments: nil,
			},
			wantContains:    []byte(`import (
	_ "golang.org/x/tools/go/packages"
)
`),
			wantErr: false,
		},
		{
			name:         "doesn't list packages when they're excluded",
			fields:       fields{
				PackagePaths:         "./...",
				ExcludedPathSegments: []string{"go/packages"},
			},
			wantContains: []byte(`import (
)
`),
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manifest{
				PackagePaths:         tt.fields.PackagePaths,
				ExcludedPathSegments: tt.fields.ExcludedPathSegments,
			}
			got, err := m.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), string(tt.wantContains)) {
				t.Errorf("Generate() got = '%s', want '%s'", got, tt.wantContains)
			}
		})
	}
}