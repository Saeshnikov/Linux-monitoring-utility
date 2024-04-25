package tests

import (
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestFindUnusedPackages(t *testing.T) {
	allPackagesExample, usedPackagesExample, unusedPackagesExample := createTestArrays()
	outputMap := allPackagesExample
	type args struct {
		allPackages  map[string]bool
		usedPackages map[string]bool
		dirPath      string
		outputMap    *map[string]bool
	}
	tests := []struct {
		name        string
		args        args
		expectedArr map[string]bool
		wantErr     bool
	}{
		{"Basic test", args{allPackagesExample, usedPackagesExample, "", &outputMap}, unusedPackagesExample, false},
		{"Negative test", args{allPackagesExample, usedPackagesExample, "", &outputMap}, usedPackagesExample, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpmLayer.FindUnusedPackages(tt.args.usedPackages, tt.args.dirPath, tt.args.outputMap)
			eq := Equal(tt.expectedArr, outputMap)
			if tt.wantErr {
				if eq {
					t.Error("Arrays with different data matched\n")
				}
			} else if !tt.wantErr {
				if !eq {
					t.Error("The resulting array did not match the expected one\n")
				}
			}
		})
	}
	os.RemoveAll("./out/")
}

func createTestArrays() (map[string]bool, map[string]bool, map[string]bool) {
	var allPackagesExample = make(map[string]bool)
	var usedPackagesExample = make(map[string]bool)
	var unusedPackagesExample = make(map[string]bool)
	for i := 1; i <= 50; i++ {
		allPackagesExample["package"+strconv.Itoa(i)] = true
		if i%2 == 0 {
			usedPackagesExample["package"+strconv.Itoa(i)] = true
		} else {
			unusedPackagesExample["package"+strconv.Itoa(i)] = true
		}
	}
	return allPackagesExample, usedPackagesExample, unusedPackagesExample
}

func Equal(unusedPackagesExp map[string]bool, unusedPackagesGot map[string]bool) bool {
	return reflect.DeepEqual(unusedPackagesExp, unusedPackagesGot)
}
