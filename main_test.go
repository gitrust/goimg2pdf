package main

import "testing"

func TestValidExtensions(t *testing.T) {
	testCases := []string{
		".jpg", ".jpeg", ".png", ".JPEG", ".PNG",
	}

	for _, tc := range testCases {
		result := isValidExtension(tc)
		if !result {
			t.Errorf("isValidExtension(%s) = %v; expected %v", tc, result, true)
		}
	}
}

func TestInvalidExtensions(t *testing.T) {
	testCases := []string{
		".txt", ".tiff", ".gif",
	}

	for _, tc := range testCases {
		result := isValidExtension(tc)
		if result {
			t.Errorf("isValidExtension(%s) = %v; expected %v", tc, result, false)
		}
	}
}
