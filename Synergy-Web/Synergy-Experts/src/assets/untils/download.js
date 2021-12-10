// const oAxios = (url) => {
//   const xhr = new XMLHttpRequest();
//   xhr.open("GET", url, true);
//   xhr.responseType = "blob";
//   //xhr.setRequestHeader('Authorization', 'Basic a2VybWl0Omtlcm1pdA==');
//   xhr.onload = () => {
//     if (xhr.status === 200) {
//       console.log(xhr.response);
//       // // 获取文件blob数据并保存
//       // var fileName = getFileName(url);
//       // saveAs(xhr.response, fileName);
//     }
//   };

//   xhr.send();
// };

// url 地址下载（非同源）
export const downloadUrl = async (videoUrl, downloadName) => {
  // await oAxios(videoUrl);
  // console.log(downloadName);
  const res = await fetch(videoUrl);
  console.log(res);
  const blob = await res.blob();
  const a = document.createElement("a");
  document.body.appendChild(a);
  a.style.display = "none";
  const url = window.URL.createObjectURL(blob);
  a.href = url;
  a.download = downloadName;
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
};

// 同原
export const downloadUrl2 = async (videoUrl) => {
  const ele = document.createElement("a");
  ele.setAttribute("href", videoUrl); //设置下载文件的url地址
  ele.setAttribute("download", "download"); //用于设置下载文件的文件名
  ele.click();
};
