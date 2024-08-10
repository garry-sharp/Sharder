package crypt

import (
	"reflect"
	"testing"
)

var testingTreeStrings []string = []string{"cat", "mat", "sit", "cot"}

func TestBuildWordListTree(t *testing.T) {

	treeBuilt := buildWordListTree(testingTreeStrings)

	treeExpected := tree{
		children: []tree{
			{
				value: "c",
				children: []tree{
					{
						value: "a",
						children: []tree{
							{
								value: "t",
								children: []tree{
									{
										value: "cat",
										leaf:  true,
									},
								},
							},
						},
					},
					{
						value: "o",
						children: []tree{
							{
								value: "t",
								children: []tree{
									{
										value: "cot",
										leaf:  true,
									},
								},
							},
						},
					},
				},
			},
			{
				value: "m",
				children: []tree{
					{
						value: "a",
						children: []tree{
							{
								value: "t",
								children: []tree{
									{
										value: "mat",
										leaf:  true,
									},
								},
							},
						},
					},
				},
			},
			{
				value: "s",
				children: []tree{
					{
						value: "i",
						children: []tree{
							{
								value: "t",
								children: []tree{
									{
										value: "sit",
										leaf:  true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !treeExpected.Equal(treeBuilt) {
		t.Errorf("Expected :\n%s\nReceived:\n%s\n", treeBuilt.String(), treeExpected.String())
	}

}

func TestGetOptions(t *testing.T) {
	treeBuilt := buildWordListTree(testingTreeStrings)

	options1 := getOptions("c", treeBuilt)
	options2 := getOptions("co", treeBuilt)
	options3 := getOptions("", treeBuilt)
	options4 := getOptions("not", treeBuilt)

	expectedOptions1 := []string{"cat", "cot"}
	expectedOptions2 := []string{"cot"}
	expectedOptions3 := []string{"cat", "cot", "mat", "sit"}
	expectedOptions4 := []string{}

	if !reflect.DeepEqual(options1, expectedOptions1) {
		t.Errorf("Expected: %s, but got: %s", expectedOptions1, options1)
	}
	if !reflect.DeepEqual(options2, expectedOptions2) {
		t.Errorf("Expected: %s, but got: %s", expectedOptions2, options2)
	}
	if !reflect.DeepEqual(options3, expectedOptions3) {
		t.Errorf("Expected: %s, but got: %s", expectedOptions3, options3)
	}
	if !reflect.DeepEqual(options4, expectedOptions4) {
		t.Errorf("Expected: %s, but got: %s", expectedOptions4, options4)
	}
}
