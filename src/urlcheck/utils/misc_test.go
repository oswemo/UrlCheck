package utils

import "testing"

func TestURLDecode(t *testing.T) {
	output, err := URLDecode("https%3A%2F%2Fwww.google.ca%2F%3Fgfe_rd%3Dcr%26ei%3Dt8AjWZjGM4PktgbqyJagDg%23q%3Durl%2Bdecode")
	if err != nil {
		t.Error("Error decoding URL")
	}

	if output != "https://www.google.ca/?gfe_rd=cr&ei=t8AjWZjGM4PktgbqyJagDg#q=url+decode" {
		t.Error("Decoded URL does not match")
	}
}

func TestValidateUrl(t *testing.T) {
	testCases := []struct {
		Url   string
		Valid bool
	}{
		{
			Url:   "www.google.com",
			Valid: false,
		},
		{
			Url:   "www.google_com:80",
			Valid: false,
		},
		{
			Url:   "www.google.com:80000",
			Valid: false,
		},
		{
			Url:   "www.google.com:0",
			Valid: false,
		},
		{
			Url:   "google.com:80",
			Valid: true,
		},
		{
			Url:   "www.google.com:80",
			Valid: true,
		},
	}

	for _, c := range testCases {
		err := ValidateUrl(c.Url)

		if err != nil {
			if c.Valid {
				t.Errorf("Validation of %s failed.  Expected it to be valid.", c.Url)
			}
		} else {
			if !c.Valid {
				t.Errorf("Validation of %s failed.  Expected it to be invalid.", c.Url)
			}
		}
	}
}
