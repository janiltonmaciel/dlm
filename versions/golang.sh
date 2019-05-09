#!/usr/bin/env bash
set -Eeuo pipefail

cd "$(dirname "$(readlink -f "$BASH_SOURCE")")"

# source '.architectures-lib'

versions=( "$@" )
if [ ${#versions[@]} -eq 0 ]; then
	versions=( */ )
fi
versions=( "${versions[@]%/}" )

# see http://stackoverflow.com/a/2705678/433558
sed_escape_lhs() {
	echo "$@" | sed -e 's/[]\/$*.^|[]/\\&/g'
}
sed_escape_rhs() {
	echo "$@" | sed -e 's/[\/&]/\\&/g' | sed -e ':a;N;$!ba;s/\n/\\n/g'
}

# https://github.com/golang/go/issues/13220
allGoVersions=()
apiBaseUrl='https://www.googleapis.com/storage/v1/b/golang/o?fields=nextPageToken,items%2Fname'
pageToken=
while [ "$pageToken" != 'null' ]; do
	page="$(curl -fsSL "$apiBaseUrl&pageToken=$pageToken")"
	allGoVersions+=( $(
		echo "$page" \
			| jq -r '.items[].name' \
			| grep -E '^go[0-9].*[.]src[.]tar[.]gz$' \
			| sed -r -e 's!^go!!' -e 's![.]src[.]tar[.]gz$!!'
	) )
	# TODO extract per-version "available binary tarballs" information while we've got it handy here?
	pageToken="$(echo "$page" | jq -r '.nextPageToken')"
done

for version in "${versions[@]}"; do
	rcVersion="${version%-rc}"
	rcGrepV='-v'
	if [ "$rcVersion" != "$version" ]; then
		rcGrepV=
	fi
	rcGrepV+=' -E'
	rcGrepExpr='beta|rc'

	fullVersion="$(
		echo "${allGoVersions[@]}" | xargs -n1 \
			| grep $rcGrepV -- "$rcGrepExpr" \
			| grep -E "^${rcVersion}([.a-z]|$)" \
			| sort -V \
			| tail -1
	)" || true
	if [ -z "$fullVersion" ]; then
		echo >&2 "warning: cannot find full version for $version"
		continue
	fi
	fullVersion="${fullVersion#go}" # strip "go" off "go1.4.2"

	# https://github.com/golang/build/commit/24f7399f96feb8dd2fc54f064e47a886c2f8bb4a
	srcSha256="$(curl -fsSL "https://storage.googleapis.com/golang/go${fullVersion}.src.tar.gz.sha256")"
	if [ -z "$srcSha256" ]; then
		echo >&2 "warning: cannot find sha256 for $fullVersion src tarball"
		continue
	fi

	for variant in \
		alpine3.7 \
	; do
		if [ -d "$version/$variant" ]; then
			tag="$variant"
			template='debian'
			case "$variant" in
				alpine*) tag="${variant#alpine}"; template='alpine' ;;
			esac

			sed -r \
				-e 's!%%VERSION%%!'"$fullVersion"'!g' \
				-e 's!%%TAG%%!'"$tag"'!g' \
				-e 's!%%SRC-SHA256%%!'"$srcSha256"'!g' \
				"Dockerfile-${template}.template" > "$version/$variant/Dockerfile"

		fi
	done

	echo "$version: $fullVersion ($srcSha256)"
done
