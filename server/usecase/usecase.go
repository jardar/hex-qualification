package usecase

import (
	"errors"
	"postapi/entity"
	"postapi/utils"
	"sync/atomic"
)

type PostBookmark struct {
	*entity.Post
	Bookmark bool `json:"bookmark"`
}

var (
	memberDb   map[string]string       //map[email] pass
	postDb     map[uint32]*entity.Post //map[post id] *Post
	bookmarkDb map[string][]uint32     //map[email][]post-id

	postIdCnt uint32
)
var (
	_allPost = []*entity.Post{
		{
			1,
			"南鯤鯓代天府",
			"是一座位於臺南市北門區之王爺廟，為台灣五府千歲的王爺總廟，主祀「代天巡狩」李、池、吳、朱、范 府千歲五府千歲",
			`全臺規模最大的王爺信仰中心\r\n南鯤鯓代天府位於北門沿海，傳說是當地漁民撿到了載有五尊神像的小船，於是將神像供奉在草寮內，早晚焚香膜拜，結果從那天起，漁民只要出海捕魚，一定滿載而歸。事蹟傳開後，香火日益鼎盛，並擇地建廟，逐漸成為當地的信仰中心。日治時期的重修，在石雕、彩繪、剪黏上皆展現高度的工藝水準，而被評定為國定古蹟。 國定民俗「南鯤鯓代天府五府千歲進香期」 五府千歲的靈驗事蹟，吸引了各地民眾前來分靈，至今全臺已有17000餘座的分靈廟。\r\n\r\n每到王爺誕辰（農曆4月中下旬、6月中旬、8月中旬、9月中旬），分靈廟便會籌組進香團前來進香謁祖，可以說是臺灣最為壯盛的王爺廟會。 到江南園林「大鯤園」一遊 由建築家漢寶德先生所設計的「大鯤園」，為一座仿中國江南式的園林，除了有優美的山水，也有傳統建築，園區內設有南鯤鯓文史館，可一覽廟宇沿革與文物。\r\n\r\n在許多台南人的記憶當中，對南鯤鯓代天府的印象就是超熱鬧的廟會，每天參拜進香的信徒及慕名而來的遊客不計其數，主要祀奉李、池、吳、朱、范等五位王爺，五位王爺與囝仔公爭地建廟的神話故事至今依然讓很多人津津樂道著，沿海的王爺信仰在這裡展露無遺，2012年11月更是完成了一件創舉，將建廟以來信徒捐贈的大大小小金牌加上廟方添購，共集結了一萬零八百兩黃金，在廟內的凌霄寶殿打造了總價約6億多，也是全球最大的黃金玉旨牌，從1984年開始，至2012年完成，可說是建廟300多年來的超大盛事。\r\n\r\n每年固定舉辦的平安鹽祭，在10~11月間就會陸陸續續的舉辦一連串的活動，每年樣子都不同的平安鹽袋更是大家排隊爭著領取收集的熱門紀念品，拿到平安鹽袋後，到廟前廣場的鹽堆拿些鹽裝入，再到香爐過爐，就是納福驅邪的小寶物。​\r\n\r\n▲友善空間：\r\n無障礙停車位：11格\r\n無障礙廁所：5間\r\n輪椅租借：10台\r\n&nbsp;\r\n主祀：李府千歲（大王）、池府千歲（二王）、吳府千歲（三王）、朱府千歲（四王）、范府千歲（五王）\r\n配祀：玉皇上帝、觀音佛祖、萬善爺（囝仔公）、中軍府、城隍爺、虎將軍、地藏王、註生娘娘、福德正神、月老神君\r\n\r\n藝品導覽：三川殿中港間「石堵」\r\n鑑賞重點：石雕技法的立體教科書\r\n\r\n南鯤鯓代天府三川殿的立面，全以石雕構成，不同於其他廟宇的是，由於廟靠近海邊，空氣中多鹽份，因此石雕會做上彩保護，靠近可以看到淺淺的色彩。\r\n\r\n除了上彩的地方特色外，每一片石雕匠師都用了不同的雕造技法。以中門為例，最上層的「知章騎馬」是透雕，將背景剔除，只保留圖案；第二層的「孔明收姜維」屬於「深浮雕」，層次較多、起伏較大的技法；第三層的「長坂坡」與「孔明舌戰群儒」用的是俗稱「內枝外葉」的技法，將石材鑿穿，但平面背景相接，在功能性上也兼具窗戶通風的效果；接在前三堵的熱鬧場景後，第四層以「陰刻」雕出花鳥形狀，在視覺上也有喘息效果；最後一層則是以「淺浮雕」雕出麒麟、鳳凰、牡丹，起伏較為平緩。\r\n「知章騎馬」使用透雕技法\r\n\r\n「孔明收姜維」使用深浮雕技法\r\n\r\n\r\n「長坂坡」與「孔明舌戰群儒」使用內枝外葉技法\r\n\r\n「杏林春燕．賜福添壽」使用陰刻技法\r\n\r\n「三王獻瑞：麒麟、鳳凰、牡丹」使用淺浮雕技法\r\n\r\n藝品導覽：正殿「蜘蛛結網藻井」\r\n鑑賞重點：在天花板結出蜘蛛網\r\n\r\n在南鯤鯓代天府的神龕前抬頭一看，會發現別有洞天&mdash;&mdash;這是藻井，傳統廟宇建築室內頂棚常見的形式之一。藻井的做法，是由一組組木頭「斗拱」組件排列成八角形或圓形，層層向上堆疊，直到屋頂處，數目眾多，令人眼花撩亂。\r\n\r\n這樣子的造型可不是匠師自由發揮，而是從《周易》64卦384爻，配合一年四季十二月令、二十四節氣、七十二候所形成的卦氣卦候圖而來。構成的形狀如同蜘蛛結網，看起來繁複又華麗，因此是觀察是匠師功力的絕佳所在。\r\n註：龕（ㄎㄢ）\r\n\r\n蜘蛛結網藻井\r\n\r\n藝品導覽：正殿「金錢壁」\r\n鑑賞重點：摸牆壁求財氣\r\n\r\n在正殿後方、青山寺的對面，有一堵人氣頗旺的牆壁，不同於一般廟宇以磚砌造，而是用澎湖運來的硓𥑮石砌成，建造於日大正15年（1926），由澎湖西嶼內塹宮的信眾所捐獻，見證了兩地間王爺信仰的密切關係。\r\n\r\n這堵牆全長6尺2乘於10尺56（約182*356公分），石材被砌成「八卦龜錦紋」樣式，又似古錢幣，枚枚相疊，塊塊相扣，而且做工細緻，幾乎看不到疊痕與縫隙，是在手工雕琢時代的佳作。由於烏龜與錢幣分別象徵長壽與財富，從下方被摸得光亮的地方就可知人氣之高。\r\n金錢壁正面一景\r\n\r\n金錢壁上的「八卦龜錦紋」樣式\r\n\r\n藝品導覽：廊道「潘麗水樑枋彩繪」\r\n鑑賞重點：國寶級匠師的彩繪廊道\r\n\r\n聚集了來自全臺的彩繪名師在此獻藝，更有「對場做」（兩匠師於不同空間、同一裝飾部位進行彩繪，較量技巧高下）的競技，使得南鯤鯓代天府有著很高的彩繪水準。\r\n\r\n而在這些作品中，高度最低、最容易欣賞的，莫過於薪傳獎藝師潘麗水在正殿兩側廊道留下的樑枋彩繪「曹操進劍」的故事來自《三國演義》，曹操以進劍之名欲暗殺專政的董卓；「舉案齊眉」出自《後漢書》，描述夫妻間相敬如賓；「驪姬巧計害申生」出自《東周列國志》，畫面中寵妾驪姬將蜂蜜淋在頭髮上，並請她欲除掉的世子申生協助驅趕蜜蜂，塑造出申生非禮的畫面，使得晉獻公廢除申生的繼承資格。在這18爿共36幅，幅幅相連的彩繪廊道裡，可以看到畫師對於文學、神話傳說、歷史故事的深厚涵養。\r\n潘麗水樑枋彩繪\r\n\r\n藝品導覽：正殿屋頂「屋脊剪黏」\r\n鑑賞重點：屋脊下的栩栩如生\r\n\r\n南鯤鯓代天府的屋頂屬於「歇山重簷式」，在每一條屋脊收尾的地方，匠師安排了各式剪黏，有人物、有動物，也有博古圖，使得屋頂看起來熱鬧非凡。\r\n\r\n中脊的燕尾下方，裝飾的是「仙翁騎鶴」，仙鶴以百片以上的瓷片組成，老翁騎在上方，跟著上揚的屋脊一起騰空；而其他的屋脊收尾處，則可以見到一手舞劍、一手拿扇的仙女，甚至是人身猴臉的孫悟空，一手拿著金箍棒，一手向上高舉，目光炯炯的望向遠方。除了人物外，也可以見到俯衝狀的「倒拋獅」、寓意「四季平安」的花瓶，都是安平剪黏名匠「葉鬃」佈置在屋脊下方的要角。\r\n「仙女舞劍」屋脊剪黏\r\n\r\n「孫悟空」屋脊剪黏\r\n\r\n「仙翁騎鶴」屋脊剪黏\r\n\r\n\r\n藝品導覽：正殿水車堵「蘆花河」剪黏\r\n鑑賞重點：走一圈欣賞隱藏版的剪黏\r\n\r\n提到剪黏，第一印象往往是在屋頂上華麗繁複的裝飾，而在南鯤鯓代天府，除了屋頂上，環繞著正殿外牆、連綿不斷的水車堵剪黏，則是由安平名匠葉鬃所作、隱藏版的欣賞重點之一。\r\n\r\n水車堵細長的空間，特別適合角色眾多的武戲場景。就算不熟悉故事的來歷，只著肚兜的蚌殼精、岸邊射箭的武將，每個人偶奮力演繹的模樣，還是值得細細品味。\r\n\r\n水車堵「蘆花河」剪黏\r\n\r\n鑑賞重點：著名的齣目「蘆花河」\r\n\r\n常見於各廟宇的「蘆花河」，說的是河神薛應龍的故事。蘆花河神薛應龍，轉世為山賊後，成為樊梨花義子協助作戰，戰亡後魂魄想要返回蘆花河，卻發現河已被惡龍所佔。不敵惡龍的薛應龍，只好托夢義父母薛丁山和樊梨花。於是薛丁山協同樊梨花等三位夫人前來助戰，瞬間整條河殺聲漫天，水翻波滾中，薛丁山一箭射中惡龍，薛應龍終於重返蘆花河。\r\n「蘆花河剪黏」特寫\r\n\r\n「蘆花河剪黏」特寫`,
			"https://unsplash.com/photos/_RBcxo9AU-U",
		},
		{
			2,
			"黑面琵鷺生態展示館",
			"在通往黑面琵鷺保護區的廣闊魚塭之中，會見到一棟造型非常特殊的水上屋，是臺南市重要的環境教育設施場所",
			`在通往黑面琵鷺保護區的廣闊魚塭之中，會見到一棟造型非常特殊的房屋，優雅的座落於水面上，這是台南第一座「水上屋」，也是臺南市重要的環境教育設施場所，更是七股黑面琵鷺保護區的重要據點。\r\n\r\n\r\n黑面琵鷺生態展示館採用綠建築的建築工法，不但節能減碳也融入自然環境當中，讓嬌貴且易受驚嚇的鳥類們不會因為覓食地中突然冒出的建築物而嚇得遷移棲地，展示館空間包含「常態展示區」、「多媒體室」、「特展區」、「會議室」及「戶外觀景平台」，每天夕陽西下時總是有許多賞鳥完的遊客來此繼續欣賞美麗的夕陽，在此取景拍夕陽還能有美麗的生態展示館當作前景，怎麼拍都好看，且距離七股的各景點都非常近，是來七股黑面琵鷺保護區旅行必來的中心景點。`,
			"https://unsplash.com/photos/_RBcxo9AU-U",
		},
		{
			3,
			"頑皮世界野生動物園",
			"園區內模擬原棲環境圈養方式，以低矮柵欄展示，讓人與動物的距離縮為零；各種不同主題動物區，可愛動物區、非洲動物觀賞區",
			`頑皮世界是台南目前唯一的動物園，也是台南相當老字號的動物園，照顧動物的用心及經營的理念都頗受遊客們支持喜愛，從踏入園區就可以感受到與動物的親近感，大門口的金剛鸚鵡搖頭晃腦的歡迎著你，動物園裡超人氣的明星水豚、長頸鹿、企鵝、藪貓、狐獴等也都有機會可以近距離的觀賞，是一處以生態教育為理念的專業動物園，過往的動物表演雖然已走入歷史，但園方現在傳達的動物友善及生命教育宣導卻更深入遊客的心中。\r\n\r\n園區內除了與動物親近以外也有不少的遊樂設施可以玩，隨處都有的涼亭及座椅讓遊客可以隨時找到地方休息，相當貼心，也有餐廳等可以隨時填飽五臟廟。人與動物之間，相處久了真的會有很微妙的感情與默契存在著，動物園不只是小朋友戶外教學的地方，喜歡動物的情侶朋友也可以來走走，體驗一下跟動物的互動，或許可以讓疲憊的身心獲得一些療癒。`,
			"https://unsplash.com/photos/_RBcxo9AU-U",
		},
	}
)

func init() {
	memberDb = make(map[string]string)
	memberDb["user@ddt.com"] = "123"  // 一般用戶
	memberDb["admin@ddt.com"] = "456" // 後台管理者

	postDb = make(map[uint32]*entity.Post)
	bookmarkDb = make(map[string][]uint32)

	// init post data
	for _, post := range _allPost {
		postDb[post.Id] = post
	}
	//
	postIdCnt = uint32(len(postDb))
}

func ReadPostsBookmarked(email string) ([]*PostBookmark, error) {
	ret := make([]*PostBookmark, 0)
	if v, ok := bookmarkDb[email]; ok {
		for _, postId := range v {
			if p, err := ReadPostById(postId); err != nil {
				continue
			} else {
				ret = append(ret, &PostBookmark{p, true})
			}
		}
	}

	return ret, nil
}
func ToggleBookmark(email string, postId uint32) bool {
	finalState := false
	if v, ok := bookmarkDb[email]; !ok {
		finalState = true
		bookmarkDb[email] = make([]uint32, 1)
		bookmarkDb[email][0] = postId
	} else {
		foundIdx := utils.SliceFindIndex(postId, v)
		if foundIdx == -1 {
			finalState = true
			bookmarkDb[email] = append(bookmarkDb[email], postId)
		} else {
			bookmarkDb[email] = append(bookmarkDb[email][:foundIdx], bookmarkDb[email][foundIdx+1:]...)
		}
	}

	return finalState
}

func ReadPostById(postId uint32) (*entity.Post, error) {

	if v, ok := postDb[postId]; !ok {
		return nil, errors.New("no data")
	} else {
		return v, nil
	}

}

func ReadPosts() ([]*entity.Post, error) {

	ret := make([]*entity.Post, 0)
	for _, v := range postDb {
		ret = append(ret, v)
	}

	return ret, nil
}

func DeletePost(postId uint32) (*entity.Post, error) {
	if _, ok := postDb[postId]; !ok {
		return nil, errors.New("not found")
	}

	oldPost := postDb[postId]
	delete(postDb, postId)

	return oldPost, nil
}

func UpdatePost(post *entity.Post) (*entity.Post, error) {

	if v, ok := postDb[post.Id]; ok {
		v.Title = post.Title
		v.Summary = post.Summary
		v.Body = post.Body

		return v, nil
	} else {
		return nil, errors.New("not found")
	}

}

func CreatePost(post *entity.Post) (*entity.Post, error) {

	id := atomic.AddUint32(&postIdCnt, 1)

	postDb[id] = &entity.Post{
		Id:      id,
		Title:   post.Title,
		Summary: post.Summary,
		Body:    post.Body,
		PicUrl:  post.PicUrl,
	}

	return postDb[id], nil
}

func GetBookmarkState(email string, postId uint32) bool {
	if _, ok := bookmarkDb[email]; !ok {
		return false
	}
	for _, bkPostId := range bookmarkDb[email] {
		if bkPostId == postId {
			return true
		}
	}

	return false
}
func Login(email, pass string) error {
	if _, ok := memberDb[email]; !ok {
		return errors.New("not exist")
	}
	if memberDb[email] != pass {
		return errors.New("invalid user")
	}

	return nil
}
func SignUp(email, pass string) error {
	if _, ok := memberDb[email]; ok {
		return errors.New("existed")
	}
	memberDb[email] = pass

	return nil
}
