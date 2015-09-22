package main

const NGINX_BUILD_VERSION = "0.5.0"

// nginx
const (
	NGINX_VERSION             = "1.9.5"
	NGINX_DOWNLOAD_URL_PREFIX = "http://nginx.org/download"
)

// pcre
const (
	PCRE_VERSION             = "8.37"
	PCRE_DOWNLOAD_URL_PREFIX = "http://ftp.csx.cam.ac.uk/pub/software/programming/pcre"
)

// openssl
const (
	OPENSSL_VERSION             = "1.0.2d"
	OPENSSL_DOWNLOAD_URL_PREFIX = "http://www.openssl.org/source"
)

// zlib
const (
	ZLIB_VERSION             = "1.2.8"
	ZLIB_DOWNLOAD_URL_PREFIX = "http://zlib.net"
)

// openResty
const (
	OPENRESTY_VERSION             = "1.9.3.1"
	OPENRESTY_DOWNLOAD_URL_PREFIX = "http://openresty.org/download"
)

// tengine
const (
	TENGINE_VERSION             = "2.1.1"
	TENGINE_DOWNLOAD_URL_PREFIX = "http://tengine.taobao.org/download"
)

// component enumerations
const (
	COMPONENT_NGINX = iota
	COMPONENT_OPENRESTY
	COMPONENT_TENGINE
	COMPONENT_PCRE
	COMPONENT_OPENSSL
	COMPONENT_ZLIB
	COMPONENT_MAX
)
