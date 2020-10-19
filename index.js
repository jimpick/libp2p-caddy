import pako from 'pako'

document.addEventListener('DOMContentLoaded', async () => {
  // UI elements
  const status = document.getElementById('status')
  const output = document.getElementById('output')

  output.textContent = ''

  function log (txt) {
    console.info(txt)
    output.textContent += `${txt.trim()}\n`
  }

  status.innerText = 'Loading Go WASM bundle...'

  // Use gzip: https://dstoiko.github.io/posts/go-pong-wasm/

  const go = new Go()
  const url = 'go-wasm/main.wasm.gz' // the gzip-compressed wasm file
  let wasm = pako.ungzip(await (await fetch(url)).arrayBuffer())
  // A fetched response might be decompressed twice on Firefox.
  // See https://bugzilla.mozilla.org/show_bug.cgi?id=610679
  if (wasm[0] === 0x1f && wasm[1] === 0x8b) {
    wasm = pako.ungzip(wasm)
  }
  WebAssembly.instantiate(wasm, go.importObject).then(result => {
    status.innerText = 'All systems good! Go WASM loaded.'
    go.run(result.instance)

    // FIXME Add callback for ready state
    const goPingButton = document.querySelector('#goPingBtn')
    goPingButton.disabled = false
    goPingButton.onclick = async function () {
      const target = document.querySelector('#maddr').value
      log(`Go Ping: ${target}`)
      const latency = await window.ping(target)
      log(`Go Pong: ${latency}ms`)
    }

    const goGraphSyncButton = document.querySelector('#goGraphSyncBtn')
    goGraphSyncButton.disabled = false
    goGraphSyncButton.onclick = async function () {
      const target = document.querySelector('#maddr').value
      if (!target.match(/ipfs/)) {
        alert('Use ipfs maddr')
        return
      }
      const cid = 'QmeqtCLGLNWsK5djgEN76F2z7gLodWCaWesupKrnGA4TWf'
      log(`Go GraphSync fetch:\n`)
      log(`> Maddr: ${target}`)
      log(`> CID: ${cid}`)
      const data = await window.graphSyncFetch(target, cid)
      log(`> Data:\n${data}\n`)
    }
  })
})
