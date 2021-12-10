/**
 * 时间格式转换
 * 秒 转换 00:00:00
 * */
export const TimeTransitionFormatOne = (time) => {
  let b = "";
  let h = parseInt(time / 3600);
  let m = parseInt((time % 3600) / 60);
  let s = parseInt((time % 3600) % 60);
  if (h > 0) {
    h = h < 10 ? "0" + h : h;
    b += h + ":";
  }
  m = m < 10 ? "0" + m : m;
  s = s < 10 ? "0" + s : s;
  b += m + ":" + s;
  return b;
};

/**
 * 时间格式转换
 * 秒 转换 年/月/日 00:00:00
 */
export const TimeTransitionFormatTwo = (time, dateSplitStr = "/") => {
  time = String(time).length == 10 ? Number(time) * 1000 : Number(time);
  let oDatetimer = new Date(time);
  let oMonth = oDatetimer.getMonth() + 1;
  let oDay = oDatetimer.getDate();
  let oH = oDatetimer.getHours();
  let oM = oDatetimer.getMinutes();
  let oS = oDatetimer.getSeconds();
  oMonth = oMonth < 10 ? "0" + oMonth : oMonth;
  oDay = oDay < 10 ? "0" + oDay : oDay;
  oH = oH < 10 ? "0" + oH : oH;
  oM = oM < 10 ? "0" + oM : oM;
  oS = oS < 10 ? "0" + oS : oS;
  return (
    oDatetimer.getFullYear() +
    dateSplitStr +
    oMonth +
    dateSplitStr +
    oDay +
    " " +
    oH +
    ":" +
    oM +
    ":" +
    oS
  );
};
