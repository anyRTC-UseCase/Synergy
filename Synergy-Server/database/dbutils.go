package database

import (
	"ARTeamViewService/global"
	"ARTeamViewService/utils"
)


/**
 * 生成房间id
 */
func GetRoomId() string {
	tempRoomId := generateRoomId()
	roomCount, err := QueryRoomCountByRoomId(tempRoomId)
	if err != nil {
		global.GLogger.Error(err)
		GetRoomId()
	}

	if roomCount > 0 {
		GetRoomId()
	}
	return tempRoomId
}

/**
 * 生成房间id
 */
func generateRoomId() string {
	return utils.RandomNoZeroNumberString(8)
}
