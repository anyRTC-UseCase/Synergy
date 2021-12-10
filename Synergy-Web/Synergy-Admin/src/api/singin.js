import { post } from "./http";
// 登录
export const signIn = (data) => post("/teamview/signIn", data);

// 退出
