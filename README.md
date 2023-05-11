# WebAssembly Port Scanner
Written in Go with target WASM/WASI.<br>

## Demo:
http://ports.sh/

## QuickStart
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


<img width="765" alt="" src="https://user-images.githubusercontent.com/19243302/126895841-99ad3ca7-fcc1-42e5-8094-50516b73ec21.png">
<img width="765" alt="" src="https://user-images.githubusercontent.com/19243302/145462240-56038b75-0bfd-4fcb-95c3-f60c3ab3b3e8.png">
<img width="765" alt="" src="https://user-images.githubusercontent.com/19243302/126895866-4cc8d000-69b4-4a78-b970-682403ffbe0b.png">
<img width="765" alt="" src="https://user-images.githubusercontent.com/19243302/126895879-97af4744-2f93-4b8d-b71e-a78ebc03ce46.png">
