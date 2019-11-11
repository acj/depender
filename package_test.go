package main

import (
	"reflect"
	"testing"
)

func Test_dependenciesForPackages(t *testing.T) {
	type args struct {
		paths                string
		excludedPathSegments []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "lists package dependencies for current package",
			args:    args{
				paths:                ".",
				excludedPathSegments: nil,
			},
			want:    []string{"golang.org/x/tools/go/packages"},
			wantErr: false,
		},
		{
			name:    "doesn't list dependencies when they're excluded",
			args:    args{
				paths:                ".",
				excludedPathSegments: []string{"go/packages"},
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dependenciesForImportPaths(tt.args.paths, tt.args.excludedPathSegments)
			if (err != nil) != tt.wantErr {
				t.Errorf("dependenciesForImportPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dependenciesForImportPaths() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNonStandardDependency(t *testing.T) {
	type args struct {
		packagePath          string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correctly identifies package from standard library",
			args: args{
				packagePath:          "io/ioutil",
			},
			want: false,
		},
		{
			name: "correctly identifies third party package",
			args: args{
				packagePath: "github.com/golang/x/tools/go/packages",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNonStandardDependency(tt.args.packagePath); got != tt.want {
				t.Errorf("isNonStandardDependency() = %v, want %v", got, tt.want)
			}
		})
	}
}