package services

import "xhyovo.cn/community/server/dao"

var (
	messageDao       dao.MessageDao
	articleDao       dao.Article
	fileDao          dao.File
	codeDao          dao.InviteCode
	typeDao          dao.Type
	userDao          dao.UserDao
	commentDao       dao.CommentDao
	subscriptionDao  dao.SubscriptionDao
	memberDao        dao.MemberDao
	aiNewsDao        dao.AiNewsDao
	crawlerConfigDao dao.CrawlerConfigDao
	// 积分系统相关DAO实例
	userPointsDao    dao.UserPointsDao
	pointRecordDao   dao.PointRecordDao
	inviteRelationDao dao.InviteRelationDao
	exchangeRequestDao dao.ExchangeRequestDao
	pointProductDao  dao.PointProductDao
	userInviteCodeDao dao.UserInviteCodeDao
)
