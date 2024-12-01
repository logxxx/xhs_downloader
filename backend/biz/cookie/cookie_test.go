package cookie

import "testing"

func TestChangeCookie(t *testing.T) {

	for i := 0; i < 10; i++ {
		req := GetCookieName(GetCookie())
		t.Logf("req%v:%v", i, req)

		ChangeCookie()
	}

}
