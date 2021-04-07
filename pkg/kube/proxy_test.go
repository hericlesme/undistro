package kube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

func newProxy(target *url.URL) http.Handler {
	p := proxy.NewUpgradeAwareHandler(target, http.DefaultTransport, false, false, &responder{})
	p.UseRequestLocation = true
	return p
}

func TestAPIRequests(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s %s %s", r.Method, r.RequestURI, string(b))
	}))
	defer ts.Close()

	// httptest.NewServer should always generate a valid URL.
	target, _ := url.Parse(ts.URL)
	target.Path = "/"
	proxy := newProxy(target)

	tests := []struct{ name, method, body string }{
		{"test1", "GET", ""},
		{"test2", "DELETE", ""},
		{"test3", "POST", "test payload"},
		{"test4", "PUT", "test payload"},
	}

	const path = "/api/test?fields=ID%3Dfoo&labels=key%3Dvalue"
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, path, strings.NewReader(tt.body))
			if err != nil {
				t.Errorf("error creating request: %v", err)
				return
			}
			w := httptest.NewRecorder()
			proxy.ServeHTTP(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("%d: proxy.ServeHTTP w.Code = %d; want %d", i, w.Code, http.StatusOK)
			}
			want := strings.Join([]string{tt.method, path, tt.body}, " ")
			if w.Body.String() != want {
				t.Errorf("%d: response body = %q; want %q", i, w.Body.String(), want)
			}
		})
	}
}

func TestPathHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	}))
	defer ts.Close()

	table := []struct {
		name       string
		prefix     string
		reqPath    string
		expectPath string
	}{
		{"test1", "/api/", "/metrics", "404 page not found\n"},
		{"test4", "/", "/metrics", "/metrics"},
		{"test5", "/", "/api/v1/pods/", "/api/v1/pods/"},
		{"test6", "/custom/", "/metrics", "404 page not found\n"},
		{"test7", "/custom/", "/api/metrics", "404 page not found\n"},
		{"test8", "/custom/", "/api/v1/pods/", "404 page not found\n"},
		{"test9", "/custom/", "/custom/api/metrics", "/api/metrics"},
		{"test10", "/custom/", "/custom/api/v1/pods/", "/api/v1/pods/"},
	}

	cc := &rest.Config{
		Host: ts.URL,
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProxyHandler(tt.prefix, cc, time.Duration(0))
			if err != nil {
				t.Fatalf("%#v: %v", tt, err)
			}
			pts := httptest.NewServer(p)
			defer pts.Close()

			r, err := http.Get(pts.URL + tt.reqPath)
			if err != nil {
				t.Fatalf("%#v: %v", tt, err)
			}
			body, err := ioutil.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				t.Fatalf("%#v: %v", tt, err)
			}
			if e, a := tt.expectPath, string(body); e != a {
				t.Errorf("%#v: Wanted %q, got %q", tt, e, a)
			}
		})
	}
}
