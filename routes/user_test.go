package routes

import "testing"

func TestCreateUserValidator(t *testing.T) {
	var testFailingCases = []User{
		User{FirstName: ""},
		User{FirstName: "test", LastName: ""},
		User{FirstName: "test", LastName: "test", Email: ""},
		User{FirstName: "test", LastName: "test", Email: "test@test.gov", Role: ""},
		User{FirstName: "test", LastName: "test", Email: "badguy@yourmomshouse.com;test@test.gov", Role: "nonmember"},
		User{FirstName: "test", LastName: "test", Email: "badguy@yourmomshouse.com,test@test.gov", Role: "nonmember"},
		User{FirstName: "test", LastName: "test", Email: "test@test.com", Role: "nonmember"},
		User{FirstName: "test", LastName: "test", Email: "test@test.gov", Role: "admin"},
		User{FirstName: "test", LastName: "test", Email: "notanemail.gov", Role: "nonmember"},
	}
	for _, testCase := range testFailingCases {
		if err := validateExternalUserRequest(testCase); err == nil {
			t.Fatalf("Expected failed validation, test case: %v", testCase)
		}
	}

	var testSuccessCases = []User{
		User{FirstName: "test", LastName: "test", Email: "CapitalY@test.gov", Role: "member"},
		User{FirstName: "test", LastName: "test-hyphen", Email: "yo@test.gov", Role: "member"},
		User{FirstName: "test", LastName: "test", Email: "CapitalY@test.gov", Role: "nonmember"},
	}
	for _, testCase := range testSuccessCases {
		if err := validateExternalUserRequest(testCase); err != nil {
			t.Fatalf("Expected validation, test case: %v", testCase)
		}
	}
}
