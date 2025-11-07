# UI 细节优化报告

**优化时间**: 2025-11-07  
**状态**: ✅ 已完成

---

## 📋 优化项目

### 1. ✅ 移除任务 item 鼠标悬停时的拖动条

**问题**：鼠标放到任务卡片上时，出现左右拖动条，影响美观。

**解决方案**：
- 优化滚动条样式，默认隐藏
- 只在容器悬停时显示半透明滚动条
- 减小滚动条宽度从 6px 到 4px

**修改文件**：`frontend/src/styles/index.css`

```css
.scrollbar-thin {
  /* 默认隐藏滚动条 */
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}

.scrollbar-thin:hover {
  /* 悬停时显示滚动条 */
  scrollbar-color: rgba(148, 163, 184, 0.3) transparent;
}

.scrollbar-thin::-webkit-scrollbar {
  width: 4px;  /* 从 6px 减小到 4px */
  height: 4px;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: transparent;  /* 默认透明 */
  border-radius: 9999px;
}

.scrollbar-thin:hover::-webkit-scrollbar-thumb {
  background-color: rgba(148, 163, 184, 0.3);  /* 悬停时半透明 */
}
```

---

### 2. ✅ 修复任务 item 按钮高度不一致

**问题**：完成按钮、编辑按钮、更多按钮高度不一样。

**解决方案**：
- 给所有按钮添加统一的固定高度 `h-9` (36px)
- 保持内边距 `py-2` 一致

**修改文件**：`frontend/src/modules/tasks/components/TaskCard.vue`

```vue
<!-- 所有按钮都添加 h-9 -->
<button class="action-btn h-9 rounded-lg ...">
```

---

### 3. ✅ 降低更多按钮弹框的透明度

**问题**：点击更多按钮弹出的菜单透明度太高，和背后的元素有干扰。

**解决方案**：
- 背景色从 `bg-slate-800/98` 改为 `bg-slate-900/95`（更深的颜色）
- 边框从 `border-white/20` 改为 `border-white/30`（更明显）
- 背景模糊从 `backdrop-blur-sm` 改为 `backdrop-blur-md`（更强的模糊）

**修改文件**：`frontend/src/modules/tasks/components/TaskCard.vue`

```vue
<div class="... bg-slate-900/95 border-white/30 backdrop-blur-md ...">
```

**效果对比**：

| 属性 | 优化前 | 优化后 |
|------|--------|--------|
| 背景色 | slate-800/98 | slate-900/95 |
| 边框 | white/20 | white/30 |
| 背景模糊 | blur-sm | blur-md |

---

### 4. ✅ 置顶功能暂时隐藏

**问题**：点击置顶没有生效。

**原因**：
- 后端 API 尚未实现 pinned 字段
- TaskDTO 类型定义中没有 pinned 属性

**解决方案**：
- 暂时注释掉置顶按钮
- 添加注释说明需要后端支持

**修改文件**：`frontend/src/modules/tasks/components/TaskCard.vue`

```vue
<!-- 置顶功能需要后端支持，暂时隐藏 -->
<!-- <button ...> 置顶 </button> -->
```

**后续计划**：
需要后端实现以下功能：
1. 在 Task 模型中添加 `pinned` 字段（boolean）
2. 添加 API 接口：`PUT /api/tasks/:uuid/pin`
3. 更新排序逻辑，置顶任务在最前面

---

### 5. ✅ 优化"移动"按钮文案

**问题**：更多按钮中的"移动"不够清楚，不知道移动到哪里。

**解决方案**：
- 根据当前任务状态动态显示文案
- "现在"任务显示"移到未来"
- "未来"任务显示"移到现在"

**修改文件**：`frontend/src/modules/tasks/components/TaskCard.vue`

```vue
<!-- 优化前 -->
<span>🔁</span> 移动

<!-- 优化后 -->
<span>🔁</span> 
<span>{{ task.status === 'now' ? '移到未来' : '移到现在' }}</span>
```

---

### 6. ✅ 修复编辑和新增按钮重叠

**问题**：编辑弹框和新增按钮在某些情况下有重叠。

**解决方案**：
- 调整 FloatingAddButton 的 z-index 从 40 降到 30
- 确保层级清晰：
  - FloatingAddButton: z-30（最底层）
  - TaskDrawer: z-40（中间）
  - TaskFormModal: z-50（最上层）
  - ToastStack: z-50（最上层）

**修改文件**：`frontend/src/modules/tasks/components/FloatingAddButton.vue`

```vue
<button class="... z-30 ...">
```

**z-index 层级规划**：

```
z-0    - 普通内容
z-10   - 下拉菜单
z-20   - 任务卡片菜单
z-30   - 浮动添加按钮
z-40   - 抽屉面板
z-50   - 模态框和 Toast
```

---

## 📊 优化效果对比

### 视觉改进

**滚动条**：
- 优化前：总是显示，6px 宽度
- 优化后：默认隐藏，悬停显示，4px 宽度

**按钮一致性**：
- 优化前：高度不一致，视觉不整齐
- 优化后：统一 36px 高度，对齐完美

**菜单可读性**：
- 优化前：背景太透明，内容模糊
- 优化后：背景不透明度 95%，清晰易读

**文案清晰度**：
- 优化前："移动"（不清楚移动到哪）
- 优化后："移到未来" / "移到现在"（明确目标）

---

## 🛠️ 修改文件清单

1. **frontend/src/styles/index.css**
   - 优化滚动条样式

2. **frontend/src/modules/tasks/components/TaskCard.vue**
   - 统一按钮高度
   - 增强菜单背景
   - 优化移动按钮文案
   - 隐藏置顶按钮

3. **frontend/src/modules/tasks/components/FloatingAddButton.vue**
   - 调整 z-index 层级

---

## ✅ 测试验证

### 功能测试

- [x] 鼠标悬停任务卡片，滚动条默认隐藏
- [x] 鼠标悬停滚动区域，滚动条显示
- [x] 所有按钮高度一致
- [x] 点击更多按钮，菜单清晰可读
- [x] "移动"按钮文案根据状态动态变化
- [x] 浮动添加按钮不与其他元素重叠
- [x] 置顶按钮已隐藏

### 视觉测试

- [x] 任务卡片整体美观
- [x] 按钮对齐完美
- [x] 菜单背景不透明
- [x] 层级关系正确

---

## 📝 注意事项

### 置顶功能后续实现

当后端支持 pinned 字段后，需要：

1. **更新类型定义**（`frontend/src/services/types.ts`）：
```typescript
export interface TaskDTO {
  ...
  pinned?: boolean;
}
```

2. **添加 API 方法**（`frontend/src/services/taskApi.ts`）：
```typescript
async pin(uuid: string, pinned: boolean) {
  return request.put(`/tasks/${uuid}/pin`, { pinned });
}
```

3. **取消注释置顶按钮**（`TaskCard.vue`）：
```vue
<button @click="handlePin">
  <span>⭐</span> {{ task.pinned ? '取消置顶' : '置顶' }}
</button>
```

4. **更新排序逻辑**（`frontend/src/utils/sort.ts`）：
```typescript
// 置顶任务优先
tasks.sort((a, b) => {
  if (a.pinned && !b.pinned) return -1;
  if (!a.pinned && b.pinned) return 1;
  // 其他排序逻辑...
});
```

---

## 🎯 用户体验提升

1. **更简洁的界面**：
   - 滚动条不再干扰视觉
   - 只在需要时显示

2. **更一致的设计**：
   - 按钮高度统一
   - 视觉更整齐

3. **更清晰的交互**：
   - 菜单背景不透明，易读
   - 操作文案明确

4. **更合理的层级**：
   - z-index 规划清晰
   - 元素不重叠

---

## 📈 性能影响

所有优化都是纯 CSS 和模板改动，没有引入新的 JavaScript 逻辑，对性能影响可以忽略。

---

**完成时间**: 2025-11-07  
**所有优化项**: 6/6 ✅  
**测试状态**: 全部通过 ✅

