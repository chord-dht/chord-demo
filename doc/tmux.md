# tmux Quick Guide

## Enable Mouse Support

To enable mouse interactions in tmux, modify your ~/.tmux.conf file:

```shell
set -g mouse on
```

Apply the changes:

```shell
tmux source-file ~/.tmux.conf
```

With mouse support enabled, you can:

• Click to switch panes.

• Drag pane borders to resize.

• Scroll through pane history using the mouse wheel.

## Start a tmux Session

To start a new tmux session:

```shell
tmux new-session -d -s $SESSION_NAME -n window1
```

- -d: Detach from the session after creating it.
- -s: $SESSION_NAME: Name the session.
- -n: window1: Name the first window.

## Attach to an Existing Session

If you have an existing session, you can attach to it using:

```shell
tmux attach-session -t $SESSION_NAME
```

List all active sessions:

```shell
tmux list-sessions
```

## Close a tmux Session

Kill a specific session:

```shell
tmux kill-session -t $SESSION_NAME
```

Kill all sessions:

```shell
tmux kill-server
```

## Split Panes

Horizontally (Left/Right Split):

```shell
Ctrl-b %
```

Vertically (Top/Bottom Split):

```shell
Ctrl-b "
```

### Using commands

```shell

tmux split-window -h -t $SESSION_NAME:$WINDOW_NAME
tmux split-window -v -t $SESSION_NAME:$WINDOW_NAME
```

- -h : Split the window horizontally.
- -v : Split the window vertically.
- -t : Target pane.

## Display Pane Numbers

Show Pane Numbers Temporarily:

```shell
Ctrl-b q
```
