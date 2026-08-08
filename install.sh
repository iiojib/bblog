#!/usr/bin/env sh
set -eu

VERSION="${1:-0.3.1}"
BINDIR="/usr/local/bin"

case "$(uname -s)" in Linux) os=linux ;; Darwin) os=darwin ;; *) exit 1 ;; esac
case "$(uname -m)" in x86_64|amd64) arch=amd64 ;; arm64|aarch64) arch=arm64 ;; *) exit 1 ;; esac

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' 0
curl -fsSL "https://github.com/iiojib/bblog/releases/download/v${VERSION#v}/bblog-${os}-${arch}.tar.gz" -o "$tmp/bblog.tar.gz"
tar -xzf "$tmp/bblog.tar.gz" -C "$tmp"

if ! { [ -w "$BINDIR" ] || { [ ! -e "$BINDIR" ] && [ -w "$(dirname "$BINDIR")" ]; }; }; then
  BINDIR="$HOME/.local/bin"
fi

mkdir -p "$BINDIR"
cp "$tmp/bblog" "$BINDIR/bblog"

echo "installed bblog ${VERSION#v} to ${BINDIR}/bblog"
if [ "$BINDIR" = "$HOME/.local/bin" ]; then
  case ":$PATH:" in
    *":$HOME/.local/bin:"*) ;;
    *) echo "add $HOME/.local/bin to PATH" ;;
  esac
fi
