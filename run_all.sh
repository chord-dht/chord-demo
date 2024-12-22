#!/bin/bash

# Name of the tmux session
SESSION_NAME="chord"

# Directory of the current script
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

LOCAL_IP="127.0.0.1"

# Base port for services
BASE_PORT=4170

# Check if TLS is enabled and set the corresponding option
if [ "$TLS" == "true" ]; then
  TLS=-tls -cacert "cacert.pem" -servercert "peer.crt" -serverkey "peer.key"
else
  TLS=""
fi

# Start a new tmux session with the specified session name
tmux new-session -d -s $SESSION_NAME -n window1

# Create eight panes by splitting the first window
tmux split-window -h -t $SESSION_NAME:window1
tmux split-window -h -t $SESSION_NAME:window1
tmux split-window -h -t $SESSION_NAME:window1
tmux split-window -h -t $SESSION_NAME:window1
tmux split-window -h -t $SESSION_NAME:window1
tmux split-window -v -t $SESSION_NAME:window1
tmux split-window -v -t $SESSION_NAME:window1

# Adjust the layout to tiled so all panes are evenly distributed
tmux select-layout -t $SESSION_NAME:window1 tiled

# Loop through 8 panes and send commands to each one
for i in $(seq 0 7); do
    # Calculate the worker index and assign server and port
    index=$((i+1))
    PORT=$((BASE_PORT + i))

    # Start Process with different configurations for the first and other workers
    if [[ $i -eq 0 ]]; then
      # For the worker in create mode
      tmux send-keys -t $SESSION_NAME:window1."$i" "cd $SCRIPT_DIR/pack/peer_$index && ./chord -a $LOCAL_IP -p $PORT --ts 3000 --tff 1000 --tcp 3000 -r 4 $TLS" C-m
    else
      tmux send-keys -t $SESSION_NAME:window1."$i" "sleep 10" C-m
      # For other workers in join mode
      tmux send-keys -t $SESSION_NAME:window1."$i" "cd $SCRIPT_DIR/pack/peer_$index && ./chord --ja $LOCAL_IP --jp $BASE_PORT -a $LOCAL_IP -p $PORT --ts 3000 --tff 1000 --tcp 3000 -r 4 $TLS" C-m
    fi
done

# Attach to the tmux session to monitor the panes
tmux attach-session -t $SESSION_NAME