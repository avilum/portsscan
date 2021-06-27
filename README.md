# WebAssembly Port Scanner
Written in Go with target WASM/WASI.<br>
The WASM main function scans all the open ports in the specified range (see <code>main.go</code>), via 0.0.0.0 with no-cors fetch mode in Javascript level.<br>
* Discovers any TCP open port available on the visiting host.<br>
* One byte of response / filtered port is enough
* Scans TCP only (WASM has no UDP support yet)
* Uses golang 'http' API rather then 'net' API (better browser compatibility)
#### Setup
Please see <code>./build.sh</code>
#### Build and Run
Simply start an HTTP server locally, for example:
<br><code>python3 -m http.server 5000</code><br>Or:<br><code>npm i -g serve && serve</code><br>

