package cookie

var (
	cookie     = ""
	rawCookie3 = `abRequestId=4100c073-c559-5dbd-aad8-9c3571ed9e72; webBuild=4.43.0; xsecappid=xhs-pc-web; a1=1933887278ejfybalpvwpg93064jme8czjvz0i4k550000360023; webId=5876019f89cedb46ca026ad2c38e66ee; gid=yjqqYYWq0qujyjqqYYWJWV34YdVivD01UhEUhT2kjq8K4u28V6dYS8888qK88Jq8JdyJYq22; websectiga=cffd9dcea65962b05ab048ac76962acee933d26157113bb213105a116241fa6c; sec_poison_id=3efb6356-00a1-485a-b962-44ef19ed6f1c; acw_tc=0a0b142917318663602441006eda0dbbe1abf8cee14ed9d6b2c7b6a04255e1; web_session=040069b655ca1916c02e7e670f354be51806f8; unread={%22ub%22:%2267160937000000001b01085b%22%2C%22ue%22:%22673478fa0000000019019d8e%22%2C%22uc%22:28}`
	rawCookie2 = `abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; xsecappid=xhs-pc-web; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; web_session=0400697999f01a77b8f86ff0c4344ba1154db9; webBuild=4.42.2; acw_tc=5fe974ce225a7c335dead981020593d627319d23b3428bdbb44c066c3bad40d9; websectiga=82e85efc5500b609ac1166aaf086ff8aa4261153a448ef0be5b17417e4512f28; sec_poison_id=9761a1b2-6063-40e7-a930-3cb8e99a6da7; unread={%22ub%22:%22671f642c000000001600f6e6%22%2C%22ue%22:%22670cf3050000000026036a72%22%2C%22uc%22:31}`
	rawCookie  = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; webBuild=4.41.1; xsecappid=xhs-pc-web; unread={%22ub%22:%2267248721000000003c01cb8f%22%2C%22ue%22:%226720fe31000000001a01c420%22%2C%22uc%22:24}; websectiga=29098a4cf41f76ee3f8db19051aaa60c0fc7c5e305572fec762da32d457d76ae; sec_poison_id=7c260fd8-07c8-4c65-ab7b-d691d49e6edc; acw_tc=790da7f023598d78001e4b86094b417debc9f72c8abfa4fdd2df9afc8fec1fa3`
	//rawCookie  = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; xsecappid=xhs-pc-web; webBuild=4.40.2; acw_tc=285a7f043bd201132463e17a6dd3a04597f96f4623495b665535948318398e2e; unread={%22ub%22:%226712f98a000000001b011c5a%22%2C%22ue%22:%22670126e1000000002c0297f8%22%2C%22uc%22:48}; websectiga=f3d8eaee8a8c63016320d94a1bd00562d516a5417bc43a032a80cbf70f07d5c0; sec_poison_id=f3faef49-096c-4c4a-a409-e116400502b2`
	//rawCookie  = `abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; xsecappid=xhs-pc-web; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; webBuild=4.35.0; websectiga=29098a4cf41f76ee3f8db19051aaa60c0fc7c5e305572fec762da32d457d76ae; sec_poison_id=57e4b6a5-dba8-4c3a-a77f-7f521ba20348; acw_tc=2ed945131369672a44b2137468ace1efdf97f8c723fac726363060a0f692a992; unread={%22ub%22:%22646f35ff0000000013008d18%22%2C%22ue%22:%2264a276b50000000034015247%22%2C%22uc%22:26}; web_session=0400697999f01a77b8f86ff0c4344ba1154db9`
)

func init() {
	cookie = rawCookie
}

func ChangeCookie() {
	if cookie == rawCookie {
		cookie = rawCookie2
	} else {
		cookie = rawCookie
	}
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
