<template>
  <div class="h-full flex justify-center items-center">
    <div class="w-full bg-purple-50 bg-opacity-50" style="height: 34.125rem">
      <div class="container m-auto flex h-full justify-center">
        <div class="h-full mr-8">
          <img
            src="@/assets/img/signin_bg.png"
            alt=""
            class="w-full h-full"
            draggable="false"
          />
        </div>

        <div
          class="bg-white rounded p-16 my-12 w-1/3 flex flex-col justify-center"
        >
          <img src="@/assets/img/logo.svg" draggable="false" alt="" />
          <h2 class="font-bold mb-8 text-center text-lg">在线调解平台</h2>
          <el-form
            :model="userInfo"
            :rules="rules"
            ref="ruleForm"
            label-width="0px"
            class=""
          >
            <el-form-item prop="userName">
              <el-autocomplete
                v-model.trim="userInfo.userName"
                :fetch-suggestions="querySearch"
                class="synergy_input w-full"
                placeholder="请输入昵称"
                @select="handleSelect"
              />
            </el-form-item>
            <el-form-item class="mb-0">
              <el-button class="w-full" type="primary" @click="LoginSynergy()">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
// 权限
import { ProjectPermissions } from "@/api/api_config.js";
// 生成随机数
import { randomString } from "@/assets/untils/until.js";
// 登录接口
import { signIn } from "@/api/singin.js";
// 路由跳转
import { useRouter } from "vue-router";
import { defineComponent, ref, onMounted } from "vue";
export default defineComponent({
  setup() {
    // 路由
    const oRoute = useRouter();

    // 用户昵称校验规则
    const validateuserName = (rule, value, callback) => {
      if (value === "") {
        callback(new Error("请输入昵称"));
      } else if (value.length > 16) {
        callback(new Error("昵称最大为16个字符"));
      } else {
        callback();
      }
    };
    // 用户昵称校验(页面使用)
    const rules = ref({
      userName: [{ validator: validateuserName, trigger: "change" }],
    });

    // 本地记录(页面使用)
    const LocalRecordsUser = ref([]);
    // 根据输入内容提供对应的输入建议(页面使用)
    const querySearch = (queryString, cb) => {
      const results = queryString
        ? LocalRecordsUser.value.filter(createFilter(queryString))
        : LocalRecordsUser.value;
      cb(results);
    };
    // 查找本地存储的类似用户信息
    const createFilter = (queryString) => {
      return (LocalRecords) => {
        return (
          LocalRecords.value
            .toLowerCase()
            .indexOf(queryString.toLowerCase()) === 0
        );
      };
    };
    // 获取本地记录的用户登录信息
    const getLocalRecordsUser = () => {
      // 获取本地存储昵称
      const oLocalUser = localStorage.getItem("Synergy_LocalHistory");
      if (oLocalUser) {
        return JSON.parse(oLocalUser);
      } else {
        return [];
      }
    };
    onMounted(() => {
      LocalRecordsUser.value = getLocalRecordsUser();
    });
    const handleSelect = (item) => {
      userInfo.value = item;
    };

    // 用户登录信息
    const userInfo = ref({
      uid: "",
      userName: "",
      userType: ProjectPermissions,
      workName: "",
      value: "",
    });
    //  登录点击校验
    const ruleForm = ref(null);
    // 本地记录当前登录信息
    const LoginSynergy = () => {
      ruleForm.value.validate((valid) => {
        if (valid) {
          LocalRecordsUser.value = getLocalRecordsUser();
          // 不存在本地记录的用户登录信息
          if (LocalRecordsUser.value.length > 0) {
            // 存在本地记录的用户登录信息
            const oM = LocalRecordsUser.value.filter((item) => {
              return item.userName === userInfo.value.userName;
            });
            if (oM.length > 0) {
              // 选中本地记录
              userInfo.value = oM[0];
            } else {
              // 设置
              userInfo.value = Object.assign(userInfo.value, {
                uid: randomString(9, "number"),
                value: userInfo.value.userName,
              });
              LocalRecordsUser.value.push(userInfo.value);
              localStorage.setItem(
                "Synergy_LocalHistory",
                JSON.stringify(LocalRecordsUser.value)
              );
            }
          } else {
            // 设置
            userInfo.value = Object.assign(userInfo.value, {
              uid: randomString(9, "number"),
              value: userInfo.value.userName,
            });
            // 第一次用户
            localStorage.setItem(
              "Synergy_LocalHistory",
              JSON.stringify([userInfo.value])
            );
          }
          // 登录
          Singin();
        } else {
          console.log("error submit!!");
          // this.singinShow = false;
          return false;
        }
      });
    };
    const Singin = async () => {
      const { code, data } = await signIn(userInfo.value);
      if (code == 0) {
        // 设置用户信息
        sessionStorage.setItem("SynergyLogin_UserInfo", JSON.stringify(data));
        // 跳转至首页
        oRoute.replace("/");
      }
    };
    return {
      ruleForm, // 登录点击校验(ref)
      rules, // 用户昵称校验(数据)
      userInfo, // 用户登录信息(数据)
      LocalRecordsUser, // 本地记录用户登录信息(数据)
      LoginSynergy, // 本地记录当前登录信息(方法)
      Singin, // 登录
      querySearch, // 根据输入内容提供对应的输入建议(方法)
      handleSelect,
    };
  },
});
</script>
<style lang="scss">
.synergy_input {
  .el-input__inner {
    @apply border-t-0 border-r-0 border-l-0 rounded-none;
  }
}
</style>
