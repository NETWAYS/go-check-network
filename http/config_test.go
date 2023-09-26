package checkhttpconfig

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTLSConfigEmpty(t *testing.T) {
	configTLSConfig := TLSConfig{
		InsecureSkipVerify: true,
	}

	expected := &tls.Config{
		InsecureSkipVerify: configTLSConfig.InsecureSkipVerify,
	}

	actual, err := NewTLSConfig(&configTLSConfig)
	if err != nil {
		t.Error("did not expect error", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

func TestTLSConfig(t *testing.T) {
	testcert := "testdata/selfsigned.cert.pem"
	testkey := "testdata/selfsigned.key.pem"

	configTLSConfig := TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "Test",
		CAFile:             testcert,
		CertFile:           testcert,
		KeyFile:            testkey,
	}

	actual, err := NewTLSConfig(&configTLSConfig)
	if err != nil {
		t.Error("did not expect error", err)
	}

	actualcert, err := tls.LoadX509KeyPair(testcert, testkey)
	if err != nil {
		t.Error("did not expect error", err)
	}

	cert, err := actual.GetClientCertificate(nil)
	if err != nil {
		t.Error("did not expect error", err)
	}
	if !reflect.DeepEqual(cert, &actualcert) {
		t.Error("\nActual: ", cert, "\nExpected: ", actualcert)
	}
}

func TestBearerAuthRoundTripper(t *testing.T) {
	testtoken := "footoken"

	var rt http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer footoken" {
			t.Error("did expect token to match \nActual: ", token, "\nExpected: ", testtoken)
		}
	}))

	defer svr.Close()

	actual := NewAuthorizationCredentialsRoundTripper("Bearer", testtoken, rt)

	request, _ := http.NewRequest("GET", svr.URL, nil)

	_, err := actual.RoundTrip(request)

	if err != nil {
		t.Errorf("unexpected error while executing RoundTrip: %s", err.Error())
	}
}

func TestBasicAuthRoundTripper(t *testing.T) {
	var rt http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Basic dXNlcjE6cGFzc3dkMg==" {
			t.Error("did expect auth to match \nActual: ", auth, "\nExpected: ", "dXNlcjE6cGFzc3dkMg==")
		}
	}))

	defer svr.Close()

	actual := NewBasicAuthRoundTripper("user1", "passwd2", rt)

	request, _ := http.NewRequest("GET", svr.URL, nil)

	_, err := actual.RoundTrip(request)

	if err != nil {
		t.Errorf("unexpected error while executing RoundTrip: %s", err.Error())
	}
}
