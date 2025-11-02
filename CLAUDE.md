# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A minimalist todo list application with three core views:
- **Now** (现在): Current tasks that must be done (suggested limit: 5 tasks)
- **Future** (未来): Tasks to be done later
- **History** (历史): Completed tasks archive

## Core Data Model

**Task Object:**
- `title` (required): Task title
- `deadline` (optional): Deadline date
- `notes` (optional): Additional notes
- Task group: One of "Now", "Future", or "History"

## Key Architectural Principles

1. **Sorting Rules:**
   - **Now**: Tasks with deadlines first (nearest first), then by addition order
   - **Future**: Tasks with deadlines first (ascending), then no-deadline tasks by addition time
   - **History**: By completion time (descending)

2. **Deadline Handling:**
   - Overdue tasks display in red in their current group (no forced movement)
   - Display format: "today"/"tomorrow"/date (YYYY-MM-DD)
   - No complex calendar logic

3. **Capacity Management:**
   - "Now" view suggests max 5 tasks (warning when exceeded, still allows more)

## Essential Operations

**Task Movement:**
- Future → Now (start working)
- Now → History (complete)
- Future → History (skip/no longer needed)

**Task Actions:**
- Add (to Future or Now)
- Edit (title, deadline, notes)
- Complete (moves to History)
- Delete (permanent removal, different from complete)
- Search/filter by keyword

**Undo Support:**
- Recent operations should be undoable (3-5 second window)
- Display undo option in Toast notifications

## Interaction Design

**Primary Interactions:**
- Drag & drop between views (with target area highlighting)
- Click title for inline editing (Enter to save, Esc to cancel)
- Checkbox/button to complete tasks
- Right-click menu or three-dot menu for actions

**Keyboard Shortcuts:**
- `N`: Quick add (to Future)
- `Shift+N`: Add to Now
- `1/2/3`: Switch to Now/Future/History
- `Enter`: Edit selected task / Save in edit mode
- `C`: Mark complete
- `Del`: Delete (with confirmation)
- `/`: Focus search
- `Z` or `Ctrl+Z`: Undo

**Batch Operations:**
- Multi-select with `Shift` (range) or `Ctrl/Cmd` (individual)
- `M`: Move selected tasks
- `D`: Set deadline for selected
- `C`: Complete selected

## Empty States

- **Now empty**: Show "Pull some tasks from Future" button
- **Future empty**: Show "Add task" quick entry
- **History empty**: Simple "No history yet" message

## UI Feedback

- Lightweight Toast notifications for all operations
- Include undo option in Toasts
- No blocking dialogs except delete confirmation
- Overdue tasks: red dot + red deadline text
- History items: 60-70% opacity

## Out of Scope (Do Not Implement)

- Tags/projects/priority matrices
- Subtasks/dependencies/recurring tasks
- Reminders/push notifications/calendar sync/statistics
- Multi-user/permissions
- Rich text/attachments

## Reference Document

See `/docs/需求说明.md` for complete Chinese specification including:
- Detailed interaction patterns
- Touch gestures for mobile
- User flow scenarios
- Complete checklist for validation
