package cookie

var (
	cookie = ""

	rawCookie3 = `abRequestId=b57c32b4-f9fe-5f1c-81eb-8b624acef89e; webBuild=4.44.1; a1=19353fe73aazasu255vk1ggk7hiyih7emivgecszy50000207887; webId=27646553690137cc0d7c835863dc63e2; gid=yjq2qidj0d6Dyjq2qidWq3SM00u0M7J22hTyhuvkkTWx3k28v3xWdd888J8WYYW840KDi4Y2; web_session=040069b655ca1916c02e9cca75354b85c67f26; xsecappid=xhs-pc-web; unread={%22ub%22:%22673ec6cf0000000002038148%22%2C%22ue%22:%22673eb3eb0000000007027b15%22%2C%22uc%22:23}; websectiga=59d3ef1e60c4aa37a7df3c23467bd46d7f1da0b1918cf335ee7f2e9e52ac04cf; sec_poison_id=a854b8b9-09ae-4c63-a466-448edf61b22b; acw_tc=0a0bb22a17323400698216979e1dc680eece4405ff584ffe1544f27fa264e7`
	rawCookie2 = `abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; xsecappid=xhs-pc-web; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; web_session=0400697999f01a77b8f86ff0c4344ba1154db9; webBuild=4.42.2; acw_tc=5fe974ce225a7c335dead981020593d627319d23b3428bdbb44c066c3bad40d9; websectiga=82e85efc5500b609ac1166aaf086ff8aa4261153a448ef0be5b17417e4512f28; sec_poison_id=9761a1b2-6063-40e7-a930-3cb8e99a6da7; unread={%22ub%22:%22671f642c000000001600f6e6%22%2C%22ue%22:%22670cf3050000000026036a72%22%2C%22uc%22:31}`
	rawCookie  = `abRequestId=0a2b07d0-9a27-57d2-a5b3-9eb833257b5d; webBuild=4.46.0; xsecappid=xhs-pc-web; a1=193a47fc82bd2qy48tfefdagnwjgvlwf89v4xq2yg50000343212; webId=fc14a13205727e6bc0b19525731a3b4b; websectiga=634d3ad75ffb42a2ade2c5e1705a73c845837578aeb31ba0e442d75c648da36a; acw_tc=0a0bb37217336316598126774e2184b5f0212c58adb119aca891e9304f9244; gid=yjq04WifK8Yyyjq04WiSYjAJJDfJAv4Y9idihvkf0kIEV428kh1EiC888q4qJyJ8yfiqYjq0; sec_poison_id=cf425447-54c6-476a-b914-1d17028b2c22; web_session=040069b5ac82da56bef8076a60354b0bed6c58; unread={%22ub%22:%22674eb4ed000000000202dcc4%22%2C%22ue%22:%226750783300000000070378a4%22%2C%22uc%22:22}`
	//rawCookie  = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; xsecappid=xhs-pc-web; webBuild=4.40.2; acw_tc=285a7f043bd201132463e17a6dd3a04597f96f4623495b665535948318398e2e; unread={%22ub%22:%226712f98a000000001b011c5a%22%2C%22ue%22:%22670126e1000000002c0297f8%22%2C%22uc%22:48}; websectiga=f3d8eaee8a8c63016320d94a1bd00562d516a5417bc43a032a80cbf70f07d5c0; sec_poison_id=f3faef49-096c-4c4a-a409-e116400502b2`
	//rawCookie  = `abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; xsecappid=xhs-pc-web; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; webBuild=4.35.0; websectiga=29098a4cf41f76ee3f8db19051aaa60c0fc7c5e305572fec762da32d457d76ae; sec_poison_id=57e4b6a5-dba8-4c3a-a77f-7f521ba20348; acw_tc=2ed945131369672a44b2137468ace1efdf97f8c723fac726363060a0f692a992; unread={%22ub%22:%22646f35ff0000000013008d18%22%2C%22ue%22:%2264a276b50000000034015247%22%2C%22uc%22:26}; web_session=0400697999f01a77b8f86ff0c4344ba1154db9`
	isCookieDisabled = []bool{false, false, false}
	cookies          = []string{rawCookie, rawCookie2, rawCookie3}
)

func init() {
	cookie = rawCookie
}

func ChangeCookie() {
	idx := 0
	switch cookie {
	case rawCookie:
		idx = 0
	case rawCookie2:
		idx = 1
	case rawCookie3:
		idx = 2
	}

	idx += 1
	idx = idx % len(cookies)

	if isCookieDisabled[idx] == true {
		idx += 1
		idx = idx % len(cookies)
	}

	cookie = cookies[idx]
}

func SetCookie1Disabled() {
	isCookieDisabled[0] = true
}

func GetCookie() string {
	return cookie
}

func GetCookie1() string {
	return rawCookie
}

func GetCookie2() string {
	return rawCookie2
}

func GetCookie3() string {
	return rawCookie3
}

func GetCookieName(input string) string {
	switch input {
	case GetCookie1():
		return "cookie1"
	case GetCookie2():
		return "cookie2"
	case GetCookie3():
		return "cookie3"
	default:
		return "UNKNOWN:" + input
	}
}
