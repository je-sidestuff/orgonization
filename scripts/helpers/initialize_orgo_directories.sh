#!/bin/bash

set -euo pipefail

# If ~/notesa dfoes not exist, create it
if [ ! -d "$HOME/notes" ]; then
    mkdir "$HOME/notes"

    # Add some content to ~/notes/sample.md
    cat << EOF > "$HOME/notes/sample.md"
# This is a sample note

It contains text and a <PING:> token.

Please enjoy.
EOF
fi

# If ~/.orgo/outfiles does not exist, create it
if [ ! -d "$HOME/.orgo/outfiles" ]; then
    mkdir -p "$HOME/.orgo/outfiles"
fi
