import { post } from "./http";

// 记录用户在线心跳包信息
export const insertUserOnlineInfo = (data) =>
  post("/teamview/insertUserOnlineInfo", data);

// 获取用户信息
export const getUserInfo = (data) => post("/teamview/getUserInfo", data);

// 获取房间列表
export const getRoomList = (data) => post("/teamview/getRoomList", data);

// 获取专家列表
export const getSpecialist = (data) => post("/teamview/getSpecialist", data);
