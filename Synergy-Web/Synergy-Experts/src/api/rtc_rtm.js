import { post } from "./http";

// 加入房间
export const joinRoom = (data) => post("/teamview/joinRoom", data);

// 用户离开房间
export const leaveRoom = (data) => post("/teamview/leaveRoom", data);
