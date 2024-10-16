package rules

import (
	"article/pkg/model"
	"article/pkg/tools"
	"gorm.io/gorm"
)

type Enforcer struct {
	db *gorm.DB
}

func NewEnforcer(db *gorm.DB) *Enforcer {

	return &Enforcer{
		db: db,
	}
}

func (e *Enforcer) AddGroup(uid, policy string) error {
	// 判断有没有这个policy
	var count int64
	e.db.Model(&model.UserRole{}).Where("type = ? AND v1 = ?", POLICY, policy).Count(&count)
	if count == 0 {
		return NoMatchingPolicy
	}

	e.db.Create(&model.UserRole{
		ID:   tools.CreateID(),
		Type: GROUP,
		V1:   uid,
		V2:   policy,
	})

	return nil
}

func (e *Enforcer) AddPolicy(policy, source, action string) {
	e.db.Create(&model.UserRole{
		ID:   tools.CreateID(),
		Type: POLICY,
		V1:   policy,
		V2:   source,
		V3:   action,
	})
}

func (e *Enforcer) Enforce(uid, source, action string) error {
	var role model.UserRole
	if err := e.db.Where("type = ? AND v2 = ? AND v3 = ?", POLICY, source, action).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return NoMatchingPolicy
		} else {
			return OtherError
		}
	}

	var count int64
	e.db.Model(&model.UserRole{}).Where("type = ? AND v1 = ? AND v2 = ?", GROUP, uid, role.V1).Count(&count)
	if count == 0 {
		return InsufficientPermissions
	}

	return nil
}

func (e *Enforcer) RemovePolicy(policy, source, action string) error {
	if err := e.db.Where("type = ? AND v1 = ? AND v2 = ? AND v3 = ?", POLICY, policy, source, action).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	return nil
}

func (e *Enforcer) RemoveGroup(uid, policy string) error {
	if err := e.db.Where("type = ? AND v1 = ? AND v2 = ?", GROUP, uid, policy).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	return nil
}
