package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"

	hws_cloud "github.com/ServiceComb/auth/third_party/forked/datastream/aws"
)

// Headers for ak/sk auth
const (
	HeaderServiceAk      = "X-Service-AK"
	HeaderServiceShaAKSK = "X-Service-ShaAKSK"
	HeaderServiceProject = "X-Service-Project"
)

type AuthFunc func(*http.Request) error

var globalAuthFunc AuthFunc

// AddAuthInfo adds auth info into request
func AddAuthInfo(r *http.Request) error {
	return globalAuthFunc(r)
}

// SetAuthFunc sets a custom auth func
func SetAuthFunc(f func(*http.Request) error) {
	globalAuthFunc = f
}

// UseAKSKAuth sets and initializes the ak/sk auth func
func UseAKSKAuth(ak, sk, project string) error {
	s := &hws_cloud.Signer{
		AccessKey: ak,
		SecretKey: sk,
		Service:   "",
		Region:    "",
	}

	shaakskAuthFunc, err := newShaAKSKAuthFunc(ak, sk, project)
	if err != nil {
		return err
	}

	globalAuthFunc = func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		if err := shaakskAuthFunc(r); err != nil {
			return err
		}
		return s.Sign(r)
	}
	return nil
}

func newShaAKSKAuthFunc(ak, sk, project string) (AuthFunc, error) {
	shaAKSK, err := genShaAKSK(sk, ak)
	if err != nil {
		return nil, err
	}

	f := func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set(HeaderServiceAk, ak)
		r.Header.Set(HeaderServiceShaAKSK, shaAKSK)
		r.Header.Set(HeaderServiceProject, project)
		return nil
	}
	return f, nil
}

// UseShaAKSKAuth sets and initializes the ak/sk auth func
func UseShaAKSKAuth(ak, sk, project string) error {
	shaakskAuthFunc, err := newShaAKSKAuthFunc(ak, sk, project)
	if err != nil {
		return err
	}

	globalAuthFunc = shaakskAuthFunc
	return nil
}

func genShaAKSK(key string, data string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	if _, err := h.Write([]byte(data)); err != nil {
		return "", err
	}
	b := h.Sum(nil)
	shaaksk := ""
	for _, j := range b {
		shaaksk = shaaksk + fmt.Sprintf("%02x", j)
	}
	return shaaksk, nil
}

func init() {
	globalAuthFunc = func(*http.Request) error { return nil }
}
