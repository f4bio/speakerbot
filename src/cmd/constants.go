package nipplebot

const UseDb = "dynamo"
const DefaultWorkingDirectory = "./tmp"
const LogDirectory = "./logs"
const SessionFile = "session"
const DbDirectory = "db"
const DbCollectionWorking = "_working"
const DbCollectionTargets = "targets"
const DbCollectionUsers = "users"
const MaxFollowingsCount = 100

var ProxyList = []string{
	"socks5://x8233405:Tfc8exbRdZ@proxy-nl.privateinternetaccess.com:1080",
	"socks5://195.201.36.98:1080",
	"socks5://149.249.0.210:1080",
	"socks5://78.47.225.59:9050",
}
var distribution = []int{1, 1, 1, 2, 2, 3, 4, 5, 4, 3, 2, 2, 1, 1, 1}
var CfgFile = "config.json"
