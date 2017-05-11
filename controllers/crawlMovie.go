package controllers
import(
	"github.com/astaxie/beego"
	"crawl_movie/models"
	"github.com/astaxie/beego/httplib"
	"time"
)

type CrawlMovieController struct {
	beego.Controller
}
func (c *CrawlMovieController) CrawlMovie(){
	//连接Redis
	models.ConnectRedis("127.0.0.1:6379")
	//爬虫入口Url
	sUrl := "https://movie.douban.com/subject/26260853/"
	
	//加入队列中
	models.PutinQueue(sUrl)

	var movieInfo models.MovieInfo

	//循环抓取爬虫
	for {
		length := models.GetQueueLength()
		if length == 0 {
			break
		}

		sUrl = models.PopfromQueue()

		//先判断sUrl是否已被访问过
		if models.IsVisit(sUrl) {
			continue
		}

		rsp := httplib.Get(sUrl)
		//设置User-agent以及cookie防止豆瓣网的403
		rsp.Header("User-Agent","Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.104 Safari/537.36 Core/1.53.2658.400 QQBrowser/9.6.11018.400")
		rsp.Headder("Cookie",`bid=trmhFGS4iyE; viewed="27016236"; gr_user_id=404ade19-68ec-466e-ab99-493fae87afd7; ll="118281"; ap=1; _vwo_uuid_v2=B1E37876D385A22927E6937005428FA9|9a64eca43d79ff9348a71e0adfa0e859; __utmt=1; __utma=30149280.1502636346.1489547064.1493702857.1493714395.12; __utmb=30149280.1.10.1493714395; __utmc=30149280; __utmz=30149280.1493395962.9.8.utmcsr=baidu|utmccn=(organic)|utmcmd=organic`)
		sMovieHtml,err := rsp.String()
		if err != nil{
			panic(err)
		}
		//记录电影信息
		movieInfo.Movie_name = models.GetMovieName(sMovieHtml)
		if movieInfo.Movie_name != ""{
			movieInfo.Movie_director = models.GetMovieDirector(sMovieHtml) 
			movieInfo.Movie_main_character = models.GetMovieMainCharacters(sMovieHtml) 
			movieInfo.Movie_type = models.GetMovieGenre(sMovieHtml)
			movieInfo.Movie_on_time = models.GetMovieOnTime(sMovieHtml)
			movieInfo.Movie_grade = models.GetMovieGrade(sMovieHtml)
			movieInfo.Movie_span = models.GetMovieRunningTime(sMovieHtml)

			models.AddMovie(&movieInfo)
		}

		//提取该页面的所有信息
		urls := models.GetMovieUrls(sMovieHtml)
		
		for _,url := range urls {
			models.PutinQueue(url)
			c.Ctx.WriteString("<br>" + url + "</br>")
		}

		//sUrl应当记录到访问set中
		models.AddToSet(sUrl)

		//爬虫不能太快，避免被发现
		time.Sleep(3*time.Second)	
	}

	c.Ctx.WriteString("end the Crwal")
	

	//id ,_:= models.AddMovie(&movieInfo)
	//c.Ctx.WriteString(fmt.Sprintf("%v",id))


	/*c.Ctx.WriteString(models.GetMovieDirector(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieName(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieMainCharacters(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieGrade(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieGenre(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieOnTime(sMovieHtml) + "\n")
	c.Ctx.WriteString(models.GetMovieRunningTime(sMovieHtml) + "\n")*/
}