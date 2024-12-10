package blog

import (
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"os"
	"testing"
)

func TestGetHomePage(t *testing.T) {

	xs := `XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6IjdkZDRkMjY2YjFjZGFmNTJlOTZjZGM2ZTY3ZDVjZjE4NWI3OGJkMzdkZjk2ODUzMWFlMzIwMTg4MDNjZDcwOTM1MGY4NTA5MGU0ZDAxMjdjNjgwMzU5MzI1MmQ0MDZmZmU2MjAxOGZhZmFkNDhjYTU0ZWQxY2VhZWQ0YzQzNTA2YmQ0MzViYmIzNzdkMDU4ZWNhYTkwODNkMmQ4YTZlMWJiMmY2OTRiNGE2MDQ3ZWVmYzZjYjFhYmRlZGE1NDg3MjExMWFkOWE1NDA2NGEyYjI5ZTViMDdmMWFjZWE0MDJlMmQyNGQyN2M2ZmI0MmM2ZGEzY2Q2N2MyMzczNTY2ZDFjNTIyYjE5MWJjMmM3MTUwYmNlNDE0Y2NmZDYxYWFiZTAxMTg1MmIxZWY0MWY4MjNlN2MwZGQ0ZGFhYjBiOTk1ODc5M2ZlYWVmZmE5MWU3Y2ZkNGI4Nzg4OWJiZDFmYWNmYjAyYTNkODk3YzhmODFiM2Q2YTYxZGU2Y2NhNjQ1YTIzZDkwYjMwNDAyNjc0OWM1MWI3NDA4OGRkY2QxZjQ3In0=`

	resp, err := GetHomePage(cookie.GetCookie1(), xs)
	if err != nil {
		t.Fatal(err)
	}
	fileutil.WriteJsonToFile(resp, "D:\\mytest\\mywork\\xhs_downloader\\backend\\biz\\blog\\home_test5.json")

}

func TestGetHomePage2(t *testing.T) {
	htmlData, err := os.ReadFile("D:\\mytest\\mywork\\xhs_downloader\\backend\\biz\\blog\\home_test.html")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := parseHomeNotes(string(htmlData))
	if err != nil {
		t.Fatal(err)
	}

	fileutil.WriteJsonToFile(resp, "D:\\mytest\\mywork\\xhs_downloader\\backend\\biz\\blog\\home_test.json")
}

func TestGetHomePage3(t *testing.T) {
	notes := []blogmodel.ParseNote{}
	fileutil.ReadJsonFile("D:\\mytest\\mywork\\xhs_downloader\\backend\\biz\\blog\\home_test4.json", &notes)
	for _, e := range notes {
		parseResp, err := ParseBlog(e.URL, "")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("parse resp:%+v", parseResp)
		os.Exit(1)
	}
}
