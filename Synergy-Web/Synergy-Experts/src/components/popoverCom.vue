<template>
  <div>
    <el-popover
      align="center"
      ref="countPopover"
      placement="right"
      width="200"
      trigger="click"
    >
      <span>{{ scope.row[propsName] }}</span>
      <!-- 调用popover内部的reference插槽 -->
      <template #reference>
        <span v-popover:countPopover>{{ scope.row[propsName] }}</span>
      </template>
    </el-popover>
  </div>
</template>

<script>
/* eslint-disable */

export default {
  name: "tableAndPopover",
  data() {
    return {};
  },
  props: {
    scope: Object,
    propsName: String,
  },
  methods: {},
};
</script>
<style scoped></style>
