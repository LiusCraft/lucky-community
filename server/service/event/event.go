package event

const (
	CommentUpdateEvent = iota + 1 // 文章下评论更新事件
	UserFollowingEvent            // 用户关注的人事件
	ArticleAt                     // 文章中@
	CommentAt                     // 评论中@
	ReplyComment                  // 评论回复
	Adoption                      // 采纳
	SectionComment                // 章节评论
	CourseComment                 // 课程回复
	CourseUpdate                  // 课程更新
	Meeting                       // 会议
)

var events []*event

var eventMap = make(map[int]string)

var eventPage = make(map[int]string)

func init() {
	events = append(events, nil)
	events = append(events, &event{Id: CommentUpdateEvent, Msg: "文章评论"})
	events = append(events, &event{Id: UserFollowingEvent, Msg: "用户更新"})
	events = append(events, &event{Id: ArticleAt, Msg: "文章 @"})
	events = append(events, &event{Id: CommentAt, Msg: "评论 @"})
	events = append(events, &event{Id: ReplyComment, Msg: "评论回复"})
	events = append(events, &event{Id: Adoption, Msg: "采纳"})
	events = append(events, &event{Id: SectionComment, Msg: "章节回复"})
	events = append(events, &event{Id: CourseComment, Msg: "课程回复"})
	events = append(events, &event{Id: CourseUpdate, Msg: "课程更新"})
	events = append(events, &event{Id: Meeting, Msg: "分享会"})

	eventMap[CommentUpdateEvent] = "文章评论"
	eventMap[UserFollowingEvent] = "用户更新"
	eventMap[ArticleAt] = "文章 @"
	eventMap[CommentAt] = "评论 @"
	eventMap[ReplyComment] = "评论回复"
	eventMap[Adoption] = "采纳"
	eventMap[SectionComment] = "章节回复"
	eventMap[CourseComment] = "课程回复"
	eventMap[CourseUpdate] = "课程更新"
	eventMap[Meeting] = "分享会"

	eventPage[CommentUpdateEvent] = "articleView"
	eventPage[UserFollowingEvent] = "articleView"
	eventPage[ArticleAt] = "articleView"
	eventPage[CommentAt] = "articleView"
	eventPage[ReplyComment] = "articleView"
	eventPage[Adoption] = "articleView"
	eventPage[SectionComment] = "articleView"
	eventPage[CourseComment] = ""
	eventPage[CourseUpdate] = ""
	eventPage[Meeting] = ""

}

// 事件
type event struct {
	Id  int    `json:"id"`
	Msg string `json:"msg"`
}

func GetMsg(eventId int) string {
	v := events[eventId]
	return v.Msg
}

func List() []*event {
	return events
}
func Map() map[int]string {
	return eventMap
}

func PageName() map[int]string {
	return eventPage
}
