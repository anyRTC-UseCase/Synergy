<template>
  <!-- 视频在线播放 -->
  <div class="h-full w-full relative">
    <!-- 页面关闭 -->
    <div
      class="h-16 z-50 text-white left-0 bg-black bg-opacity-50 absolute top-0 w-full flex justify-between px-12 items-center"
    >
      <span>{{ roomName }}</span>
      <div
        title="关闭"
        class=""
        @click="
          $router.push({
            name: 'Home',
            params: {
              type: 1,
            },
          })
        "
      >
        <img src="@/assets/img/close.svg" alt="" />
      </div>
      <!-- <a
        title="关闭"
        href="javascript:window.opener=null;window.open('','_self');window.close();window.open('', '_self').window.close();"
      >
        <img src="@/assets/img/close.svg" alt=""
      /></a> -->
    </div>
    <video-player
      class="h-full w-full set_video-player"
      ref="videoPlayer"
      :playsinline="false"
      :options="playerOptions"
      @timeupdate="onPlayerTimeupdate($event)"
      @ready="playerReadied"
    >
    </video-player>
    <!-- 操作 -->
    <div
      class="absolute left-1/2 bottom-20 text-white bg-black bg-opacity-50 flex flex-col items-center rounded py-4 px-16 -ml-80"
    >
      <div class="flex">
        <!-- 倒退十秒 -->
        <img
          @click="videoSpeed(false)"
          class="cursor-pointer"
          src="@/assets/img/video_retrogress.svg"
          alt=""
          title="倒退十秒"
        />
        <!-- 开始 -->
        <img
          class="cursor-pointer mx-6"
          v-show="!play_pause"
          @click="playPauseFn"
          src="@/assets/img/video_play.svg"
          alt=""
        />
        <!-- 暂停 -->
        <img
          class="cursor-pointer mx-6"
          v-show="play_pause"
          @click="playPauseFn"
          src="@/assets/img/video_pause.svg"
          alt=""
        />
        <!-- 快进十秒 -->
        <img
          @click="videoSpeed(true)"
          class="cursor-pointer"
          src="@/assets/img/video_speed.svg"
          title="快进十秒"
          alt=""
        />
      </div>
      <div class="mt-6 flex items-center">
        <!-- 当前播放时间 -->
        <span>{{ video_current_time }}</span>
        <el-slider
          class="w-96 mx-6 video_slider"
          :min="0"
          :max="100"
          v-model="video_slider"
          :show-tooltip="false"
          @change="videoSlider"
        ></el-slider>
        <!-- 总时间 -->
        <span>{{ video_total_time }}</span>
      </div>
    </div>
  </div>
</template>

<script>
// 时间转化
import { TimeTransitionFormatOne } from "@/assets/untils/time.js";
export default {
  data() {
    return {
      roomName: "",
      // 播放暂停
      play_pause: false,
      // 滑动进度
      video_slider: 0,
      // 视频当前进度
      video_current_time: "00:00:00",
      // 视频总时长
      video_total_time: "00:00:00",
      playerOptions: {},
    };
  },
  created() {
    const oData = this.$route.query;
    this.roomName = oData.roomName;
    this.$nextTick(() => {
      this.playerOptions = {
        controls: false,
        //如果true,浏览器准备好时开始回放。
        autoplay: false,
        // 默认情况下将会消除任何音频。
        muted: false,
        // 导致视频一结束就重新开始。
        loop: false,
        // 建议浏览器在<video>加载元素后是否应该开始下载视频数据。auto浏览器选择最佳行为,立即开始加载视频（如果浏览器支持）
        preload: "auto",
        language: "zh-CN",
        // 将播放器置于流畅模式，并在计算播放器的动态大小时使用该值。值应该代表一个比例 - 用冒号分隔的两个数字（例如"16:9"或"4:3"）
        // aspectRatio: "4:3",
        // fluid: true,
        sources: [
          {
            type: "video/mp4",
            // type: "video/ogg",
            // src: "https://media.w3.org/2010/05/sintel/trailer.mp4",
            src: oData.url,
          },
          // {
          //   type: "application/x-mpegURL", // 类型
          // //   // src: "https://teameeting.oss-cn-shanghai.aliyuncs.com/ar/anyrtc/5BBGaDJIbgO8Z0zEssj0roYi_123456.m3u8"
          //   src: "http://pro.vod.agrtc.cn/ar/anyrtc/f4G95r8iJGc8j1I7_12345.mp4",
          // //   // src: "http://121.37.135.206/recorder/h5/dzahm3oaH39ujkqgStoIk0ZHHv89mFexh5Yj6AxCU37hqruQ/v0YDN40Xb1uqNdM9005KtTY1_55667788.m3u8",
          // },
        ],
        notSupportedMessage: "此视频暂无法播放，请稍后再试",
      };
    });
  },
  methods: {
    // 快进、倒退
    videoSpeed(fase) {
      // 当前时长
      let oTimeIng = this.$refs.videoPlayer.player.currentTime();
      if (fase) {
        if (oTimeIng + 10 < this.$refs.videoPlayer.player.duration()) {
          //快进
          this.$refs.videoPlayer.player.currentTime(oTimeIng + 10);
        }
      } else {
        if (oTimeIng - 10 > 0) {
          // 倒退
          this.$refs.videoPlayer.player.currentTime(oTimeIng - 10);
        }
      }
    },
    // 播放、暂停
    playPauseFn() {
      this.play_pause = !this.play_pause;
      this.video_total_time = TimeTransitionFormatOne(
        this.$refs.videoPlayer.player.duration()
      );
      if (this.play_pause) {
        this.$refs.videoPlayer.player.play(); // 播放
      } else {
        this.$refs.videoPlayer.player.pause(); // 暂停
      }
    },
    // 进度(设置)
    videoSlider(val) {
      this.video_slider = val;
      // 总时长
      let oAllTime = this.$refs.videoPlayer.player.duration();
      this.$refs.videoPlayer.player.currentTime((val * oAllTime) / 100);
    },
    // 监听播放进度
    onPlayerTimeupdate(e) {
      this.video_current_time = TimeTransitionFormatOne(e.currentTime());
      this.video_slider = (e.currentTime() * 100) / e.duration();
      if (e.currentTime() == e.duration()) {
        this.play_pause = false;
      }
    },
    // 初始化
    playerReadied() {
      const oA = setInterval(async () => {
        const oTime = await TimeTransitionFormatOne(
          this.$refs.videoPlayer.player.duration()
        );
        if (oTime !== "NaN:NaN") {
          this.video_total_time = oTime;
          clearInterval(oA);
        } else {
          clearInterval(oA);
        }
      }, 500);
    },
  },
};
</script>

<style lang="scss">
.set_video-player {
  .video-js {
    @apply h-full w-full;
    .vjs-big-play-button {
      display: none !important;
    }
    // .vjs-control-bar {
    //   @apply h-36;
    // }
  }
}

// 进度条修改
.video_slider {
  .el-slider__runway {
    @apply bg-synergy-gray_400;
    .el-slider__bar {
      @apply bg-synergy-gray_200;
    }
    .el-slider__button-wrapper {
      .el-slider__button {
        @apply w-2 h-5 rounded border-0;
      }
    }
  }
}
</style>
