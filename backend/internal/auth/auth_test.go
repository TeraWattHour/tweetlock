package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGuard(t *testing.T) {
	os.Setenv("ACCESS_SECRET", "ayestrongsecret")

	testCases := []struct {
		Input    string
		Expected string
	}{{
		Input:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.-_vDglxXWFGPd4ku4hr4vrk6t4p-09PRPVBNlPebk-g",
		Expected: "1",
	}, {
		Input:    "",
		Expected: "",
	}, {
		Input:    "dsadasdsadasdsa",
		Expected: "",
	}, {
		Input:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.GRZNBFtV1sonVZdmWdP2xPyf7_uzwVLmh6tmCNxEVh0",
		Expected: "",
	}}

	for ix, test := range testCases {
		userID, _ := Guard(events.APIGatewayProxyRequest{
			Headers: map[string]string{
				"Cookie": fmt.Sprintf("x-access=%s", test.Input),
			},
		})
		if userID != test.Expected {
			t.Errorf("case %d failed, expected `%v`, got `%v`\n", ix+1, test.Expected, userID)
		}
	}
}
