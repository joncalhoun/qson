package query

import "testing"

func TestNested(t *testing.T) {
	query := "bar%5Bone%5D%5Btwo%5D=2&bar[one][red]=112"
	expected := `{"bar":{"one":{"red":112,"two":2}}}`
	actual, err := JSON(query)
	if err != nil {
		t.Error(err)
	}
	actualStr := string(actual)
	if expected != actualStr {
		t.Errorf("Expected %s, received %s", expected, actualStr)
	}
}

func TestPlain(t *testing.T) {
	query := "cat=1&dog=2"
	expected := `{"cat":1,"dog":2}`
	actual, err := JSON(query)
	if err != nil {
		t.Error(err)
	}
	actualStr := string(actual)
	if expected != actualStr {
		t.Errorf("Expected %s, received %s", expected, actualStr)
	}
}

func TestSlice(t *testing.T) {
	query := "cat[]=1"
	expected := `{"cat":[1]}`
	actual, err := JSON(query)
	if err != nil {
		t.Error(err)
	}
	actualStr := string(actual)
	if expected != actualStr {
		t.Errorf("Expected %s, received %s", expected, actualStr)
	}
}
