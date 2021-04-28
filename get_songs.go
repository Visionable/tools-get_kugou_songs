
type InternetSong struct {
	Name   string
	Singer string
}

func InsertSongLibrary() {
	defer common.CheckGoPanic()
	appzaplog.Info("InsertSongLibrary start")
	maxPage := 23
	songList := make([]*model.InternetSong, 0)
	for i := 1; i <= maxPage; i++ {
		songList = append(songList, GetSongs(i)...)
		if len(songList) > 80 || i == maxPage { // 处理一次插表
			log.Info("songList", zap.Any("songList", songList))
       // insert to db
			songList = make([]*InternetSong, 0)
		}
	}
	appzaplog.Info("InsertSongLibrary end")
}

func GetSongs(page int) []*InternetSong {
	reqUrl := fmt.Sprintf(`https://www.kugou.com/yy/rank/home/%d-8888.html?from=rank`, page)
	songList := make([]*InternetSong, 0)
	doc, err := htmlquery.LoadURL(reqUrl)
	if err != nil {
		log.Error("GetSongs error", zap.Error(err))
		return songList
	}
	list, _ := htmlquery.QueryAll(doc, "//div/ul/li[@title]/a")
	for _, v := range list {
		text := htmlquery.InnerText(v)
		if len(text) <= 0 {
			continue
		}
		textList := strings.Split(text, "-")
		if len(textList) < 2 {
			continue
		}
		singer, name := textList[0], textList[1]
		name = strings.Trim(name, " ")
		singer = strings.Trim(singer, " ")
		if len(name) <= 0 || len(singer) <= 0 {
			continue
		}
		// 歌名带括号的丢弃
		if strings.Index(name, "(") != -1 || strings.Index(name, "（") != -1 {
			continue
		}
		songList = append(songList, &InternetSong{Name: name, Singer: singer})
	}
	return songList
}
