export SSH_AUTH_SOCK=$HOME/.ssh/agent.sock

# Start a new SSH agent if the socket is not reachable
ssh-add -l >/dev/null 2>&1
if [ $? -eq 2 ]; then
    rm -f "$SSH_AUTH_SOCK"
    eval $(ssh-agent -s -a "$SSH_AUTH_SOCK") >/dev/null
fi

# Add the default ed25519 key if present (harmless if already loaded)
ssh-add "$HOME/.ssh/id_ed25519" 2>/dev/null
