package tests

import (
	rpmLayer "linux-monitoring-utility/internal/rpmLayer"
	"reflect"
	"strconv"
	"testing"
)

func TestFindUnusedPackages(t *testing.T) {
	var allPackagesExample = make(map[string]bool)
	for i := 1; i < 50; i++ {
		allPackagesExample["package"+strconv.Itoa(i)] = true
	}
	var usedPackagesExample = make(map[string]bool)
	for i := 1; i < 50; i++ {
		if i%2 == 0 {
			usedPackagesExample["package"+strconv.Itoa(i)] = true
		}
	}

	var outputMap = allPackagesExample
	rpmLayer.FindUnusedPackages(allPackagesExample, usedPackagesExample, "", &outputMap)
	for packageName := range usedPackagesExample {
		if _, ok := allPackagesExample[packageName]; ok {
			delete(allPackagesExample, packageName)
		}
	}

	if !reflect.DeepEqual(allPackagesExample, outputMap) {
		t.Error("The resulting array did not match the expected one")
	}
}
