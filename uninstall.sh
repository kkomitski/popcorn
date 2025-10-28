#!/bin/bash

set -e

echo "🍿 Uninstalling Popcorn..."
echo ""

# Remove binary from ~/bin
if [ -f ~/bin/pop ]; then
    echo "�️  Removing binary from ~/bin..."
    rm ~/bin/pop
    echo "✅ Binary removed"
else
    echo "⚠️  Binary not found in ~/bin"
fi

# Ask if user wants to remove PATH entry from .zshrc
echo ""
read -p "Remove PATH entry from ~/.zshrc? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "� Removing PATH entry from ~/.zshrc..."
    # Remove the popcorn installer lines from .zshrc
    sed -i.bak '/# Added by popcorn installer/,+1d' ~/.zshrc
    echo "✅ PATH entry removed (backup saved as ~/.zshrc.bak)"
    echo "   Run: source ~/.zshrc"
else
    echo "⏭️  Skipping PATH removal"
fi

echo ""
echo "� Popcorn uninstalled. See you next time! 👋"
