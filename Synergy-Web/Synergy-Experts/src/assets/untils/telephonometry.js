/**
 * 计时功能换算
 */

let hour = "00"; //时
let minus = "00"; //分
let seconds = "00"; //秒
// 相当于计时(秒)
let x = 0;
// 相当于计时(分)
let a = 0;
// 相当于计时(时)
let b = 0;
// 时间换算
export const beginS = () => {
  //计算秒
  x++;
  if (x < 10) {
    seconds = "0" + x;
  } else if (x >= 10 && x <= 59) {
    seconds = x;
  } else if (x > 59) {
    seconds = "00";
    x = 0;
    a++;
  }

  if (a < 10) {
    minus = "0" + a;
  } else if (a >= 10 && a <= 59) {
    minus = a;
  } else if (a > 59) {
    minus = "00";
    a = 0;
    b++;
  }

  if (b < 10) {
    hour = "0" + b;
  } else if (b >= 10 && b <= 59) {
    hour = b;
  }

  return hour + ":" + minus + ":" + seconds;
};
// 时间清空
export const beginSclear = () => {
  hour = "00"; //时
  minus = "00"; //分
  seconds = "00"; //秒
  // 相当于计时(秒)
  x = 0;
  // 相当于计时(分)
  a = 0;
  // 相当于计时(时)
  b = 0;
};
