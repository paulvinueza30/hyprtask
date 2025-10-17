# HyprTask

A terminal-based task manager for Hyprland window manager.

## Requirements

### Minimum Terminal Dimensions
- **Width**: 65 characters
- **Height**: 20 lines

The application is designed to work optimally with these minimum dimensions. Smaller terminals may experience layout issues or limited functionality.

## Features

- Responsive workspace selector
- Dynamic padding based on terminal size
- Keyboard navigation and scrolling
- Real-time workspace monitoring

## Usage

```bash
./hyprtask
```

## Controls

Use your configured keybinds for navigation and selection.

## Development

Built with Go using the Bubble Tea framework for terminal user interfaces.

## Process List Design (TODO)

### Simple Single Model Approach

**Messages:**
```go
type ViewWorkspaceProcsMsg struct {
    WorkspaceID int
}

type ViewAllProcsMsg struct{}
```

**ProcessList Model:**
```go
type ProcessList struct {
    processes []Process
    viewMode string // "all" | "workspace_1" | "workspace_2"
}
```

**Flow:**
- WorkspaceSelector → ViewWorkspaceProcsMsg{WorkspaceID: 1} → ProcessList updates with filtered data
- WorkspaceSelector → ViewAllProcsMsg{} → ProcessList updates with all data

**Benefits:**
- Simple, clean, easy to understand
- Single model handles both "all processes" and "workspace-specific processes"
- Easy to implement and debug

**Future Optimization (LRU Cache):**
- Keep "all processes" always loaded
- LRU cache for workspace-specific views
- Fast switching between cached views
- Only implement if performance becomes an issue
