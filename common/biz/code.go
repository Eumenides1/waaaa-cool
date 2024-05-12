package biz

import (
	"errors"
	"framework/waError"
)

const OK = 0

var (
	Fail                        = waError.NewError(1, errors.New("请求失败"))
	RequestDataError            = waError.NewError(2, errors.New("请求数据错误"))
	SqlError                    = waError.NewError(3, errors.New("数据库操作错误"))
	InvalidUsers                = waError.NewError(4, errors.New("无效用户"))
	PermissionNotEnough         = waError.NewError(6, errors.New("权限不足"))
	SmsCodeError                = waError.NewError(7, errors.New("短信验证码错误"))
	ImgCodeError                = waError.NewError(8, errors.New("图形验证码错误")) // 图形验证码错误
	SmsSendFailed               = waError.NewError(9, errors.New("短信发送失败"))
	ServerMaintenance           = waError.NewError(10, errors.New("服务器维护"))
	NotEnoughGold               = waError.NewError(11, errors.New("钻石不足"))
	UserDataLocked              = waError.NewError(12, errors.New("用户数据被锁定"))
	NotEnoughScore              = waError.NewError(13, errors.New("积分不足"))
	AccountOrPasswordError      = waError.NewError(101, errors.New("账号或密码错误"))
	GetHallServersFail          = waError.NewError(102, errors.New("获取大厅服务器失败"))
	AccountExist                = waError.NewError(103, errors.New("账号已存在"))
	AccountNotExist             = waError.NewError(104, errors.New("帐号不存在"))
	NotFindBindPhone            = waError.NewError(105, errors.New("该手机号未绑定"))
	PhoneAlreadyBind            = waError.NewError(106, errors.New("该手机号已被绑定，无法重复绑定"))
	NotFindUser                 = waError.NewError(107, errors.New("用户不存在"))
	TokenInfoError              = waError.NewError(201, errors.New("无效的token"))
	NotEnoughVipLevel           = waError.NewError(202, errors.New("vip等级不足"))
	BlockedAccount              = waError.NewError(203, errors.New("帐号已冻结"))
	AlreadyCreatedUnion         = waError.NewError(204, errors.New("已经创建过牌友圈，无法重复创建"))
	UnionNotExist               = waError.NewError(205, errors.New("联盟不存在"))
	UserInRoomDataLocked        = waError.NewError(206, errors.New("用户在房间中，无法操作数据"))
	NotInUnion                  = waError.NewError(207, errors.New("用户不在联盟中"))
	AlreadyInUnion              = waError.NewError(208, errors.New("用户已经在联盟中"))
	InviteIdError               = waError.NewError(209, errors.New("邀请码错误"))
	NotYourMember               = waError.NewError(210, errors.New("添加的用户不是你的下级成员"))
	ForbidGiveScore             = waError.NewError(211, errors.New("禁止赠送积分"))
	ForbidInviteScore           = waError.NewError(212, errors.New("禁止玩家或代理邀请玩家"))
	CanNotCreateNewHongBao      = waError.NewError(213, errors.New("暂时无法分发新的红包"))
	CanNotLeaveRoom             = waError.NewError(305, errors.New("正在游戏中无法离开房间"))
	RoomCountReachLimit         = waError.NewError(301, errors.New("房间数量到达上线"))
	LeaveRoomGoldNotEnoughLimit = waError.NewError(302, errors.New("金币不足，无法开始游戏"))
	LeaveRoomGoldExceedLimit    = waError.NewError(303, errors.New("金币超过最大限度，无法开始游戏"))
	NotInRoom                   = waError.NewError(306, errors.New("不在该房间中"))
	RoomPlayerCountFull         = waError.NewError(307, errors.New("房间玩家已满"))
	RoomNotExist                = waError.NewError(308, errors.New("房间不存在"))
	CanNotEnterNotLocation      = waError.NewError(309, errors.New("无法进入房间，获取定位信息失败"))
	CanNotEnterTooNear          = waError.NewError(310, errors.New("无法进入房间，与房间中的其他玩家太近"))
)
