package services

import (
	"errors"
	"strings"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type CourseService struct {
}

// 发布课程
func (*CourseService) Publish(course model.Courses) {
	course.Technology = strings.Join(course.TechnologyS, ",")
	if course.ID == 0 {
		model.Course().Save(&course)
	} else {
		model.Course().Where("id = ?", course.ID).Updates(&course)
	}
}

// 获取课程详细信息
func (*CourseService) GetCourseDetail(id int) *model.Courses {
	var course *model.Courses
	model.Course().Where("id = ?", id).Find(&course)
	course.TechnologyS = strings.Split(course.Technology, ",")
	return course
}

// 获取课程列表
func (*CourseService) PageCourse(page, limit int) (courses []model.Courses, count int64) {
	model.Course().Offset((page - 1) * limit).Limit(limit).Order("created_at desc").Find(&courses)
	model.Course().Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourse(id int) {
	model.Course().Delete("id = ?", id)
	model.CoursesSection().Where("course_id = ?", id).Delete(&model.CoursesSections{})
}

// 发布章节
func (c *CourseService) PublishSection(section model.CoursesSections) error {
	if c.GetCourseDetail(section.CourseId).ID == 0 {
		return errors.New("对应课程不存在")
	}
	if section.ID == 0 {

		model.CoursesSection().Save(&section)
		var b SubscribeData
		var subscriptionService SubscriptionService
		b.UserId = section.UserId
		b.CurrentBusinessId = section.CourseId
		b.SubscribeId = section.CourseId
		b.CourseId = section.CourseId
		b.SectionId = section.ID
		subscriptionService.Do(event.CourseUpdate, b)
	} else {
		model.CoursesSection().Where("id = ?", section.ID).Updates(&section)
	}
	return nil
}

// 获取章节详细信息
func (*CourseService) GetCourseSectionDetail(id int) *model.CoursesSections {
	var sections *model.CoursesSections
	model.CoursesSection().Where("id = ?", id).Find(&sections)
	var userS UserService
	sections.UserSimple = userS.GetUserSimpleById(sections.UserId)
	return sections
}

// 获取课程列表
func (*CourseService) PageCourseSection(page, limit, courseId int) (courses []model.CoursesSections, count int64) {
	model.CoursesSection().Limit(limit).Offset((page-1)*limit).Where("course_id = ? ", courseId).Select("id", "title").Find(&courses)
	model.CoursesSection().Where("course_id = ? ", courseId).Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourseSection(id int) {
	model.CoursesSection().Delete("id = ?", id)
	// 对应评论一并删除 todo
}
