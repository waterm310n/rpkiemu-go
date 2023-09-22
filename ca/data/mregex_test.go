package data

import (
	"testing"
)

/* TODO 以后要是学了fuzz，再来改进测试 */

//测试URI分组匹配能不能匹配得上
func TestURIMatch(t *testing.T){
	testcases := []struct {
        in, wantProtocol,wantHost,wantCertName string
    }{
        {"rsync://rpki.ripe.net/ta/ripe-ncc-ta.cer", "rsync","rpki.ripe.net","ripe-ncc-ta.cer"},
		{"rsync://rpki.apnic.net/repository/B3A24F201D6611E28AC8837C72FD1FF2/vw5vTuDhfd6MSiS_iX0ZuHqldZ8.cer","rsync","rpki.apnic.net","vw5vTuDhfd6MSiS_iX0ZuHqldZ8.cer"},
		{"rsync://repository.lacnic.net/rpki/lacnic/f52b0c01-84ba-406f-919b-e3164fefc3f9/977f12dc643e493b8d079520804e926f4964772d.cer","rsync","repository.lacnic.net","977f12dc643e493b8d079520804e926f4964772d.cer"},
	}
	for _, tc := range testcases {
        parts := URI_MATCH.FindStringSubmatch(tc.in)
		if parts[1]!=tc.wantProtocol{
			t.Errorf("URL_MATCH : %q, want %q", parts[1], tc.wantProtocol)
		}
		if parts[2]!=tc.wantHost {
			t.Errorf("URL_MATCH : %q, want %q", parts[2], tc.wantHost)
		}
		if parts[3]!=tc.wantCertName{
			t.Errorf("URL_MATCH : %q, want %q", parts[3], tc.wantCertName)
		}
    }
}

//测试AS匹配能否在所有的AS号前加上'AS'
func TestAsnMatch(t *testing.T){
	testcases := []struct{
		in string
		want string
	}{
		{"[ 1704, 2560 , 1170 ,Min: 1768 max: 1769]","[ AS1704, AS2560 , AS1170 ,Min: AS1768 max: AS1769]"},
		{"[26173, 26194, 26210, 26218, 26317, 26418, 26426]","[AS26173, AS26194, AS26210, AS26218, AS26317, AS26418, AS26426]"},
	}
	for _,tc := range testcases{
		res := ASN_MATCH.ReplaceAllString(tc.in,"AS$1")
		if  res!= tc.want{
			t.Errorf("AsnMatch: %q ,want: %q",res,tc.want)
		}		
	}
}

//测试AS最小到最大匹配能否在将所有的Min: 1768 max: 1769转化成AS1768-AS1769
func TestAsnMinMaxMatch(t *testing.T){
	testcases := []struct{
		in string
		want string
	}{
		{"[ 1704, Min: 1768 max: 1769,  2563, Min: 2569 max: 2570]","[ AS1704, AS1768-AS1769,  AS2563, AS2569-AS2570]"},
	}
	for _,tc := range testcases{
		tc.in = ASN_MATCH.ReplaceAllString(tc.in,"AS$1")
		res := ASN_MINMAX_MATCH.ReplaceAllString(tc.in,"$1-$2")
		if res != tc.want{
			t.Errorf("AsnMinMaxMatch: %q ,want: %q",res,tc.want)
		}
	}
}

func TestIPV4MinMaxMatch(t *testing.T){
	testcases := []struct{
		in string
		want string
	}{
		{"[ Min: 140.109.0.0 max: 140.138.255.255, 168.95.0.0/16, 192.192.0.0/16,]","[ 140.109.0.0-140.138.255.255, 168.95.0.0/16, 192.192.0.0/16,]"},
	}
	for _,tc := range testcases{
		res := IPV4_MINMAX_MATCH.ReplaceAllString(tc.in,"$1-$2")
		if res != tc.want{
			t.Errorf("AsnMinMaxMatch: %q ,want: %q",res,tc.want)
		}
	}
}

func TestIPV6MinMaxMatch(t *testing.T){
	testcases := []struct{
		in string
		want string
	}{
		{"[ Min: 2001:7fa:3:: max: 2001:7fa:4:ffff:ffff:ffff:ffff:ffff, ]","[ 2001:7fa:3::-2001:7fa:4:ffff:ffff:ffff:ffff:ffff, ]"},
	}
	for _,tc := range testcases{
		res := IPV6_MINMAX_MATCH.ReplaceAllString(tc.in,"$1-$2")
		if res != tc.want{
			t.Errorf("AsnMinMaxMatch: %q ,want: %q",res,tc.want)
		}
	}
}