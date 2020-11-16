import Libp2p from 'libp2p'
import Websockets from 'libp2p-websockets'
import { NOISE } from 'libp2p-noise'
import Mplex from 'libp2p-mplex'
import { BrowserProvider } from './browser-provider'
import { LotusRPC } from '@filecoin-shipyard/lotus-client-rpc'
import { mainnet } from '@filecoin-shipyard/lotus-client-schema'
import { WasmProvider } from './wasm-provider'

document.addEventListener('DOMContentLoaded', async () => {
  // Create our libp2p node
  const libp2p = await Libp2p.create({
    modules: {
      transport: [Websockets],
      connEncryption: [NOISE],
      streamMuxer: [Mplex]
    }
  })

  // UI elements
  const status = document.getElementById('status')
  const output = document.getElementById('output')

  output.textContent = ''

  function log (txt) {
    console.info(txt)
    const divEl = document.createElement('div')
    divEl.textContent = txt.trim()
    output.appendChild(divEl)
  }

  // Listen for new connections to peers
  libp2p.connectionManager.on('peer:connect', connection => {
    log(`JS Connected to ${connection.remotePeer.toB58String()}`)
  })

  // Listen for peers disconnecting
  libp2p.connectionManager.on('peer:disconnect', connection => {
    log(`Disconnected from ${connection.remotePeer.toB58String()}`)
  })

  await libp2p.start()
  status.innerText = 'js libp2p started, loading go WASM bundle...'
  log(`JS libp2p id is ${libp2p.peerId.toB58String()}`)

  const jsPingButton = document.querySelector('#jsPingBtn')
  jsPingButton.disabled = false
  jsPingButton.onclick = async function () {
    const target = document.querySelector('#maddr').value
    log(`JS Ping: ${target}`)
    const latency = await libp2p.ping(target)
    log(`JS Pong: ${latency}ms`)
  }

  const setWsMaddrButton = document.querySelector('#setWsMaddrBtn')
  setWsMaddrButton.onclick = async function () {
    document.querySelector('#maddr').value =
      '/dns4/libp2p-caddy-ws.localhost/tcp/9056/wss/p2p/' +
      process.env.PEER_ID_WS
  }

  const setTcpMaddrButton = document.querySelector('#setTcpMaddrBtn')
  setTcpMaddrButton.onclick = async function () {
    document.querySelector('#maddr').value =
      '/ip4/127.0.0.1/tcp/2020/p2p/' + process.env.PEER_ID_TCP
  }

  const setIpfsMaddrButton = document.querySelector('#setIpfsMaddrBtn')
  setIpfsMaddrButton.onclick = async function () {
    document.querySelector('#maddr').value =
      '/dns4/ipfs.jimpick.com/tcp/10100/wss/p2p/QmScdku7gc3VvfZZvT8kHU77bt6bnH3PnGXkyFRZ17g9EG'
  }

  const setMinerMaddrButton = document.querySelector('#setMinerMaddrBtn')
  setMinerMaddrButton.onclick = async function () {
    document.querySelector('#maddr').value =
      '/dns4/lotus.jimpick.com/tcp/10101/wss/p2p/12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5'
  }

  // Export libp2p to the window so you can play with the API
  window.libp2p = libp2p
  console.log('.env PEER_ID_WS:', process.env.PEER_ID_WS)
  console.log('.env PEER_ID_TCP:', process.env.PEER_ID_TCP)
  document.querySelector('#maddr').value =
    '/dns4/libp2p-caddy-ws.localhost/tcp/9056/wss/p2p/' + process.env.PEER_ID_WS

  const go = new Go()
  const result = await WebAssembly.instantiateStreaming(
    fetch('go-wasm/main.wasm'),
    go.importObject
  )
  status.innerText = 'All systems good! JS and Go loaded.'
  go.run(result.instance)

  const goPingButton = document.querySelector('#goPingBtn')
  goPingButton.disabled = false
  goPingButton.onclick = async function () {
    const target = document.querySelector('#maddr').value
    log(`Go Ping: ${target}`)
    const latency = await window.ping(target)
    log(`Go Pong: ${latency}ms`)
  }

  const goStreamWsButton = document.querySelector('#goStreamWsBtn')
  goStreamWsButton.disabled = false
  goStreamWsButton.onclick = async function () {
    const target = document.querySelector('#maddr').value
    log(`Go Stream From WS: ${target}`)
    const result = await window.streamFromWs(target)
    log(`Go Stream from WS result: ${result}`)
  }

  const goStreamTcpButton = document.querySelector('#goStreamTcpBtn')
  goStreamTcpButton.disabled = false
  goStreamTcpButton.onclick = async function () {
    const target = document.querySelector('#maddr').value
    log(`Go Stream From TCP: ${target}`)
    const result = await window.streamFromTcp(target)
    log(`Go Stream from TCP result: ${result}`)
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

  const wsUrl = 'wss://lotus.jimpick.com/spacerace_api/0/node/rpc/v0'
  const browserProvider = new BrowserProvider(wsUrl)
  await browserProvider.connect()
  const requestsForLotusHandler = async (req, responseHandler) => {
    const request = JSON.parse(req)
    console.log('JSON-RPC request => Lotus', request)
    async function waitForResult () {
      const result = await browserProvider.sendWs(request)
      console.log('Jim result', result)
      responseHandler(JSON.stringify(result))
    }
    waitForResult()
    // return 'abcde'
  }

  const goChainHeadButton = document.querySelector('#goChainHeadBtn')
  goChainHeadButton.disabled = false
  goChainHeadButton.onclick = async function () {
    log(`Go ChainHead`)
    const result = await window.chainHead(requestsForLotusHandler)
    log(`Go ChainHead: ${result}`)
  }

  const schema = {
    methods: {
      HelloName: {}
    }
  }

  const wasmHelloServiceProvider = new WasmProvider(window.connectHelloService)
  const helloClient = new LotusRPC(wasmHelloServiceProvider, { schema })
  const goHelloButton = document.querySelector('#goHelloBtn')
  goHelloButton.disabled = false
  goHelloButton.onclick = async function () {
    log(`Go Hello`)
    const result = await helloClient.helloName('Jim')
    log(`Go Hello: ${JSON.stringify(result)}`)
  }

  const wasmQueryAskServiceProvider = new WasmProvider(
    window.connectQueryAskService,
    {
      environment: {
        requestsForLotusHandler
      }
    }
  )

  const queryAskClient = new LotusRPC(wasmQueryAskServiceProvider, {
    schema: mainnet.fullNode
  })

  queryAskBtn.disabled = false
  queryAskBtn.onclick = async function () {
    log(`Query Ask WSS`)
    const result = await queryAskClient.clientQueryAsk(
      '12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5',
      'f063655'
    )
    log(`Query Ask WSS: ${JSON.stringify(result)}`)
  }

  queryAskTcpBtn.disabled = false
  queryAskTcpBtn.onclick = async function () {
    log(`Query Ask TCP`)
    const result = await queryAskClient.clientQueryAsk(
      '12D3KooWDMpcct12Vb6jPXwjvLQHA2hoP8XKGbUZ2tpue1ydoZUm',
      'f02620'
    )
    log(`Query Ask TCP: ${JSON.stringify(result)}`)
  }
})
