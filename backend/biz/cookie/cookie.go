package cookie

var (
	cookie     = ""
	rawCookie2 = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; xsecappid=xhs-pc-web; webBuild=4.35.0; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; unread={%22ub%22:%2266f04d37000000001201221c%22%2C%22ue%22:%2266f0c934000000000c01a94f%22%2C%22uc%22:28}; acw_tc=0037368b76eb084b57220da6544683b6c66ab107526d14e9f88dc22415626e2a; websectiga=2a3d3ea002e7d92b5c9743590ebd24010cf3710ff3af8029153751e41a6af4a3; sec_poison_id=70aa5d78-36b6-471c-805c-5185f5f12273`
	rawCookie  = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; xsecappid=xhs-pc-web; webBuild=4.40.2; acw_tc=285a7f043bd201132463e17a6dd3a04597f96f4623495b665535948318398e2e; unread={%22ub%22:%226712f98a000000001b011c5a%22%2C%22ue%22:%22670126e1000000002c0297f8%22%2C%22uc%22:48}; websectiga=f3d8eaee8a8c63016320d94a1bd00562d516a5417bc43a032a80cbf70f07d5c0; sec_poison_id=f3faef49-096c-4c4a-a409-e116400502b2`
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
