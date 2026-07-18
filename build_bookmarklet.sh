#!/usr/bin/env bash
set -euo pipefail

IN_JS="${1}"
OUT_BOOKMARKLET="${2}"

if [[ ! -f "$IN_JS" ]]; then
	echo "input file not found: $IN_JS" >&2
	exit 1
fi

rawurlencode() {
	local input="$1"
	local out=""
	local i c hex

	LC_ALL=C
	for ((i = 0; i < ${#input}; i++)); do
		c="${input:i:1}"
		case "$c" in
		[a-zA-Z0-9.~_-] | "!" | "*" | "'" | "(" | ")")
			out+="$c"
			;;
		*)
			printf -v hex '%02X' "'$c"
			out+="%$hex"
			;;
		esac
	done

	printf '%s' "$out"
}

code="$(cat "$IN_JS")"
encoded="$(rawurlencode "$code")"
mkdir -p "$(dirname "$OUT_BOOKMARKLET")"
printf 'javascript:%s\n' "$encoded" >"$OUT_BOOKMARKLET"

echo "Built: $OUT_BOOKMARKLET"
