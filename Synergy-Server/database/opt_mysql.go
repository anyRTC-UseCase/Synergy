package database

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"ARTeamViewService/utils"
	"errors"
	"time"
)

/**
 * 通过uid获取用户信息
 */
func QueryUserInfoByUId(uid string) (tableData []map[string]interface{}, err error) {
	var strSql = "select uid, u_name as userName, u_workname as workName, u_type as userType, u_ts as userTs from user_info where uid = ?"
	rows, err := global.GDb.Query(strSql, uid)

	if err != nil {
		global.GLogger.Error("QueryUserInfoByUId, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

/**
 * 通过uid,userName,userType获取用户信息
 */
func QueryUserInfoByUserInfo(uid, userName string, userType int, workName string) (tableData []map[string]interface{}, err error) {
	var strSql = "select uid, u_name as userName, u_workname as workName, u_type as userType, u_ts as userTs " +
		"from user_info where u_name = ? and u_type = ? and u_workname = ? or uid = ? "
	rows, err := global.GDb.Query(strSql, userName, userType, workName, uid)

	if err != nil {
		global.GLogger.Error("QueryUserInfoByUserInfo, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

//插入用户信息
func InsertUserInfo(userName, uid, workName string, userType int) (cnt int, err error) {
	var strSql = "insert into user_info (uid, u_name, u_workname, u_ts, u_type) values(?,?,?,?,?)"
	stmt, err := global.GDb.Prepare(strSql)
	defer stmt.Close()
	if err != nil {
		global.GLogger.Error("InsertUserInfo, 1, ", err)
		return 0, err
	}

	result, err := stmt.Exec(uid, userName, workName, time.Now().Unix(), userType)
	if err != nil {
		global.GLogger.Error("InsertUserInfo, 2, ", err)
		return 0, err
	}

	affectRows, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("InsertUserInfo, 3, ", err)
		return 0, err
	}
	return int(affectRows), nil
}

/**
 * 查询房间数量
 * @param roomid
 */
func QueryRoomCountByRoomId(uid string) (cnt int, err error) {
	var strSql = "select count(roomid) as cnt from room_info where roomid = ?"
	err = global.GDb.QueryRow(strSql, uid).Scan(&cnt)
	if err != nil {
		global.GLogger.Error("QueryRoomCountByRoomId, 1, ", err)
		return cnt, err
	}
	return cnt, nil
}

/**
 * 创建房间
 */
func InsertRoomInfo(roomId, roomName, uid string, roomState int) (cnt int, err error) {
	var strSql = "insert into room_info (roomid, r_name, r_hostid , r_state, r_ts) values(?,?,?,?,?)"
	stmt, err := global.GDb.Prepare(strSql)
	defer stmt.Close()
	if err != nil {
		global.GLogger.Error("InsertRoomInfo, 1, ", err)
		return 0, err
	}

	result, err := stmt.Exec(roomId, roomName, uid, roomState, time.Now().Unix())
	if err != nil {
		global.GLogger.Error("InsertRoomInfo, 2, ", err)
		return 0, err
	}

	affectRows, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("InsertRoomInfo, 3, ", err)
		return 0, err
	}
	return int(affectRows), nil
}

/**
 * 房间信息
 */
func QueryRoomInfoById(roomId string) (tableData []map[string]interface{}, err error) {
	var strSql = "select roomid as roomId, " +
		"r_name as roomName, " +
		"r_hostid as roomHostId, " +
		"r_state as roomState, " +
		"r_vod_uid as roomVodUId, " +
		"r_vod_resource_id as roomVodResourceId, " +
		"r_vod_sid as roomVodSId, " +
		"r_ts as roomTs " +
		"from room_info " +
		"where roomid = ? "

	rows, err := global.GDb.Query(strSql, roomId)

	if err != nil {
		global.GLogger.Error("QueryRoomInfoById, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

/**
 * 加入房间
 */
func InsertUserJoinInfo(roomId, uid string, userRole int) (cnt int, err error) {
	var strSql = "insert into user_join_info (ar_uid, ar_roomid, ar_join_time, ar_ts, ar_user_role) values(?,?,?,?,?)"
	stmt, err := global.GDb.Prepare(strSql)
	defer stmt.Close()
	if err != nil {
		global.GLogger.Error("InsertUserJoinInfo, 1, ", err)
		return 0, err
	}

	result, err := stmt.Exec(uid, roomId, time.Now().Unix(), time.Now().Unix(), userRole)
	if err != nil {
		global.GLogger.Error("InsertUserJoinInfo, 2, ", err)
		return 0, err
	}

	affectRows, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("InsertUserJoinInfo, 3, ", err)
		return 0, err
	}
	return int(affectRows), nil
}

/**
 * 更新房间录像信息
 */
func UpdateRoomInfo(roomId, resourceId, sid, vodUId string) (cnt int, err error) {

	var strSql = "update room_info set r_vod_resource_id = ?, r_vod_sid = ?, r_vod_uid = ? where roomid = ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomInfo 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(resourceId, sid, vodUId, roomId)
	if err != nil {
		global.GLogger.Error("UpdateRoomInfo 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomInfo 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新房间录像文件url
 */
func UpdateRoomVodUrl(roomId, fileUrl string) (cnt int, err error) {

	var strSql = "update room_info set r_vod_file_url = ? where roomid = ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomVodUrl 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(fileUrl, roomId)
	if err != nil {
		global.GLogger.Error("UpdateRoomVodUrl 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomVodUrl 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新房间录制结束时间
 */
func UpdateRoomStopTs(roomId string, stopTs int64) (cnt int, err error) {

	var strSql = "update room_info set r_vod_stop_ts = ? where roomid = ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomStopTs 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(stopTs, roomId)
	if err != nil {
		global.GLogger.Error("UpdateRoomStopTs 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomStopTs 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新房间录像开始时间
 */
func UpdateRoomStarTs(roomId string, stopTs int64) (cnt int, err error) {

	var strSql = "update room_info set r_vod_start_ts = ? where roomid = ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomStarTs 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(stopTs, roomId)
	if err != nil {
		global.GLogger.Error("UpdateRoomStarTs 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomStarTs 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

// 查询房间人数
func QueryRoomUserNum(roomId string) (cnt int, err error) {

	var strSql = "select count(id) as cnt from user_join_info where ar_roomid = ?"
	err = global.GDb.QueryRow(strSql, roomId).Scan(&cnt)
	if err != nil {
		global.GLogger.Error("QueryRoomUserNum, 1, ", err)
		return cnt, err
	}
	return cnt, nil
}

/**
 * 更新用户离开房间
 */
func UpdateUserLeaveTs(roomId, uid string) (cnt int, err error) {

	var strSql = "update user_join_info set ar_leave_time = ? where ar_roomid = ? and ar_uid = ? and ar_leave_time = ?"
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTs 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(time.Now().Unix(), roomId, uid, consts.UserLeaveTime)
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTs 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTs 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新除最后一条的leave_time=0的离开时间
 */
func UpdateUserLeaveTsExceptLast(roomId, uid string) (cnt int, err error) {

	var strSql = "update user_join_info set ar_leave_time = ? where ar_roomid = ? and ar_uid = ? and ar_leave_time = ? " +
		"and id <> (select max(id) from user_join_info where ar_roomid = ? and ar_uid = ?)"
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTsExceptLast 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(time.Now().Unix(), roomId, uid, consts.UserLeaveTime, roomId, uid)
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTsExceptLast 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateUserLeaveTsExceptLast 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

//进行中房间列表
func QueryOnGoingRoomList(pageNum, pageSize, roomState int) (tableData []map[string]interface{}, err error) {

	var strSql = "select roomid as roomId, " +
		"r_name as roomName, " +
		"r_hostid as roomHostId, " +
		"(select u_name from user_info where uid = tbRm.r_hostid) as userName, " +
		"r_state as roomState, " +
		"r_ts as roomTs " +
		"from room_info tbRm " +
		"where r_state = ? " +
		"order by r_ts desc " +
		"limit ?, ?"
	rows, err := global.GDb.Query(strSql, roomState, (pageNum-1)*pageSize, pageSize)

	if err != nil {
		global.GLogger.Error("QueryOnGoingRoomList, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

//结束房间列表
func QueryFinishedRoomList(pageNum, pageSize int) (tableData []map[string]interface{}, err error) {

	var strSql = "select roomid as roomId, " +
		"r_name as roomName, " +
		"r_hostid as roomHostId, " +
		"(select u_name from user_info where uid = tbRm.r_hostid) as userName, " +
		"r_state as roomState, " +
		"r_vod_start_ts as roomStartTs, " +
		"r_vod_stop_ts as roomStopTs, " +
		"r_vod_file_url as roomFileUrl, " +
		"r_ts as roomTs " +
		"from room_info tbRm " +
		"where r_state in (?,?) " +
		"order by r_ts desc " +
		"limit ?, ?"
	rows, err := global.GDb.Query(strSql, consts.RoomStateClosed, consts.RoomStateMixed, (pageNum-1)*pageSize, pageSize)

	if err != nil {
		global.GLogger.Error("QueryFinishedRoomList, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

// 查询房间数量
func QueryRoomNum(roomState int) (cnt int, err error) {
	var strSql = ""
	if roomState == consts.RoomStateClosed {
		strSql = "select count(roomid) as cnt from room_info where r_state in (?,?)"
		err = global.GDb.QueryRow(strSql, consts.RoomStateClosed, consts.RoomStateMixed).Scan(&cnt)
	} else {
		strSql = "select count(roomid) as cnt from room_info where r_state = ?"
		err = global.GDb.QueryRow(strSql, roomState).Scan(&cnt)
	}

	if err != nil {
		global.GLogger.Error("QueryRoomNum, 1, ", err)
		return cnt, err
	}
	return cnt, nil
}

//获取专家列表
func QuerySpecialist(pageNum, pageSize int, uid string) (tableData []map[string]interface{}, err error) {
	/*
	 *if 30s没有心跳包:离线
	 *else
	 *最后加入房间的角色为观众或最后加入房间的角色为主播且离开时间>0:空闲
	 *else
	 *最后加入房间的角色为主播且离开时间=0:通话
	 *else
	 *空闲
	 */
	//前30s时间戳
	var tmpTs = utils.FormatNowUnix() - int64(global.GConfig.OnlineInterval)

	//专家状态(0:未知,1:通话中,2:空闲,3:离线)
	//用户角色(1:主播,2:观众)

	strSql := "select uid, " +
		"u_name as userName, " +
		"u_workname as workName, " +
		"u_type as userType, " +
		"u_ts as userTs, " +
		"tbUserJ.ar_roomid as roomId, " +
		"if((select count(id) from user_online_info where ar_uid = tbUser.uid and ar_opt_ts > ?) = ?, ?, " +
		"if(tbUserJ.ar_user_role = ? or tbUserJ.ar_user_role = ? and tbUserJ.ar_leave_time > ?, ?, " +
		"if(tbUserJ.ar_user_role = ? and tbUserJ.ar_leave_time = ?, ?, ?))) as userState " +
		"from user_info tbUser left join " +
		"(select tbUserJTmp1.id, ar_roomid, ar_uid, ar_user_role, ar_leave_time from user_join_info tbUserJTmp1 join " +
		"(select max(id) id from user_join_info group by ar_uid) tbUserJTmp2 on tbUserJTmp1.id = tbUserJTmp2.id) tbUserJ " +
		"on tbUser.uid = tbUserJ.ar_uid " +
		"where u_type = ? and uid <> ? " +
		"order by userState asc, userTs desc " +
		"limit ?, ?"

	rows, err := global.GDb.Query(strSql, tmpTs, consts.IntZero, consts.UserStateOffline, consts.UserRoleAudience,
		consts.UserRoleHost, consts.IntZero, consts.UserStateLeisure, consts.UserRoleHost,
		consts.UserLeaveTime, consts.UserStateBusy, consts.UserStateLeisure,
		consts.UserTypeSpecialist, uid, (pageNum-1)*pageSize, pageSize)

	if err != nil {
		global.GLogger.Error("QuerySpecialist, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}

// 查询专家数量
func QuerySpecialistNum(uid string) (cnt int, err error) {

	var strSql = "select count(uid) as cnt from user_info where u_type = ? and uid <> ? "
	err = global.GDb.QueryRow(strSql, consts.UserTypeSpecialist, uid).Scan(&cnt)
	if err != nil {
		global.GLogger.Error("QuerySpecialistNum, 1, ", err)
		return cnt, err
	}
	return cnt, nil
}

/**
 * 插入用户在线心跳包信息
 */
func InsertUserOnlineInfo(uid string, optTs int, roomId string) (cnt int, err error) {
	var strSql = "insert into user_online_info (ar_uid, ar_opt_ts, ar_ts, ar_roomid) values(?,?,?,?)"
	stmt, err := global.GDb.Prepare(strSql)
	defer stmt.Close()
	if err != nil {
		global.GLogger.Error("InsertUserOnlineInfo, 1, ", err)
		return 0, err
	}

	result, err := stmt.Exec(uid, optTs, time.Now().Unix(), roomId)
	if err != nil {
		global.GLogger.Error("InsertUserOnlineInfo, 2, ", err)
		return 0, err
	}

	affectRows, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("InsertUserOnlineInfo, 3, ", err)
		return 0, err
	}
	return int(affectRows), nil
}

/**
 * 更新房间状态
 */
func UpdateRoomState(roomId string, state int) (cnt int, err error) {

	var strSql = "update room_info set r_state = ? where roomid = ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomState 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(state, roomId)
	if err != nil {
		global.GLogger.Error("UpdateRoomState 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomState 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新房间内所有用户离开时间
 */
func UpdateAllUserLeaveTs(roomId string) (cnt int, err error) {

	var strSql = "update user_join_info set ar_leave_time = ? where ar_roomid = ? and ar_leave_time = ?"
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateAllUserLeaveTs 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(time.Now().Unix(), roomId, consts.UserLeaveTime)
	if err != nil {
		global.GLogger.Error("UpdateAllUserLeaveTs 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateAllUserLeaveTs 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

// 查询用户在线心跳包数量
func QueryUserOnlineNum(uid string) (cnt int, err error) {

	//前30s时间戳
	var tmpTs = utils.FormatNowUnix() - int64(global.GConfig.OnlineInterval*time.Second)

	var strSql = "select count(id) as cnt from user_online_info where ar_uid = ? and ar_opt_ts > ?"
	err = global.GDb.QueryRow(strSql, uid, tmpTs).Scan(&cnt)
	if err != nil {
		global.GLogger.Error("QueryUserOnlineNum, 1, ", err)
		return cnt, err
	}
	return cnt, nil
}

//每天0点50分执行任务删掉前一天的心跳包
func DayDeleteUserOnlineInfo(ts int64) (cnt int, err error) {

	var strSql = "delete from user_online_info where ar_ts < ? "
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("DayDeleteUserOnlineInfo, 1, ", err)
		return -1, errors.New(consts.ErrDbError)
	}
	result, err := stmt.Exec(ts)
	if err != nil {
		global.GLogger.Error("DayDeleteUserOnlineInfo, 2, ", err)
		return -1, errors.New(consts.ErrDbError)
	}
	defer stmt.Close()
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("DayDeleteUserOnlineInfo, 3, ", err)
		return -1, errors.New(consts.ErrDbError)
	}

	return int(rowsAffected), err
}

/**
 * 更新房间状态为转码
 */
func UpdateRoomStateMixed(roomId string) (cnt int, err error) {

	var strSql = "update room_info set r_state = ? where roomid = ? and r_state = ?"
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomStateMixed 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(consts.RoomStateMixed, roomId, consts.RoomStateOpen)
	if err != nil {
		global.GLogger.Error("UpdateRoomStateMixed 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomStateMixed 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

/**
 * 更新房间状态为关闭
 */
func UpdateRoomStateClosed(roomId string) (cnt int, err error) {

	var strSql = "update room_info set r_state = ? where roomid = ? and r_state = ?"
	stmt, err := global.GDb.Prepare(strSql)
	if err != nil {
		global.GLogger.Error("UpdateRoomStateClosed 1, error: ", err)
		return -1, err
	}

	result, err := stmt.Exec(consts.RoomStateClosed, roomId, consts.RoomStateMixed)
	if err != nil {
		global.GLogger.Error("UpdateRoomStateClosed 2, error: ", err)
		return -1, err
	}

	affectRaws, err := result.RowsAffected()
	if err != nil {
		global.GLogger.Error("UpdateRoomStateClosed 3, error: ", err)
		return -1, err
	}
	return int(affectRaws), err
}

//根据roomId和uid查询用户最后一条加入房间的信息
func QueryUserJoinInfo(roomId, uid string) (tableData []map[string]interface{}, err error) {

	strSql := "select ar_roomid as roomId, " +
		"ar_uid as uid, " +
		"ar_user_role as userRole, " +
		"ar_join_time as joinTime, " +
		"ar_leave_time as leaveTime " +
		"from user_join_info  " +
		"where ar_roomid = ? and ar_uid = ? " +
		"order by ar_ts desc " +
		"limit 1"

	rows, err := global.GDb.Query(strSql, roomId, uid)

	if err != nil {
		global.GLogger.Error("QueryUserJoinInfo, 1, ", err)
		return nil, errors.New(consts.ErrDbError)
	}
	defer rows.Close()

	tableData = make([]map[string]interface{}, 0)
	err = utils.RowToMap(rows, &tableData)
	return tableData, err
}
