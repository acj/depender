package main

import (
	"strings"
	"testing"
)

func TestSourceFileWithDependenciesImported(t *testing.T) {
	type args struct {
		packagePaths           string
		excludedPathSubstrings []string
	}
	tests := []struct {
		name    string
		args    args
		wantContains    []byte
		wantErr bool
	}{
		{
			name:    "lists dependencies for all packages under current directory",
			args:    args{
				packagePaths:           "./...",
				excludedPathSubstrings: nil,
			},
			wantContains:    []byte(`import (
	_ "golang.org/x/tools/go/packages"
)
`),
			wantErr: false,
		},
		{
			name:    "doesn't list packages when they're excluded",
			args:    args{
				packagePaths:           "./...",
				excludedPathSubstrings: []string{"go/packages"},
			},
			wantContains:    []byte(`import (
)
`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateDummySource(tt.args.packagePaths, tt.args.excludedPathSubstrings)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDummySource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), string(tt.wantContains)) {
				t.Errorf("GenerateDummySource() got = %v, want %v", string(got), string(tt.wantContains))
			}
		})
	}
}