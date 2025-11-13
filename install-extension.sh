#!/bin/bash

set -e

echo "🍿 Installing Popcorn VS Code Extension..."

# Navigate to extension directory
cd "$(dirname "$0")/extension"

echo ""
echo "🔨 Packaging extension..."
# Install vsce if not present
if ! command -v vsce &> /dev/null; then
    echo "📦 Installing vsce (VS Code Extension Manager)..."
    npm install -g @vscode/vsce
fi

# Package the extension
vsce package --allow-missing-repository

# Get the VSIX file name
VSIX_FILE=$(ls -t *.vsix | head -n 1)

if [ -z "$VSIX_FILE" ]; then
    echo "❌ Failed to create VSIX package"
    exit 1
fi

echo "✅ Created package: $VSIX_FILE"

# Prompt user to select profile or use default
echo ""
echo "The extension will be installed to your currently active VS Code profile."
read -p "Continue? [Y/n]: " confirm
if [ "$confirm" = "n" ] || [ "$confirm" = "N" ]; then
    echo "Installation cancelled."
    exit 0
fi

# Install the extension
echo "🔗 Installing extension to VS Code..."
code --install-extension "$VSIX_FILE" --force

# Find the installed extension directory
INSTALLED_EXT=$(find "$HOME/.vscode/extensions" -name "popcorn.popcorn-highlight-*" -type d 2>/dev/null | head -n 1)

if [ -z "$INSTALLED_EXT" ]; then
    echo "❌ Extension installation failed - extension not found in .vscode/extensions"
    exit 1
fi

echo "📍 Extension installed at: $INSTALLED_EXT"

# Add to all profile extensions.json files
PROFILES_DIR="$HOME/Library/Application Support/Code/User/profiles"
if [ -d "$PROFILES_DIR" ]; then
    echo ""
    echo "📝 Registering extension in profile(s)..."
    
    # Check if jq is available
    if ! command -v jq &> /dev/null; then
        echo "⚠️  jq not found. Installing jq..."
        brew install jq 2>/dev/null || {
            echo "❌ Failed to install jq. Please install it manually: brew install jq"
            exit 1
        }
    fi
    
    for profile_dir in "$PROFILES_DIR"/*; do
        if [ -d "$profile_dir" ]; then
            EXTENSIONS_JSON="$profile_dir/extensions.json"
            if [ -f "$EXTENSIONS_JSON" ]; then
                PROFILE_NAME=$(basename "$profile_dir")
                
                # Check if extension already exists
                if jq -e '.[] | select(.identifier.id == "popcorn.popcorn-highlight")' "$EXTENSIONS_JSON" > /dev/null 2>&1; then
                    echo "   ℹ️  Already registered in profile: $PROFILE_NAME"
                else
                    # Create new entry and append to array
                    TIMESTAMP=$(date +%s)000
                    REL_LOCATION=$(basename "$INSTALLED_EXT")
                    
                    jq --arg path "$INSTALLED_EXT" \
                       --arg rel "$REL_LOCATION" \
                       --arg ts "$TIMESTAMP" \
                       '. += [{
                           "identifier": {
                               "id": "popcorn.popcorn-highlight",
                               "uuid": "00000000-0000-0000-0000-popcorn00001"
                           },
                           "version": "0.0.1",
                           "location": {
                               "$mid": 1,
                               "path": $path,
                               "scheme": "file"
                           },
                           "relativeLocation": $rel,
                           "metadata": {
                               "installedTimestamp": ($ts | tonumber),
                               "pinned": false,
                               "source": "vsix",
                               "targetPlatform": "undefined",
                               "updated": false,
                               "private": false,
                               "isPreReleaseVersion": false,
                               "hasPreReleaseVersion": false
                           }
                       }]' "$EXTENSIONS_JSON" > "${EXTENSIONS_JSON}.tmp" && mv "${EXTENSIONS_JSON}.tmp" "$EXTENSIONS_JSON"
                    
                    echo "   ✅ Registered in profile: $PROFILE_NAME"
                fi
            fi
        fi
    done
fi

echo ""
echo "✅ Popcorn extension installed successfully!"
echo ""
echo "📝 Please reload VS Code to activate the extension."
echo "   Press Cmd+Shift+P → 'Developer: Reload Window'"
echo "   Then open any .pop file to see syntax highlighting!"
