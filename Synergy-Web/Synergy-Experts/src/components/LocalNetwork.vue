<template>
  <!-- 本地网络断开 -->
  <div
    class="fixed top-0 left-0 w-full h-full bg-black bg-opacity-50 flex justify-center items-center z-50"
  >
    <div
      class="flex flex-col justify-center items-center bg-white w-40 py-8 rounded animate-pulse select-none"
    >
      <img src="@/assets/img/local_network.svg" class="mb-2" alt="" />
      <span> 连网中...</span>
    </div>
  </div>
</template>

<script>
export default {};
</script>

<style></style>
