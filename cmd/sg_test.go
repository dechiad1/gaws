package cmd

import "testing"

func TestListSGs(t *testing.T) {

}

/*
TODO: set up some tests with sample.out test data
uncomment the DEBUG print in list to view a sample of the printed struct to be mocked & used as input
*/

func TestCalculateCidrRange(t *testing.T) {
	ip := "10.10.10.12"
	cidr := "24"
	expected := "10.10.10.0/24"
	result := calculateCidrRange(ip, cidr)
	if result != expected {
		t.Logf("expected %s, got %s", expected, result)
		t.Fail()
	}
}
