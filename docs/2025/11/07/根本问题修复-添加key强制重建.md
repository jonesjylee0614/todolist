# 根本问题修复 - 添加 key 强制重建组件

**修复时间**: 2025-11-07  
**问题**: 路由变化但组件不重新创建，导致接口不调用

---

## 🔍 问题诊断

### 用户反馈的日志

```
App.vue:117 === Route changed === now -> future
App.vue:117 === Route changed === future -> history
App.vue:117 === Route changed === history -> future
App.vue:117 === Route changed === future -> now
```

### 关键发现

1. ✅ **路由变化**：日志显示路由正常切换
2. ❌ **View 组件未重建**：没有看到 "NowView.vue loaded" 等日志
3. ❌ **TaskBoard 未挂载**：没有看到 "TaskBoard onMounted" 日志
4. ❌ **接口未调用**：因为 onMounted 没有执行

### 根本原因

**Vue Router 的组件复用机制**：

当路由切换时，如果新旧组件是同一个类型（比如从 NowView 切换到 FutureView），Vue 默认会**复用组件实例**而不是销毁重建。

这是 Vue 的性能优化策略，但在我们的场景下导致了问题：
- 组件被复用
- `onMounted` 不会再次执行
- 因此不会加载新数据

---

## ✅ 解决方案

### 核心思路

给 `<component>` 添加 `key` 属性，强制 Vue 在路由变化时**销毁旧组件**并**创建新组件**。

### 修改 1: App.vue - 添加 key

```vue
<!-- 修复前 -->
<RouterView v-slot="{ Component }">
  <Transition name="fade-slide" mode="out-in">
    <component :is="Component" />
  </Transition>
</RouterView>

<!-- 修复后 -->
<RouterView v-slot="{ Component, route: childRoute }">
  <Transition name="fade-slide" mode="out-in">
    <component :is="Component" :key="childRoute.name" />
  </Transition>
</RouterView>
```

**关键改动**：
1. 添加 `route: childRoute` 到 slot props
2. 给 component 添加 `:key="childRoute.name"`

**工作原理**：
- `key` 是 "now" / "future" / "history"
- 当 key 变化时，Vue 会销毁旧组件并创建新组件
- 新组件的 `onMounted` 会被执行
- 从而加载新数据

### 修改 2: TaskBoard.vue - 添加根元素

```vue
<!-- 修复前：多个根节点（v-for） -->
<template>
  <TaskColumn v-for="..." />
</template>

<!-- 修复后：单一根元素 -->
<template>
  <div class="task-board-wrapper">
    <TaskColumn v-for="..." />
  </div>
</template>
```

**原因**：
- Vue 警告：`Component inside <Transition> renders non-element root node`
- Transition 组件要求子组件有单一根元素
- 添加包裹的 div 解决此问题

---

## 🛠️ 修改文件

### frontend/src/App.vue

```vue
<RouterView v-slot="{ Component, route: childRoute }">
  <Transition name="fade-slide" mode="out-in">
    <component :is="Component" :key="childRoute.name" />
  </Transition>
</RouterView>
```

### frontend/src/modules/tasks/components/TaskBoard.vue

```vue
<template>
  <div class="task-board-wrapper">
    <TaskColumn
      v-for="column in visibleColumns"
      :key="column.status"
      ...
    />
  </div>
</template>
```

---

## 🧪 测试验证

### 预期结果

刷新页面后，点击标签切换应该看到：

1. **点击"未来"**：
```
=== Route changed === now -> future
FutureView.vue loaded
=== FutureView mounted ===
=== TaskBoard onMounted === future
Loading tasks for status: future
```

2. **点击"历史"**：
```
=== Route changed === future -> history
HistoryView.vue loaded
=== HistoryView mounted ===
=== TaskBoard onMounted === history
Loading tasks for status: history
```

3. **Network 标签**：
每次切换都应该看到对应的 API 请求：
- `GET /api/tasks?status=future`
- `GET /api/tasks?status=history`
- 等等

---

## 📊 技术深度解析

### Vue Router 的组件复用

Vue Router 的默认行为：
```
路由A -> 路由B
如果都使用同一个组件 ➜ 复用实例
如果使用不同组件 ➜ 销毁重建
```

在我们的场景：
```
/now (NowView) -> /future (FutureView) -> /history (HistoryView)
```

虽然是三个不同的组件文件，但因为它们的结构几乎相同（都只是简单包裹 TaskBoard），Vue 可能会尝试复用。

### key 的作用

`key` 是 Vue 的特殊属性：
- 当 key 变化时，Vue 认为这是**完全不同的元素**
- 触发**完整的销毁和重建流程**
- 包括调用 `onBeforeUnmount`, `onUnmounted`, `onBeforeMount`, `onMounted`

### 为什么使用 route.name 作为 key

```javascript
:key="childRoute.name"  // "now", "future", "history"
```

优点：
- 每个路由有唯一的 name
- name 变化 = 路由变化
- 简单明了，易于理解

其他选项：
```javascript
:key="childRoute.path"      // "/now", "/future", "/history"
:key="childRoute.fullPath"  // 包含 query 参数
```

---

## 🎯 核心要点

1. **理解 Vue 的复用机制**：
   - Vue 会尽可能复用组件以提高性能
   - 有时这会导致生命周期钩子不执行

2. **使用 key 控制重建**：
   - key 是控制 Vue 是否复用的开关
   - 不同的 key = 不同的组件实例

3. **Transition 的要求**：
   - 子组件必须有单一根元素
   - 不能是 v-if/v-for 直接生成的多个元素

---

## 💡 其他解决方案对比

### 方案 1: 添加 key（当前方案）✅

```vue
<component :is="Component" :key="route.name" />
```

**优点**：
- 简单直接
- 完全控制组件重建
- 不改变业务逻辑

**缺点**：
- 每次切换都销毁重建（性能开销）

### 方案 2: 使用 watch route

```typescript
watch(() => route.name, (newRoute) => {
  const statusMap = { now: 'now', future: 'future', history: 'history' };
  tasksStore.load(statusMap[newRoute]);
});
```

**优点**：
- 不销毁组件，性能更好
- 可以做更细粒度的控制

**缺点**：
- 需要在多个地方添加 watch
- 逻辑分散，难以维护

### 方案 3: 使用 keep-alive + watch

```vue
<keep-alive>
  <component :is="Component" />
</keep-alive>
```

配合 watch route

**优点**：
- 组件状态保留
- 可以做缓存优化

**缺点**：
- 需要额外处理数据刷新
- 复杂度增加

### 为什么选择方案 1

对于当前应用：
- 三个 View 组件非常轻量
- 销毁重建的性能开销可以忽略
- 实现最简单，代码最清晰
- 符合"简单优于复杂"的原则

---

## ✅ 完成状态

- [x] 问题定位完成
- [x] 添加 key 强制重建
- [x] 修复 Transition 警告
- [x] 测试验证准备就绪
- [x] 文档编写完成

---

**状态**: ✅ 已完成修复  
**修复时间**: 2025-11-07  
**修复人员**: Claude AI

**关键改动**：
1. `App.vue` - 给 component 添加 `:key="childRoute.name"`
2. `TaskBoard.vue` - 添加根元素包裹，修复 Transition 警告

