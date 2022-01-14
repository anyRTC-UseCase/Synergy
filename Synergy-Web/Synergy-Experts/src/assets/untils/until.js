// 生成随机数字
export const randomString = (n, type = "str") => {
  let str = "abcdefghijklmnopqrstuvwxyz9876543210";
  if (type == "number") {
    str = "9876543210";
  }
  let tmp = "";
  let l = str.length;
  let i = 0;
  for (i = 0; i < n; i++) {
    tmp += str.charAt(Math.floor(Math.random() * l));
  }
  return tmp;
};

// 节流
export const throttle = (fn, delay = 100) => {
  let timer = null;
  let ctx = this;
  return function () {
    if (timer) {
      return;
    }
    timer = setTimeout(() => {
      fn.apply(ctx, arguments);
      timer = null;
    }, delay);
  };
};

// 容器固定比例 w:h=4:3
export const containerRatio = (id, applyid) => {
  const oId = document.getElementById(id);

  const oApplyId = document.getElementById(applyid);
  const iHeight = oId.offsetHeight;
  const iWidth = oId.offsetWidth;
  // 宽度 > 高度
  if (iWidth > (4 * iHeight) / 3) {
    oApplyId.style.width = (4 * iHeight) / 3 + "px";
    oApplyId.style.height = iHeight + "px";
  } else {
    oApplyId.style.width = iWidth + "px";
    oApplyId.style.height = (3 * iWidth) / 4 + "px";
  }
};
