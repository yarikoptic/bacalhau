#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
SESSION=bacalhau-dashboard
export APP=${APP:=""}
export PREDICTABLE_API_PORT=1

function start() {
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    echo "Session $SESSION already exists. Attaching..."
    sleep 1
    tmux -2 attach -t $SESSION
    exit 0;
  fi

  echo "Creating tmux session $SESSION..."

  rm -f /tmp/bacalhau-devstack.{port,pid}

  # get the size of the window and create a session at that size
  local screensize=$(stty size)
  local width=$(echo -n "$screensize" | awk '{print $2}')
  local height=$(echo -n "$screensize" | awk '{print $1}')
  tmux -2 new-session -d -s $SESSION -x "$width" -y "$(($height - 1))"

  # the right hand col with a 50% vertical split
  tmux split-window -h -d
  tmux select-pane -t 1
  tmux split-window -v -d
  tmux select-pane -t 0
  tmux split-window -v -d
  tmux split-window -v -d

  tmux send-keys -t 0 'make devstack' C-m
  tmux send-keys -t 1 'cd dashboard/frontend' C-m
  tmux send-keys -t 1 'yarn dev' C-m
  tmux send-keys -t 2 'cd dashboard/api' C-m
  tmux send-keys -t 2 'sleep 5' C-m
  tmux send-keys -t 2 'go run main.go 127.0.0.1 20000 20000 127.0.0.1 20001 20001 127.0.0.1 20002 20002' C-m

  tmux -2 attach-session -t $SESSION
}

function stop() {
  echo "Stopping tmux session $SESSION..."
  tmux kill-session -t $SESSION
}

command="$@"

if [ -z "$command" ]; then
  command="start"
fi

eval "$command"