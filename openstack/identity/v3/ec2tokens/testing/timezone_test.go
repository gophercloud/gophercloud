package testing

import (
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSignatureV4Timezone(t *testing.T) {
	locations := []*time.Location{
		time.UTC,
		time.FixedZone("UTC+1", 3600),
		time.FixedZone("UTC+2", 2*3600),
		time.FixedZone("UTC-7", -7*3600),
		time.FixedZone("UTC+5:45", 5*3600+45*60),
		time.FixedZone("", 90*60),
		time.Local,
	}
	for _, instant := range []time.Time{
		time.Date(2032, 12, 31, 23, 22, 22, 0, time.UTC),
		time.Date(2033, 1, 1, 0, 22, 22, 0, time.UTC),
	} {
		t.Run(instant.Format(time.RFC3339), func(t *testing.T) {
			bodyHash := "body-hash"
			opts := auth.EC2TokenOpts{
				Access: "access", Secret: "secret", Region: "region", Service: "ec2",
				Verb: "GET", Path: "/", BodyHash: &bodyHash, Timestamp: &instant,
				Headers: map[string]string{"Host": "localhost", "X-Amz-SignedHeaders": "host"},
			}
			expectedKey := auth.EC2CredentialsBuildSignatureKeyV4(opts.Secret, opts.Region, opts.Service, instant)
			expectedString := auth.EC2CredentialsBuildStringToSignV4(opts, "host", bodyHash, instant)
			expectedSignature := auth.EC2CredentialsBuildSignatureV4(expectedKey, expectedString)
			expectedHeader := auth.EC2CredentialsBuildAuthorizationHeaderV4(opts, "host", expectedSignature, instant)
			expectedBody, err := opts.ToAuthBody()
			th.AssertNoErr(t, err)
			for _, loc := range locations {
				t.Run(loc.String(), func(t *testing.T) {
					timestamp := instant.In(loc)
					original := timestamp
					key := auth.EC2CredentialsBuildSignatureKeyV4(opts.Secret, opts.Region, opts.Service, timestamp)
					th.CheckDeepEquals(t, expectedKey, key)
					stringToSign := auth.EC2CredentialsBuildStringToSignV4(opts, "host", bodyHash, timestamp)
					th.CheckDeepEquals(t, expectedString, stringToSign)
					lines := strings.Split(string(stringToSign), "\n")
					th.AssertEquals(t, 4, len(lines))
					th.CheckEquals(t, instant.Format("20060102T150405Z"), lines[1])
					th.CheckEquals(t, instant.Format("20060102")+"/region/ec2/aws4_request", lines[2])
					signature := auth.EC2CredentialsBuildSignatureV4(key, stringToSign)
					th.CheckEquals(t, expectedSignature, signature)
					th.CheckEquals(t, expectedHeader, auth.EC2CredentialsBuildAuthorizationHeaderV4(opts, "host", expectedSignature, timestamp))
					localOpts := opts
					localOpts.Timestamp = &timestamp
					body, err := localOpts.ToAuthBody()
					th.AssertNoErr(t, err)
					th.CheckDeepEquals(t, expectedBody, body)
					headers := body["credentials"]["headers"].(map[string]string)
					th.CheckEquals(t, instant.Format("20060102T150405Z"), headers["X-Amz-Date"])
					th.CheckDeepEquals(t, original, timestamp)
					th.CheckDeepEquals(t, map[string]string{"Host": "localhost", "X-Amz-SignedHeaders": "host"}, opts.Headers)
				})
			}
		})
	}
}
