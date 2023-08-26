package controller

var NgrokHost string

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type User struct {
	//Id            int64  `json:"id,omitempty"`
	Id          int64  `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Name        string `json:"name,omitempty"`
	FollowCount int64  `json:"follow_count,omitempty"`
	Follows     string `json:"follows,omitempty" gorm:"type:text;"`
	//Followws        []*User `gorm:"many2many:follows;"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	Followers       string `json:"followers,omitempty" gorm:"type:text;"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Password        string `json:"password,omitempty"`
	Token           string `json:"token,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  string `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	// 用户收藏的视频ID数组
	FavoritedVideos string `json:"favorited_videos,omitempty" gorm:"type:text;"`
	//TableName     string `gorm:"tablename:users"`
}

/*func (u *User) BeforeSave() error {
	// 序列化为 JSON 格式存储
	bytes, err := json.Marshal(u.FavoritedVideos)
	if err != nil {
		return err
	}
	u.FavoritedVideos = string(bytes)
	return nil
}

func (u *User) AfterFind() error {
	// 读取后需要反序列化
	if err := json.Unmarshal([]byte(u.FavoritedVideos), &u.FavoritedVideos); err != nil {
		return err
	}
	return nil
}*/

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Author        User   `json:"author" gorm:"ForeignKey:UserID"`
	UserID        int64  `json:"user_id"` // 用户ID
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	PublishTime   string `json:"publish_time,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user" gorm:"ForeignKey:UserID"`
	UserID     int64  `json:"user_id"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoID    int64  `json:"video_id"`
	Video      Video  `json:"video" gorm:"ForeignKey:VideoID"`
}

type Message struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	//CreateTime string `json:"create_time,omitempty"`
	CreateTime int64 `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
