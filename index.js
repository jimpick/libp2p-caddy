import 'babel-polyfill'
import Libp2p from 'libp2p'
import Websockets from 'libp2p-websockets'
import { NOISE } from 'libp2p-noise'
import Mplex from 'libp2p-mplex'

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
    output.textContent += `${txt.trim()}\n`
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
  status.innerText = 'js libp2p started!'
  log(`JS libp2p id is ${libp2p.peerId.toB58String()}`)

  const jsPingButton = document.querySelector('#jsPingBtn')
  jsPingButton.onclick = async function () {
    const target = document.querySelector('#maddr').value
    log(`JS Ping: ${target}`)
    const latency = await libp2p.ping(target)
    log(`JS Pong: ${latency}ms`)
  }

  const goPingButton = document.querySelector('#goPingBtn')
  const state = {
    counter: 0
  }
  goPingButton.onclick = async function () {
    // const target = document.querySelector('#maddr').value
    // log(`Go Ping: ${target}`)
    state.counter = window.add(state.counter, 1)
    log(`Go Counter: ${state.counter}`)
    /*
    const latency = await libp2p.ping(target)
    log(`JS Pong: ${latency}ms`)
    */
  }

  // Export libp2p to the window so you can play with the API
  window.libp2p = libp2p
})
