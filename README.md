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


<img width="531" alt="Screen Shot 2021-07-24 at 10 58 07" src="https://user-images.githubusercontent.com/19243302/126895841-99ad3ca7-fcc1-42e5-8094-50516b73ec21.png">
<img width="765" alt="Screen Shot 2021-07-24 at 10 58 29" src="https://user-images.githubusercontent.com/19243302/126895843-55d8f87d-2c3b-4f24-a70d-174ed85bec4c.png">


<img width="560" alt="Screen Shot 2021-06-28 at 18 14 35" src="https://user-images.githubusercontent.com/19243302/126895866-4cc8d000-69b4-4a78-b970-682403ffbe0b.png">
<img width="544" alt="Screen Shot 2021-06-06 at 10 34 58" src="https://user-images.githubusercontent.com/19243302/126895879-97af4744-2f93-4b8d-b71e-a78ebc03ce46.png">
